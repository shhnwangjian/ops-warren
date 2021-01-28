package module

import (
	"net"
	"strconv"
	"time"
)

// ProbeCommon is an interface that defines the Probe function
type ProbeCommon interface {
	Probe() (bool, error)
	RetryProbe() (bool, error)
}

type tcpProbe struct {
	retry   int
	host    string
	port    int
	timeout time.Duration
}

func NewTcpProbe(host string, port int, timeout time.Duration) ProbeCommon {
	return tcpProbe{
		host:    host,
		port:    port,
		timeout: timeout,
		retry:   3,
	}
}

// Probe checks that a TCP socket to the address can be opened.
// If the socket can be opened, it returns Success
// If the socket fails to open, it returns Failure.
// This is exported because some other packages may want to do direct TCP probes.
func (p tcpProbe) Probe() (bool, error) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(p.host, strconv.Itoa(p.port)), p.timeout)
	if err != nil {
		// Convert errors to failures to handle timeouts.
		return false, err
	}
	err = conn.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p tcpProbe) RetryProbe() (status bool, err error) {
	for num := 0; num <= p.retry; num++ {
		status, err = p.Probe()
		if status {
			return status, nil
		}
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
	}
	return false, err
}

func (p tcpProbe) SetRetry(num int) {
	p.retry = num
}
