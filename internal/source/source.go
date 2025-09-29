package source

import (
	"fmt"

	"github.com/casantosmu/port-monitor/internal/config"
)

func Get(src config.Source) (string, error) {
	switch src.Type {
	case config.SourceTypeFile:
		return fileSource(src)
	default:
		return "", fmt.Errorf("unknown source type: %s", src.Type)
	}
}
