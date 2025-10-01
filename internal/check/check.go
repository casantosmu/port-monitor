package check

import (
	"context"
	"fmt"

	"github.com/casantosmu/port-monitor/internal/config"
)

func Verify(ctx context.Context, ip, port string, check config.Check) (bool, error) {
	switch check.Type {
	case config.CheckTypeHTTP:
		return httpCheck(ctx, ip, port, check)
	default:
		return false, fmt.Errorf("unknown check type: %s", check.Type)
	}
}
