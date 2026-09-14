package channels

import (
	"sync"
	"sync/atomic"
)

var (
	abortChannel = make(chan struct{})
	abortEvent   sync.Once
	abortStatus  atomic.Int32
)

func GetAbortChannel() chan struct{} {
	return abortChannel
}

func SignalAbort() {
	abortEvent.Do(func() {
		abortStatus.Store(1)
		close(abortChannel)
	})
}

func CheckAbortStatus() bool {
	return abortStatus.Load() == 1
}
