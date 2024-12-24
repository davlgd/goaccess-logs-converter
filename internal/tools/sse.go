package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/davlgd/goaccess-logs-converter/internal/config"
	"github.com/davlgd/goaccess-logs-converter/internal/types"
)

func ParseSSE(reader io.Reader, messageChan chan types.SSEMessage) {
	scanner := bufio.NewScanner(reader)
	var message types.SSEMessage

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if message.Data != "" || message.Event != "" {
				messageChan <- message
				message = types.SSEMessage{}
			}
			continue
		}

		switch {
		case strings.HasPrefix(line, "data:"):
			message.Data = strings.TrimPrefix(line, "data:")
		case strings.HasPrefix(line, "event:"):
			message.Event = strings.TrimPrefix(line, "event:")
		case strings.HasPrefix(line, "id:"):
			message.ID = strings.TrimPrefix(line, "id:")
		}
	}
}

func ProcessSSE(cfg *config.Config) error {
	cmd := exec.Command("bash", "-c", buildCommand(cfg))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating pipe: %v", err)
	}

	messageChan := make(chan types.SSEMessage)
	errorChan := make(chan error)

	go func() {
		ParseSSE(stdout, messageChan)
		close(messageChan)
	}()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	file, err := os.Create(cfg.Output)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	heartbeatCount := 0
	checkHeartbeats := config.ShouldCheckHeartbeats(cfg)
	heartbeatTimer := time.NewTimer(30 * time.Second)
	defer heartbeatTimer.Stop()

	for {
		select {
		case message, ok := <-messageChan:
			if !ok {
				return nil // Channel closed, processing complete
			}
			heartbeatTimer.Reset(30 * time.Second)

			switch message.Event {
			case "HEARTBEAT":
				heartbeatCount++
				if checkHeartbeats && heartbeatCount >= MAX_HEARTBEATS {
					return fmt.Errorf("too many consecutive heartbeats, stopping")
				}
			case "ACCESS_LOG":
				heartbeatCount = 0
				if err := processAccessLog(message.Data, file); err != nil {
					return err
				}
			}

		case <-heartbeatTimer.C:
			return fmt.Errorf("no messages received for 30 seconds, connection may be lost")

		case err := <-errorChan:
			return fmt.Errorf("stream error: %v", err)
		}
	}
}

func processAccessLog(data string, file *os.File) error {
	var entry types.LogEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		fmt.Printf("JSON parsing error: %v\nProblematic JSON: %s\n", err, data)
		return nil // Continue processing other logs
	}

	if logLine := convertLogEntry(entry); logLine != "" {
		fmt.Println(logLine)
		if _, err := file.WriteString(logLine + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
		file.Sync()
	}
	return nil
}
