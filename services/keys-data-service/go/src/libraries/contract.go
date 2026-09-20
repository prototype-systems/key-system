package libraries

import (
	"errors"
	"sync"
)

type Contract[Value any, Breach any] struct {
	control sync.RWMutex

	setups []func() error

	breach *Signal[Breach]

	pipes []chan Value
}

var (
	IssueContractResolved = errors.New("contract resolved")

	IssueContractBreached = errors.New("contract breached")

	IssueContractInitiated = errors.New("contract initiated")

	IssueInvalidContractValue = errors.New("invalid contract value")
)

func ConstructContract[Value any, Breach any]() *Contract[Value, Breach] {
	return &Contract[Value, Breach]{
		breach: ConstructSignal[Breach](),
		setups: []func() error{},
		pipes:  []chan Value{},
	}
}

func (contract *Contract[Value, Breach]) Setup(initiator func() error) error {
	contract.control.Lock()
	defer contract.control.Unlock()

	if contract.breach == nil {
		return IssueContractBreached
	}

	if contract.pipes == nil {
		return IssueContractResolved
	}

	if contract.setups == nil {
		return IssueContractInitiated
	}

	contract.setups = append(contract.setups, initiator)

	return nil
}

func (contract *Contract[Value, Breach]) Trigger(condition func(value Value, breach func(issue Breach)) <-chan any) (<-chan any, error) {
	contract.control.Lock()

	if contract.breach == nil {
		contract.control.Unlock()
		return nil, IssueContractBreached
	}

	if contract.pipes == nil {
		contract.control.Unlock()
		return nil, IssueContractResolved
	}

	if contract.setups == nil {
		contract.control.Unlock()
		return nil, IssueContractInitiated
	}

	pipe := make(chan Value, 1)
	contract.pipes = append(contract.pipes, pipe)

	breach := contract.breach

	contract.control.Unlock()

	done := make(chan any, 1)

	go func() {
		value := <-pipe

		resultPipe := condition(value, func(issue Breach) {
			contract.control.Lock()

			if contract.breach != nil {
				contract.breach = nil
				contract.control.Unlock()

				breach.Push(issue)
			} else {
				contract.control.Unlock()
			}
		})

		if resultPipe == nil {
			panic("non-responsive condition on contract breach")
		}

		result := <-resultPipe
		done <- result
	}()

	return done, nil
}

func (contract *Contract[Value, Breach]) Check(handler func(issue Breach) any) (<-chan any, error) {
	contract.control.RLock()

	if contract.breach == nil {
		contract.control.RUnlock()
		return nil, IssueContractBreached
	}

	result, issue := contract.breach.Listen(handler)

	contract.control.RUnlock()

	return result, issue
}

func (contract *Contract[Value, Breach]) Initiate() error {
	contract.control.Lock()

	if contract.breach == nil {
		contract.control.Unlock()
		return IssueContractBreached
	}

	if contract.pipes == nil {
		contract.control.Unlock()
		return IssueContractResolved
	}

	if contract.setups == nil {
		contract.control.Unlock()
		return IssueContractInitiated
	}

	setups := contract.setups
	contract.setups = nil

	contract.control.Unlock()

	for _, initiator := range setups {
		if issue := initiator(); issue != nil {
			panic(issue)
		}
	}

	return nil
}

func (contract *Contract[Value, Breach]) Push(value Value) error {
	if any(value) == nil {
		return IssueInvalidContractValue
	}

	contract.control.Lock()

	if contract.breach == nil {
		contract.control.Unlock()
		return IssueContractBreached
	}

	if contract.pipes == nil {
		contract.control.Unlock()
		return IssueContractResolved
	}

	pipes := contract.pipes
	contract.pipes = nil

	contract.control.Unlock()

	for _, pipe := range pipes {
		pipe <- value
		close(pipe)
	}

	return nil
}

func (contract *Contract[Value, Breach]) Data() (<-chan Value, error) {
	contract.control.Lock()
	defer contract.control.Unlock()

	if contract.breach == nil {
		return nil, IssueContractBreached
	}

	if contract.pipes == nil {
		return nil, IssueContractResolved
	}

	pipe := make(chan Value, 1)

	contract.pipes = append(contract.pipes, pipe)

	return pipe, nil
}
