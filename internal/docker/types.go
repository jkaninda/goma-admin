package docker

type GomaConfig struct {
	Routes []Route `json:"routes" yaml:"routes"`
}

type Route struct {
	Name           string           `yaml:"name" json:"name"`
	Path           string           `yaml:"path" json:"path"`
	Rewrite        string           `yaml:"rewrite,omitempty" json:"rewrite,omitempty"`
	Priority       int              `yaml:"priority,omitempty" json:"priority,omitempty"`
	Enabled        bool             `yaml:"enabled,omitempty" default:"true" json:"enabled,omitempty"`
	Hosts          []string         `yaml:"hosts,omitempty" json:"hosts,omitempty"`
	Methods        []string         `yaml:"methods,omitempty" json:"methods,omitempty"`
	Target         string           `yaml:"target,omitempty" json:"target,omitempty"`
	HealthCheck    RouteHealthCheck `yaml:"healthCheck,omitempty" json:"healthCheck,omitempty"`
	Security       Security         `yaml:"security,omitempty" json:"security,omitempty"`
	DisableMetrics bool             `yaml:"disableMetrics,omitempty" json:"disableMetrics,omitempty"`
	Middlewares    []string         `yaml:"middlewares,omitempty" json:"middlewares,omitempty"`
}

type RouteHealthCheck struct {
	Path            string `yaml:"path,omitempty" json:"path,omitempty"`
	Interval        string `yaml:"interval,omitempty" json:"interval,omitempty"`
	Timeout         string `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	HealthyStatuses []int  `yaml:"healthyStatuses,omitempty" json:"healthyStatuses,omitempty"`
}

type Security struct {
	ForwardHostHeaders      bool        `yaml:"forwardHostHeaders" json:"forwardHostHeaders" default:"true"`
	EnableExploitProtection bool        `yaml:"enableExploitProtection" json:"enableExploitProtection"`
	TLS                     SecurityTLS `yaml:"tls" json:"tls"`
}

type SecurityTLS struct {
	InsecureSkipVerify bool `yaml:"insecureSkipVerify,omitempty" json:"insecureSkipVerify,omitempty"`
}
