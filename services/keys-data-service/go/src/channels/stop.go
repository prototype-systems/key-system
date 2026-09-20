package channels

import "sync"

var (
	stopChannel = make(chan struct{})
	stopEvent   sync.Once
)

func GetStopChannel() chan struct{} {
	return stopChannel
}

func SignalStop() {
	stopEvent.Do(func() {
		close(stopChannel)
	})
}
