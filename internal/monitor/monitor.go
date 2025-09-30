package monitor

import (
	"fmt"

	"github.com/casantosmu/port-monitor/internal/config"
	"github.com/casantosmu/port-monitor/internal/source"
)

type Result struct {
	IP   string
	Port string
}

func Start(svc config.Service) (Result, error) {
	res := Result{}

	ip, err := source.Get(svc.IPSource)
	if err != nil {
		return res, fmt.Errorf("ip_source failed: %w", err)
	}

	port, err := source.Get(svc.PortSource)
	if err != nil {
		return res, fmt.Errorf("port_source failed: %w", err)
	}
	res.IP = ip
	res.Port = port

	return res, nil
}
