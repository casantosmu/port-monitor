package main

import (
	"net"
	"time"
)

func isOpen(ip string, port string) bool {
	address := net.JoinHostPort(ip, port)

	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}
