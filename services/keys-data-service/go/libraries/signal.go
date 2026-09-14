package libraries

import (
	"errors"
	"sync"
)

type Signal[Value any] struct {
	control sync.RWMutex
	pipes   []chan Value
}

var (
	IssueResolvedSignal = errors.New("signal resolved")

	IssueInvalidSignalValue = errors.New("invalid signal value")
)

func ConstructSignal[Value any]() *Signal[Value] {
	return &Signal[Value]{
		pipes: []chan Value{},
	}
}

func (signal *Signal[Value]) Push(value Value) error {
	if any(value) == nil {
		return IssueInvalidSignalValue
	}

	signal.control.Lock()

	if signal.pipes == nil {
		signal.control.Unlock()
		return IssueResolvedSignal
	}

	pipes := signal.pipes
	signal.pipes = nil

	signal.control.Unlock()

	for _, pipe := range pipes {
		pipe <- value
		close(pipe)
	}

	return nil
}

func (signal *Signal[Value]) Listen(handler func(Value) any) (<-chan any, error) {
	dataPipe, issue := signal.Data()
	if issue != nil {
		return nil, issue
	}

	resultPipe := make(chan any, 1)

	go func() {
		resultPipe <- handler(<-dataPipe)
	}()

	return resultPipe, nil
}

func (signal *Signal[Value]) Data() (<-chan Value, error) {
	signal.control.Lock()
	defer signal.control.Unlock()

	if signal.pipes == nil {
		return nil, IssueResolvedSignal
	}

	pipe := make(chan Value, 1)

	signal.pipes = append(signal.pipes, pipe)

	return pipe, nil
}

func (signal *Signal[Value]) Status() bool {
	signal.control.RLock()
	defer signal.control.RUnlock()

	return signal.pipes == nil
}
