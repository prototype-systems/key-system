//go:build !windows

package libraries

import (
	"net"
	"time"
)

func init() {
	listenNative = func(address string) (net.Listener, error) {
		return net.Listen("unix", address)
	}

	dialNative = func(address string, timeout time.Duration) (net.Conn, error) {
		if timeout > 0 {
			return net.DialTimeout("unix", address, timeout)
		}

		return net.Dial("unix", address)
	}
}
