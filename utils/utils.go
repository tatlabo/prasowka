package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GenerateID creates a unique identifier.
func GenerateID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func ConvertBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func CleanInput(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	re := regexp.MustCompile(`[^a-zA-Z0-9.,;\/-_!$ ]+ `)
	// re := regexp.MustCompile(`[^a-zA-Z0-9.,\-_ ]+`)
	return re.ReplaceAllString(input, "")
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05") // Customize the format as needed
}

func Not(b bool) bool {
	return !b
}

func Equals(a, b any) bool {
	return a == b
}

func Notequals(a, b any) bool {
	return a != b
}
