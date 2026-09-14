package channels

import "sync"

var (
	startChannel = make(chan struct{})
	startEvent   sync.Once
)

func GetStartChannel() chan struct{} {
	return startChannel
}

func SignalStart() {
	startEvent.Do(func() {
		close(startChannel)
	})
}
