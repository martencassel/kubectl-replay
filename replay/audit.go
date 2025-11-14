package replay

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// StreamAudit reads an audit log file line by line and prints kubectl commands
func StreamAudit(filePath string, speedMultiplier int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event AuditEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			continue
		}
		// Only show ResponseComplete stage to avoid duplicates and get actual response codes
		if event.Stage != "ResponseComplete" {
			continue
		}
		fmt.Println(ToKubectlCommand(event))
		time.Sleep(time.Second / time.Duration(speedMultiplier))
	}
	return scanner.Err()
}

// ToKubectlCommand maps an AuditEvent into a kubectl command with HTTP details
func ToKubectlCommand(event AuditEvent) string {
	// Generate kubectl command
	verb := verbToKubectl(event.Verb)
	cmd := fmt.Sprintf("kubectl %s %s", verb, event.ObjectRef.Resource)
	if event.ObjectRef.Name != "" {
		cmd += " " + event.ObjectRef.Name
	}
	if event.ObjectRef.Namespace != "" {
		cmd += " -n " + event.ObjectRef.Namespace
	}

	// HTTP method and status
	method := verbToHTTPMethod(event.Verb)
	statusCode := event.ResponseStatus.Code
	statusText := event.ResponseStatus.Reason
	if event.ResponseStatus.Message != "" {
		statusText = event.ResponseStatus.Message
	}

	// Parse and format timestamp (just time portion)
	timeStr := "--:--:--"
	if event.StageTimestamp != "" {
		if len(event.StageTimestamp) >= 19 {
			// Extract HH:MM:SS from ISO8601 timestamp
			timeStr = event.StageTimestamp[11:19]
		}
	}

	// Color code based on status
	var statusColor, cmdColor string
	switch {
	case statusCode >= 200 && statusCode < 300:
		statusColor = "\033[32m" // Green
		cmdColor = "\033[90m"    // Gray for successful commands
	case statusCode == 403:
		statusColor = "\033[33m" // Yellow
		cmdColor = "\033[33m"
	case statusCode == 404:
		statusColor = "\033[35m" // Magenta
		cmdColor = "\033[35m"
	case statusCode >= 400:
		statusColor = "\033[31m" // Red
		cmdColor = "\033[31m"
	default:
		statusColor = "\033[90m" // Gray
		cmdColor = "\033[90m"
	}
	dimColor := "\033[90m"
	reset := "\033[0m"

	// Format: metadata line (dimmed) + command line (colored by status)
	return fmt.Sprintf("%s# [%s] %s %s → %s%d %s%s\n%s%s%s",
		dimColor,
		timeStr,
		method,
		event.RequestURI,
		statusColor,
		statusCode,
		statusText,
		reset,
		cmdColor,
		cmd,
		reset,
	)
}

func verbToHTTPMethod(verb string) string {
	switch verb {
	case "get", "list", "watch":
		return "GET"
	case "create":
		return "POST"
	case "update":
		return "PUT"
	case "patch":
		return "PATCH"
	case "delete", "deletecollection":
		return "DELETE"
	default:
		return verb
	}
}

func verbToKubectl(verb string) string {
	switch verb {
	case "list":
		return "get"
	case "watch":
		return "get --watch"
	default:
		return verb
	}
}
