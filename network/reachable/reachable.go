package reachable

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"time"

	iurl "github.com/wuxler/goutils/network/url"
)

const (
	// ReachableDefaultPort is the default port used if no port is defined in a reachable checker
	ReachableDefaultPort = "80"
	// ReachableDefaultNetwork is the default network used in the reachable checker
	ReachableDefaultNetwork = "tcp"
	// ReachableDefaultTimeout is the default timeout used when reachable is checking the URL
	ReachableDefaultTimeout = time.Duration(3) * time.Second
)

type ReachableDialer func(network, address string, timeout time.Duration) (net.Conn, error)

type ReachableConfig struct {
	// Dialer is optional and defaults to use net.DialTimeout.
	Dialer ReachableDialer
	// Timeout is optional and defaults to "3s".
	Timeout time.Duration
	// Network is optional and defaults to "tcp". It should be one of "tcp", "tcp4", "tcp6",
	// "unix", "unixpacket", "udp", "udp4", "udp6", "unixgram" or an IP transport. The IP
	// transports are "ip", "ip4", or "ip6" followed by a colon and a literal protocol number
	// or a protocol name, as in "ip:1" or "ip:icmp".
	Network string
}

type ReachableChecker struct {
	dialer  ReachableDialer
	timeout time.Duration
	network string
}

func New(cfg *ReachableConfig) *ReachableChecker {
	if cfg == nil {
		cfg = &ReachableConfig{}
	}
	timeout := ReachableDefaultTimeout
	if cfg.Timeout != 0 {
		timeout = cfg.Timeout
	}
	dialer := net.DialTimeout
	if cfg.Dialer != nil {
		dialer = cfg.Dialer
	}
	network := ReachableDefaultNetwork
	if cfg.Network != "" {
		network = cfg.Network
	}

	return &ReachableChecker{
		dialer:  dialer,
		timeout: timeout,
		network: network,
	}
}

func (r *ReachableChecker) Check(addr string) error {
	address, err := r.resolve(addr)
	if err != nil {
		return err
	}

	conn, err := r.dialer(r.network, address, r.timeout)
	if err != nil {
		return err
	}
	if conn != nil {
		if errClose := conn.Close(); errClose != nil {
			return err
		}
	}

	return nil
}

func (r *ReachableChecker) resolve(addr string) (string, error) {
	if len(addr) == 0 {
		return "", errors.New("empty addr")
	}

	if strings.Index(addr, "://") < 0 {
		addr = "dump://" + addr
	}

	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	// We must provide a port so when a port is not set in the URL provided use
	// the default port (80)
	port := iurl.Port(u)
	if len(port) == 0 {
		port = ReachableDefaultPort
	}
	return net.JoinHostPort(iurl.Hostname(u), port), nil
}

func Check(addr string) error {
	checker := New(nil)
	return checker.Check(addr)
}
