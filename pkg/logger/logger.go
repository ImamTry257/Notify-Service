package logger

import (
	"encoding/json"
	"fmt"
	"log"
)

// Info logs a general informational message with optional key-value fields.
// Usage: logger.Info("FlowName", "message", "key1", "val1", "key2", "val2")
func Info(flow, message string, fields ...any) {
	log.Printf("[%s] %s%s", flow, message, buildFields(fields...))
}

// Request logs the incoming request payload for a given flow.
func Request(flow string, data any) {
	b, _ := json.Marshal(data)
	log.Printf("[%s] REQUEST: %s", flow, string(b))
}

// Response logs the outgoing response payload for a given flow.
func Response(flow string, data any) {
	b, _ := json.Marshal(data)
	log.Printf("[%s] RESPONSE: %s", flow, string(b))
}

// Error logs an error message with optional key-value context fields.
// Usage: logger.Error("FlowName", "description", err, "email", "x@y.com")
func Error(flow, message string, err error, fields ...any) {
	log.Printf("[%s] ERROR %s: %v%s", flow, message, err, buildFields(fields...))
}

// buildFields formats key-value pairs into " | key: value" format.
// Expects pairs: key1, val1, key2, val2, ...
func buildFields(fields ...any) string {
	if len(fields) == 0 {
		return ""
	}
	result := ""
	for i := 0; i+1 < len(fields); i += 2 {
		result += fmt.Sprintf(" | %v: %v", fields[i], fields[i+1])
	}
	// handle odd trailing key without value
	if len(fields)%2 != 0 {
		result += fmt.Sprintf(" | %v: <missing>", fields[len(fields)-1])
	}
	return result
}
