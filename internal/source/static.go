package source

import (
	"github.com/casantosmu/port-monitor/internal/config"
)

func staticSource(src config.Source) (string, error) {
	return src.Value, nil
}
