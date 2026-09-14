//go:build windows

package libraries

import (
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

func init() {
	listenNative = func(address string) (net.Listener, error) {
		return winio.ListenPipe(address, nil)
	}

	dialNative = func(address string, timeout time.Duration) (net.Conn, error) {
		if timeout > 0 {
			return winio.DialPipe(address, &timeout)
		}

		return winio.DialPipe(address, nil)
	}
}
