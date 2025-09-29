package source

import (
	"fmt"
	"os"
	"regexp"

	"github.com/casantosmu/port-monitor/internal/config"
)

func fileSource(src config.Source) (string, error) {
	content, err := os.ReadFile(src.Path)
	if err != nil {
		return "", err
	}

	if src.Pattern == "" {
		return string(content), nil
	}

	re, err := regexp.Compile(src.Pattern)
	if err != nil {
		return "", fmt.Errorf("invalid pattern: %w", err)
	}

	matches := re.FindSubmatch(content)
	if len(matches) == 0 {
		return "", fmt.Errorf("pattern '%s' not found in file '%s'", src.Pattern, src.Path)
	}

	if src.MatchGroup == nil {
		return "", fmt.Errorf("match_group is required")
	}

	if *src.MatchGroup >= len(matches) {
		return "", fmt.Errorf("match_group %d out of range (only %d groups captured)", *src.MatchGroup, len(matches)-1)
	}

	match := matches[*src.MatchGroup]
	return string(match), nil
}
