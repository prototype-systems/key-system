package channels

import (
	"sync"
	"sync/atomic"
)

var (
	killChannel = make(chan struct{})
	killEvent   sync.Once
	killStatus  atomic.Int32
)

func GetKillChannel() chan struct{} {
	return killChannel
}

func SignalKill() {
	killEvent.Do(func() {
		killStatus.Store(1)
		close(killChannel)
	})
}

func CheckKillStatus() bool {
	return killStatus.Load() == 1
}
