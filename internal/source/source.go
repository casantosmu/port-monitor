package source

import (
	"context"
	"fmt"

	"github.com/casantosmu/port-monitor/internal/config"
)

func Get(ctx context.Context, src config.Source) (string, error) {
	switch src.Type {
	case config.SourceTypeHTTP:
		return httpSource(ctx, src)
	case config.SourceTypeFile:
		return fileSource(src)
	case config.SourceTypeStatic:
		return staticSource(src)
	default:
		return "", fmt.Errorf("unknown source type: %s", src.Type)
	}
}
