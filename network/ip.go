package network

import (
	"net"
	"strings"
)

func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}
