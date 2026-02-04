package config

// TargetConfig represents a phishing target configuration
type TargetConfig struct {
	Name                string
	TargetDomain        string
	PhishingDomain      string
	Description         string
	ExternalOrigins     []string
	ContentReplacements []ContentReplacement
	LoginPaths          []string
	CredentialPatterns  []CredentialPattern
	AuthSessionURLs     []string
	TriggerCookies      []string
	LogPath             string
	SSLCertPath         string
	SSLKeyPath          string
	NecroBrowserConfig  string
}

// ContentReplacement represents a content replacement rule
type ContentReplacement struct {
	From string
	To   string
}

// CredentialPattern represents a credential capture pattern
type CredentialPattern struct {
	Label string
	Start string
	End   string
}

// MuraenaConfig represents the complete Muraena configuration
type MuraenaConfig struct {
	Proxy     ProxyConfig     `toml:"proxy"`
	TLS       TLSConfig       `toml:"tls"`
	Log       LogConfig       `toml:"log"`
	Origins   OriginsConfig   `toml:"origins"`
	Transform TransformConfig `toml:"transform"`
	Tracking  TrackingConfig  `toml:"tracking"`
	Crawler   CrawlerConfig   `toml:"crawler"`
	Redis     RedisConfig     `toml:"redis"`
	Necro     NecroConfig     `toml:"necrobrowser"`
}

// ProxyConfig represents proxy configuration
type ProxyConfig struct {
	Phishing    string          `toml:"phishing"`
	Destination string          `toml:"destination"`
	IP          string          `toml:"IP"`
	Listener    string          `toml:"listener"`
	HTTPtoHTTPS HTTPtoHTTPSConf `toml:"HTTPtoHTTPS"`
}

// HTTPtoHTTPSConf represents HTTP to HTTPS redirect config
type HTTPtoHTTPSConf struct {
	Enable   bool `toml:"enable"`
	HTTPPort int  `toml:"HTTPport"`
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	Enable      bool   `toml:"enable"`
	Expand      bool   `toml:"expand"`
	Certificate string `toml:"certificate"`
	Key         string `toml:"key"`
	Root        string `toml:"root"`
	MinVersion  string `toml:"minVersion"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	Enable   bool   `toml:"enable"`
	FilePath string `toml:"filePath"`
}

// OriginsConfig represents external origins configuration
type OriginsConfig struct {
	ExternalOriginPrefix string   `toml:"externalOriginPrefix"`
	ExternalOrigins      []string `toml:"externalOrigins"`
}

// TransformConfig represents transformation configuration
type TransformConfig struct {
	Base64   Base64Config   `toml:"base64"`
	Request  RequestConfig  `toml:"request"`
	Response ResponseConfig `toml:"response"`
}

// Base64Config represents base64 encoding configuration
type Base64Config struct {
	Enable bool `toml:"enable"`
}

// RequestConfig represents request transformation configuration
type RequestConfig struct {
	Headers []string `toml:"headers"`
}

// ResponseConfig represents response transformation configuration
type ResponseConfig struct {
	SkipContentType []string             `toml:"skipContentType"`
	Headers         []string             `toml:"headers"`
	CustomContent   [][]string           `toml:"customContent"`
	Remove          ResponseRemoveConfig `toml:"remove"`
}

// ResponseRemoveConfig represents headers to remove
type ResponseRemoveConfig struct {
	Headers []string `toml:"headers"`
}

// TrackingConfig represents tracking configuration
type TrackingConfig struct {
	Enable  bool          `toml:"enable"`
	Trace   TraceConfig   `toml:"trace"`
	Secrets SecretsConfig `toml:"secrets"`
}

// TraceConfig represents trace configuration
type TraceConfig struct {
	Identifier string        `toml:"identifier"`
	Validator  string        `toml:"validator"`
	Landing    LandingConfig `toml:"landing"`
}

// LandingConfig represents landing page configuration
type LandingConfig struct {
	Type string `toml:"type"`
}

// SecretsConfig represents secrets tracking configuration
type SecretsConfig struct {
	Paths    []string        `toml:"paths"`
	Patterns []PatternConfig `toml:"patterns"`
}

// PatternConfig represents a credential pattern
type PatternConfig struct {
	Label string `toml:"label"`
	Start string `toml:"start"`
	End   string `toml:"end"`
}

// CrawlerConfig represents crawler configuration
type CrawlerConfig struct {
	Enable bool `toml:"enable"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Enable bool   `toml:"enable"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
}

// NecroConfig represents NecroBrowser configuration
type NecroConfig struct {
	Enable   bool        `toml:"enable"`
	Endpoint string      `toml:"endpoint"`
	Profile  string      `toml:"profile"`
	URLs     URLsConfig  `toml:"urls"`
	Trigger  TriggerConf `toml:"trigger"`
}

// URLsConfig represents URL configuration
type URLsConfig struct {
	AuthSession []string `toml:"authSession"`
}

// TriggerConf represents trigger configuration
type TriggerConf struct {
	Type   string   `toml:"type"`
	Values []string `toml:"values"`
	Delay  int      `toml:"delay"`
}
