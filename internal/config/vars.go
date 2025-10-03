package config

import (
	"os"
	"regexp"
)

var (
	placeholderPattern = regexp.MustCompile(`{{\s*(.+?)\s*}}`)

	reservedPlaceholders = map[string]bool{
		"ip":   true,
		"port": true,
	}
)

func expandVars(content []byte) []byte {
	return placeholderPattern.ReplaceAllFunc(content, func(match []byte) []byte {
		matches := placeholderPattern.FindSubmatch(match)

		if len(matches) < 2 {
			return match
		}

		key := string(matches[1])
		if reservedPlaceholders[key] {
			return match
		}

		val := os.Getenv(key)
		if val == "" {
			return match
		}

		return []byte(val)
	})
}
