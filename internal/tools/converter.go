package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/davlgd/goaccess-logs-converter/internal/config"
	"github.com/davlgd/goaccess-logs-converter/internal/types"
)

const MAX_HEARTBEATS = 20

func processLogEntry(entry types.LogEntry) string {
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

func Process(cfg *config.Config) error {
	cmd := exec.Command("bash", "-c", buildCommand(cfg))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating pipe: %v", err)
	}

	file, err := os.Create(cfg.Output)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	heartbeatCount := 0
	checkHeartbeats := config.ShouldCheckHeartbeats(cfg)
	var currentBlock []string

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if len(currentBlock) >= 2 && strings.HasPrefix(currentBlock[1], "event:") {
				event := strings.TrimPrefix(currentBlock[1], "event:")

				if event == "HEARTBEAT" {
					heartbeatCount++
					if checkHeartbeats && heartbeatCount >= MAX_HEARTBEATS {
						return fmt.Errorf("too many consecutive heartbeats, stopping")
					}
				} else if event == "ACCESS_LOG" {
					heartbeatCount = 0

					var entry types.LogEntry
					if err := json.Unmarshal([]byte(strings.TrimPrefix(currentBlock[0], "data:")), &entry); err != nil {
						fmt.Printf("JSON parsing error: %v\nProblematic JSON: %s\n", err, currentBlock[0])
						currentBlock = nil
						continue
					}

					if logLine := processLogEntry(entry); logLine != "" {
						fmt.Println(logLine)
						if _, err := file.WriteString(logLine + "\n"); err != nil {
							return fmt.Errorf("error writing to file: %v", err)
						}
						file.Sync()
					}
				}
			}
			currentBlock = nil
			continue
		}
		currentBlock = append(currentBlock, line)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command execution error: %v", err)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stream reading error: %v", err)
	}

	return nil
}
