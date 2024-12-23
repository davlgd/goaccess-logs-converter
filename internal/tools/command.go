package tools

import (
	"fmt"
	"strings"

	"github.com/davlgd/goaccess-logs-converter/internal/config"
)

func buildCommand(cfg *config.Config) string {
	params := []string{fmt.Sprintf("since=%s", cfg.Since)}
	if cfg.Until != "" {
		params = append(params, fmt.Sprintf("until=%s", cfg.Until))
	}
	if cfg.Limit != 0 {
		params = append(params, fmt.Sprintf("limit=%d", cfg.Limit))
	}
	if cfg.DeploymentID != "" {
		params = append(params, fmt.Sprintf("deploymentId=%s", cfg.DeploymentID))
	}
	if cfg.InstanceID != "" {
		params = append(params, fmt.Sprintf("instanceId=%s", cfg.InstanceID))
	}

	url := fmt.Sprintf(
		"https://api.clever-cloud.com/v4/accesslogs/organisations/%s/applications/%s/accesslogs?throttleElements=42&%s",
		cfg.OrgID, cfg.AppID, strings.Join(params, "&"),
	)
	return fmt.Sprintf(`clever curl "%s"`, url)
}
