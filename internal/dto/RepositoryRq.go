package dto

type CreateRepositoryRq struct {
	Body struct {
		Name      string `json:"name" required:"true" minLength:"2" description:"Repository name" example:"gateway-configs"`
		URL       string `json:"url" required:"true" description:"Git clone URL" example:"https://github.com/org/configs.git"`
		Branch    string `json:"branch,omitempty" description:"Branch to track" example:"main"`
		AuthType  string `json:"authType,omitempty" enum:"token,basic,ssh" description:"Authentication type"`
		AuthValue string `json:"authValue,omitempty" description:"Auth credentials (PAT, user:pass, or SSH key)"`
	} `json:"body"`
}

type UpdateRepositoryRq struct {
	ID   int `param:"id" required:"true" description:"Repository ID"`
	Body struct {
		Name      string `json:"name" required:"true" minLength:"2" description:"Repository name"`
		URL       string `json:"url" required:"true" description:"Git clone URL"`
		Branch    string `json:"branch,omitempty" description:"Branch to track"`
		AuthType  string `json:"authType,omitempty" enum:"token,basic,ssh" description:"Authentication type"`
		AuthValue string `json:"authValue,omitempty" description:"Auth credentials (leave empty to keep existing)"`
	} `json:"body"`
}

type RepositoryByIDRq struct {
	ID int `param:"id" required:"true" description:"Repository ID"`
}

type BrowseEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Path  string `json:"path"`
}

type SyncResult struct {
	Status string `json:"status"`
	Commit string `json:"commit,omitempty"`
}
