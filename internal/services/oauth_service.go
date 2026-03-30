package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/okapi"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	config       *config.Config
	userRepo     *repository.UserRepository
	providerRepo *repository.OAuthProviderRepository
	states       *stateStore
}

// stateStore is a simple in-memory CSRF state store with expiry.
type stateStore struct {
	mu     sync.Mutex
	states map[string]time.Time
}

func newStateStore() *stateStore {
	return &stateStore{states: make(map[string]time.Time)}
}

func (s *stateStore) Generate() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	state := base64.URLEncoding.EncodeToString(b)
	s.mu.Lock()
	s.states[state] = time.Now().Add(10 * time.Minute)
	s.mu.Unlock()
	return state, nil
}

func (s *stateStore) Validate(state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.states[state]
	if !ok {
		return false
	}
	delete(s.states, state)
	return time.Now().Before(exp)
}

func NewOAuthService(conf *config.Config) *OAuthService {
	return &OAuthService{
		config:       conf,
		userRepo:     repository.NewUserRepository(conf.Database.DB),
		providerRepo: repository.NewOAuthProviderRepository(conf.Database.DB),
		states:       newStateStore(),
	}
}

//  Admin endpoints

// GetProvider returns the configured OAuth provider (admin view).
func (s *OAuthService) GetProvider(c *okapi.Context) error {
	p, err := s.providerRepo.Get(c.Request().Context())
	if err != nil {
		return c.AbortInternalServerError("Failed to fetch OAuth provider", err)
	}
	if p == nil {
		return c.OK(nil)
	}
	return c.OK(toProviderDetail(p))
}

// SaveProvider creates or updates the OAuth provider config.
func (s *OAuthService) SaveProvider(c *okapi.Context, input *dto.SaveOAuthProviderRequest) error {
	ctx := c.Request().Context()

	// If updating, preserve the existing secret when the placeholder is sent
	clientSecret := input.Body.ClientSecret
	if clientSecret == "__unchanged__" {
		existing, err := s.providerRepo.Get(ctx)
		if err == nil && existing != nil {
			clientSecret = existing.ClientSecret
		}
	}

	p := &models.OAuthProvider{
		Name:         input.Body.Name,
		DisplayName:  input.Body.DisplayName,
		ClientID:     input.Body.ClientID,
		ClientSecret: clientSecret,
		AuthURL:      input.Body.AuthURL,
		TokenURL:     input.Body.TokenURL,
		UserInfoURL:  input.Body.UserInfoURL,
		Scopes:       models.StringArray(input.Body.Scopes),
		UserIDField:  input.Body.UserIDField,
		EmailField:   input.Body.EmailField,
		NameField:    input.Body.NameField,
		AvatarField:  input.Body.AvatarField,
		Enabled:      input.Body.Enabled,
	}
	p.Defaults()

	if err := s.providerRepo.Save(ctx, p); err != nil {
		return c.AbortInternalServerError("Failed to save OAuth provider", err)
	}

	return c.OK(toProviderDetail(p))
}

// DeleteProvider removes the OAuth provider config.
func (s *OAuthService) DeleteProvider(c *okapi.Context) error {
	if err := s.providerRepo.Delete(c.Request().Context()); err != nil {
		return c.AbortInternalServerError("Failed to delete OAuth provider", err)
	}
	return c.NoContent()
}

//  Public endpoints

// ProviderInfo returns the public OAuth provider info (name + auth URL).
func (s *OAuthService) ProviderInfo(c *okapi.Context) error {
	p, err := s.providerRepo.Get(c.Request().Context())
	if err != nil {
		return c.AbortInternalServerError("Failed to fetch OAuth provider", err)
	}
	if p == nil || !p.Enabled {
		return c.OK(nil)
	}

	state, err := s.states.Generate()
	if err != nil {
		return c.AbortInternalServerError("Failed to generate state", err)
	}

	cfg := s.buildOAuth2Config(p)
	return c.OK(dto.OAuthProviderInfo{
		Name:        p.Name,
		DisplayName: p.DisplayName,
		AuthURL:     cfg.AuthCodeURL(state, oauth2.AccessTypeOffline),
		Enabled:     true,
	})
}

// Authorize redirects to the OAuth provider's authorization page.
func (s *OAuthService) Authorize(c *okapi.Context) error {
	p, err := s.providerRepo.Get(c.Request().Context())
	if err != nil || p == nil || !p.Enabled {
		return c.AbortNotFound("OAuth provider not configured or disabled")
	}

	state, err := s.states.Generate()
	if err != nil {
		return c.AbortInternalServerError("Failed to generate state", err)
	}

	cfg := s.buildOAuth2Config(p)
	url := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
	return nil
}

