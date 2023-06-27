package url

import (
	stdurl "net/url"

	"github.com/wuxler/goutils/network"
)

// Hostname returns url.Host when url.Host is IPv6, otherwise returns url.Hostname().
func Hostname(url *stdurl.URL) string {
	if network.IsIPv6(url.Host) {
		return url.Host
	}
	return url.Hostname()
}

// Port returns "" when url.Host is IPv6, otherwise returns url.Port().
func Port(url *stdurl.URL) string {
	if network.IsIPv6(url.Host) {
		return ""
	}
	return url.Port()
}
