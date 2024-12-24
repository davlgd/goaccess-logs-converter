package tools

import (
	"fmt"
	"time"

	"github.com/davlgd/goaccess-logs-converter/internal/types"
)

const MAX_HEARTBEATS = 20

func convertLogEntry(entry types.LogEntry) string {
	t, err := time.Parse(time.RFC3339, entry.Date)
	if err != nil {
		fmt.Printf("Date parsing error: %v\n", err)
		return ""
	}

	method, path, status := "-", "-", 0
	if entry.HTTP != nil {
		method = entry.HTTP.Request.Method
		path = entry.HTTP.Request.Path
		status = entry.HTTP.Response.StatusCode
	}
	return fmt.Sprintf("%s - - [%s] \"%s %s\" %d %d",
		entry.Source.IP,
		t.Format("02/Jan/2006:15:04:05 -0700"),
		method, path, status,
		entry.BytesOut)
}