// Callback handles the OAuth callback and issues a JWT.
func (s *OAuthService) Callback(c *okapi.Context, input *dto.OAuthCallbackRequest) error {
	p, err := s.providerRepo.Get(c.Request().Context())
	if err != nil || p == nil || !p.Enabled {
		return c.AbortNotFound("OAuth provider not configured or disabled")
	}

	if !s.states.Validate(input.State) {
		return c.AbortBadRequest("Invalid or expired OAuth state", nil)
	}

	cfg := s.buildOAuth2Config(p)

	// Exchange code for token
	token, err := cfg.Exchange(c.Request().Context(), input.Code)
	if err != nil {
		return c.AbortBadRequest("Failed to exchange authorization code", err)
	}

	// Fetch user info
	userInfo, err := s.fetchUserInfo(c.Request().Context(), cfg, p.UserInfoURL, token)
	if err != nil {
		return c.AbortInternalServerError("Failed to fetch user info", err)
	}

	oauthID := extractString(userInfo, p.UserIDField)
	email := extractString(userInfo, p.EmailField)
	name := extractString(userInfo, p.NameField)
	avatar := extractString(userInfo, p.AvatarField)

	if oauthID == "" || email == "" {
		return c.AbortBadRequest("OAuth provider did not return required user info (id, email)", nil)
	}

	// Find or create user
	user, isNew, err := s.findOrCreateUser(c, p.Name, oauthID, email, name, avatar)
	if err != nil {
		return c.AbortInternalServerError("Failed to process OAuth user", err)
	}

	if !user.Active {
		return c.AbortForbidden("Account is disabled")
	}

	_ = s.userRepo.UpdateLastLogin(c.Request().Context(), user.ID, c.Request().RemoteAddr)

	// Issue JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"iss":   s.config.JWT.Issuer,
		"aud":   "goma-admin",
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
		"email": user.Email,
		"role":  user.Role,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return c.AbortInternalServerError("Token generation failed", err)
	}

	return c.OK(dto.OAuthLoginResponse{
		AccessToken: tokenString,
		ExpiresAt:   expirationTime.Unix(),
		TokenType:   "Bearer",
		User: dto.UserResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			Name:          user.Name,
			Roles:         user.Role,
			OAuthProvider: user.OAuthProvider,
		},
		NewUser: isNew,
	})
}

//  helpers

func (s *OAuthService) buildOAuth2Config(p *models.OAuthProvider) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  p.AuthURL,
			TokenURL: p.TokenURL,
		},
		RedirectURL: fmt.Sprintf("%s/auth/oauth/callback", s.config.OAuth.BaseURL),
		Scopes:      []string(p.Scopes),
	}
}

func (s *OAuthService) fetchUserInfo(ctx context.Context, cfg *oauth2.Config, userInfoURL string, token *oauth2.Token) (map[string]interface{}, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("userinfo request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read userinfo response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var info map[string]interface{}
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo JSON: %w", err)
	}
	return info, nil
}

func (s *OAuthService) findOrCreateUser(c *okapi.Context, provider, oauthID, email, name, avatar string) (*models.User, bool, error) {
	ctx := c.Request().Context()

	// Try finding by OAuth provider + ID first
	user, err := s.userRepo.GetByOAuth(ctx, provider, oauthID)
	if err == nil {
		return user, false, nil
	}

	// Try finding by email — link the OAuth account
	user, err = s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		user.OAuthProvider = provider
		user.OAuthID = oauthID
		if avatar != "" && user.Avatar == "" {
			user.Avatar = avatar
		}
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, false, err
		}
		return user, false, nil
	}

	// Create new user
	newUser := &models.User{
		Email:         email,
		Name:          name,
		Avatar:        avatar,
		Role:          string(models.RoleViewer),
		EmailVerified: true,
		Active:        true,
		OAuthProvider: provider,
		OAuthID:       oauthID,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, false, err
	}
	return newUser, true, nil
}

func extractString(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		// JSON numbers are float64; format integers without decimals
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toProviderDetail(p *models.OAuthProvider) dto.OAuthProviderDetailResponse {
	return dto.OAuthProviderDetailResponse{
		ID:          p.ID,
		Name:        p.Name,
		DisplayName: p.DisplayName,
		ClientID:    p.ClientID,
		AuthURL:     p.AuthURL,
		TokenURL:    p.TokenURL,
		UserInfoURL: p.UserInfoURL,
		Scopes:      []string(p.Scopes),
		UserIDField: p.UserIDField,
		EmailField:  p.EmailField,
		NameField:   p.NameField,
		AvatarField: p.AvatarField,
		Enabled:     p.Enabled,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
