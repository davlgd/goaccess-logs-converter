package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	AppID        string
	DeploymentID string
	InstanceID   string
	Limit        uint
	OrgID        string
	Output       string
	Since        string
	Until        string
}

var (
	cfg         Config
	showVersion bool
)

func InitFlags() {
	flag.StringVar(&cfg.AppID, "app", "", "Application ID")
	flag.StringVar(&cfg.DeploymentID, "deployment", "", "Filter by deployment ID")
	flag.StringVar(&cfg.InstanceID, "instance", "", "Filter by instance ID")
	flag.UintVar(&cfg.Limit, "limit", 0, "Limit number of results (min: 1)")
	flag.StringVar(&cfg.OrgID, "org", "", "Organization ID")
	flag.StringVar(&cfg.Output, "out", "goaccess_logs.txt", "Output file name")
	flag.StringVar(&cfg.Since, "since", "", "Start date (ISO8601 format)")
	flag.StringVar(&cfg.Until, "until", "", "End date (ISO8601 format)")
	flag.BoolVar(&showVersion, "v", false, "Show version")
}

func Usage(name, version, description string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, `%s (v%s)
%s

Required:
  -org string        Organization ID
  -app string        Application ID
  -since string      Start date (ISO 8601 format)

Optional:

  -deployment string Filter by deployment ID
  -instance string   Filter by instance ID
  -limit int         Limit number of results (min: 1)
  -until string      End date (ISO 8601 format)

  -out string        Output file name (default "goaccess_logs.txt")

Examples:
  # Basic usage
  %s -org=orga_xxx -app=app_xxx -since=2024-12-01T00:00:00Z

  # With until, limit
  %s -org=orga_xxx -app=app_xxx -since=2024-12-05T00:00:00Z -until=2024-12-02T00:00:00Z -limit=1000
`,
			name, version, description,
			filepath.Base(os.Args[0]),
			filepath.Base(os.Args[0]))
	}
}

func New() (*Config, error) {

	if showVersion {
		return &cfg, nil
	}

	if cfg.OrgID == "" || cfg.AppID == "" || cfg.Since == "" {
		return nil, fmt.Errorf("missing required parameters")
	}

	if _, err := time.Parse(time.RFC3339, cfg.Since); err != nil {
		return nil, fmt.Errorf("invalid since date format, must be ISO 8601")
	}

	if cfg.Until != "" {
		if _, err := time.Parse(time.RFC3339, cfg.Until); err != nil {
			return nil, fmt.Errorf("invalid until date format, must be ISO 8601")
		}
	}

	if cfg.Limit != 0 && cfg.Limit < 1 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	return &cfg, nil
}

func ShouldCheckHeartbeats(cfg *Config) bool {
	return cfg.Until != "" || cfg.Limit != 0
}

func ShouldShowVersion() bool {
	return showVersion
}
