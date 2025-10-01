package monitor

import (
	"context"
	"fmt"

	"github.com/casantosmu/port-monitor/internal/check"
	"github.com/casantosmu/port-monitor/internal/config"
	"github.com/casantosmu/port-monitor/internal/source"
)

type Result struct {
	IP      string
	Port    string
	Success bool
}

func Start(ctx context.Context, svc config.Service) (Result, error) {
	res := Result{}

	ip, err := source.Get(ctx, svc.IPSource)
	if err != nil {
		return res, fmt.Errorf("ip_source failed: %w", err)
	}

	port, err := source.Get(ctx, svc.PortSource)
	if err != nil {
		return res, fmt.Errorf("port_source failed: %w", err)
	}

	ok, err := check.Verify(ctx, ip, port, svc.ConnectivityCheck)
	if err != nil {
		return res, fmt.Errorf("check failed: %w", err)
	}

	res.IP = ip
	res.Port = port
	res.Success = ok

	return res, nil
}
