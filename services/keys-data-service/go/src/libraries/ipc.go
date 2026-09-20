package libraries

import (
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type IpcSetting struct {
	Address string
	Timeout time.Duration
	Invoke  func(data string) error
}

type Ipc struct {
	Listener    net.Listener
	Connections IpcConnections
	Setting     IpcSetting
	Control     *sync.Mutex
}

type IpcConnections map[net.Conn]bool

var listenNative func(address string) (net.Listener, error)
var dialNative func(address string, timeout time.Duration) (net.Conn, error)

var (
	ipcService *Ipc
	ipcControl sync.RWMutex
)

func SetIpc(ipc *Ipc) {
	ipcControl.Lock()
	defer ipcControl.Unlock()

	ipcService = ipc
}

func GetIpc() *Ipc {
	ipcControl.RLock()
	defer ipcControl.RUnlock()

	return ipcService
}

func SetupIpc(setting IpcSetting) *Contract[*Ipc, string] {
	ipcContract := ConstructContract[*Ipc, string]()

	trigger := func(condition func(ipc *Ipc, breach func(issue string)) <-chan any) {
		if _, issue := ipcContract.Trigger(condition); issue != nil {
			panic("setup ipc contract on trigger with unrelated issue, " + issue.Error())
		}
	}

	ipcDone := make(chan bool)

	if issue := ipcContract.Setup(func() error {
		control := &sync.Mutex{}
		connections := make(IpcConnections)

		listener, issue := listenNative(setting.Address)
		if issue != nil {
			return issue
		}

		go func() {
			defer close(ipcDone)

			for {
				connection, issue := listener.Accept()
				if issue != nil {
					return
				}

				control.Lock()
				connections[connection] = true
				control.Unlock()

				go handleConnection(connection, connections, control, setting)
			}
		}()

		return ipcContract.Push(&Ipc{
			Listener:    listener,
			Connections: connections,
			Setting:     setting,
			Control:     control,
		})
	}); issue != nil {
		panic("setup ipc contract done with issue " + issue.Error())
	}

	trigger(func(ipc *Ipc, breach func(issue string)) <-chan any {
		done := make(chan any, 1)

		go func() {
			<-ipcDone

			StopIpc(ipc)
			breach("ipc done")

			done <- nil
		}()

		return done
	})

	trigger(func(ipc *Ipc, breach func(issue string)) <-chan any {
		done := make(chan any, 1)

		go func() {
			sigterm := make(chan os.Signal, 1)
			signal.Notify(sigterm, syscall.SIGTERM)
			defer signal.Stop(sigterm)

			select {
			case <-sigterm:
				StopIpc(ipc)
				breach("ipc received sigterm")
			case <-ipcDone:
			}

			done <- nil
		}()

		return done
	})

	trigger(func(ipc *Ipc, breach func(issue string)) <-chan any {
		done := make(chan any, 1)

		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, syscall.SIGINT)
			defer signal.Stop(sigint)

			select {
			case <-sigint:
				StopIpc(ipc)
				breach("ipc received sigint")
			case <-ipcDone:
			}

			done <- nil
		}()

		return done
	})

	return ipcContract
}

func handleConnection(
	connection net.Conn,
	connections IpcConnections,
	control *sync.Mutex,
	setting IpcSetting,
) {
	defer func() {
		if issue := connection.Close(); issue != nil && !errors.Is(issue, net.ErrClosed) {
			panic("handle connection on close with issue, " + issue.Error())
		}

		control.Lock()

		delete(connections, connection)

		control.Unlock()
	}()

	if setting.Timeout > 0 {
		if issue := connection.SetDeadline(time.Now().Add(setting.Timeout)); issue != nil {
			if errors.Is(issue, net.ErrClosed) {
				return
			}

			panic("handle connection on set timeout with issue, " + issue.Error())
		}
	}

	var writeBuffer strings.Builder
	readBuffer := make([]byte, 4096)
	var data string

	for {
		bytesRead, issue := connection.Read(readBuffer)

		if issue != nil {
			if errors.Is(issue, io.EOF) {
				return
			}

			if errors.Is(issue, net.ErrClosed) {
				return
			}

			if networkIssue, status := issue.(net.Error); status && networkIssue.Timeout() {
				return
			}

			panic("handle connection on read with issue, " + issue.Error())
		}

		writeBuffer.Write(readBuffer[:bytesRead])

		if setting.Timeout > 0 {
			if issue := connection.SetDeadline(time.Now().Add(setting.Timeout)); issue != nil {
				if errors.Is(issue, net.ErrClosed) {
					return
				}

				panic("handle connection on set timeout with issue, " + issue.Error())
			}
		}

		data = writeBuffer.String()

		if !strings.Contains(data, "\n") {
			continue
		}

		segments := strings.Split(data, "\n")

		remainder := segments[len(segments)-1]
		segments = segments[:len(segments)-1]

		writeBuffer.Reset()
		writeBuffer.WriteString(remainder)

		for _, message := range segments {
			if message == "" {
				continue
			}

			if issue := setting.Invoke(message); issue != nil {
				panic("handle connection on invoke with issue, " + issue.Error())
			}
		}
	}
}

func StopIpc(ipc *Ipc) {
	ipc.Control.Lock()

	if len(ipc.Connections) > 0 {
		for connection := range ipc.Connections {
			if issue := connection.Close(); issue != nil && !errors.Is(issue, net.ErrClosed) {
				panic("stop ipc on close connection with issue, " + issue.Error())
			}

			delete(ipc.Connections, connection)
		}
	}

	if ipc.Listener == nil {
		ipc.Control.Unlock()
		return
	}

	listener := ipc.Listener
	ipc.Listener = nil

	ipc.Control.Unlock()

	if issue := listener.Close(); issue != nil {
		panic("stop ipc on close listener with issue, " + issue.Error())
	}
}

func ResetIpc(ipc *Ipc) (*Ipc, error) {
	StopIpc(ipc)

	ipcContract := SetupIpc(ipc.Setting)

	ipcPipe, issue := ipcContract.Data()
	if issue != nil {
		return nil, issue
	}

	if issue := ipcContract.Initiate(); issue != nil {
		return nil, issue
	}

	return <-ipcPipe, nil
}

func ConfigureIpc(ipc *Ipc, setting IpcSetting) (*Ipc, error) {
	ipc.Setting = setting

	return ResetIpc(ipc)
}

func SendIpcData(address string, data string, setting IpcSetting) *Signal[string] {
	sendSignal := ConstructSignal[string]()

	go func() {
		connection, issue := dialNative(address, setting.Timeout)

		if issue != nil {
			sendSignal.Push(issue.Error())

			return
		}

		defer func() {
			if issue := connection.Close(); issue != nil && !errors.Is(issue, net.ErrClosed) {
				panic("send data on connection close with issue, " + issue.Error())
			}
		}()

		sendSignal.Listen(func(issue string) any {
			if issue := connection.Close(); issue != nil && !errors.Is(issue, net.ErrClosed) {
				panic("send data on connection close with issue, " + issue.Error())
			}

			return nil
		})

		if setting.Timeout > 0 {
			if issue := connection.SetDeadline(time.Now().Add(setting.Timeout)); issue != nil {
				sendSignal.Push(issue.Error())

				return
			}
		}

		_, issue = connection.Write([]byte(data + "\n"))

		if issue != nil {
			sendSignal.Push(issue.Error())

			return
		}

		sendSignal.Push("")
	}()

	return sendSignal
}

func ListenIpc(
	address string,
	handler func(data string) error,
	setting IpcSetting,
) *Signal[string] {
	listenSignal := ConstructSignal[string]()

	go func() {
		connection, issue := dialNative(address, setting.Timeout)

		if issue != nil {
			if issue := listenSignal.Push(issue.Error()); issue != nil {
				panic("listen ipc then signal push on connect with issue, " + issue.Error())
			}

			return
		}

		defer func() {
			if issue := connection.Close(); issue != nil && !errors.Is(issue, net.ErrClosed) {
				panic("listen ipc on connection close with issue, " + issue.Error())
			}
		}()

		listenSignal.Listen(func(issue string) any {
			if issue := connection.Close(); issue != nil && !errors.Is(issue, net.ErrClosed) {
				panic("listen ipc on connection close with issue, " + issue.Error())
			}

			return nil
		})

		if setting.Timeout > 0 {
			if issue := connection.SetDeadline(time.Now().Add(setting.Timeout)); issue != nil {
				listenSignal.Push(issue.Error())

				return
			}
		}

		var writeBuffer strings.Builder
		readBuffer := make([]byte, 4096)
		var data string

		for {
			bytesRead, issue := connection.Read(readBuffer)

			if issue != nil {
				if errors.Is(issue, io.EOF) {
					listenSignal.Push(issue.Error())

					return
				}

				if errors.Is(issue, net.ErrClosed) {
					listenSignal.Push(issue.Error())

					return
				}

				if networkIssue, status := issue.(net.Error); status && networkIssue.Timeout() {
					listenSignal.Push(issue.Error())

					return
				}

				if issue := listenSignal.Push(issue.Error()); issue != nil {
					panic("listen ipc then signal push on read data with issue, " + issue.Error())
				}

				return
			}

			writeBuffer.Write(readBuffer[:bytesRead])

			if setting.Timeout > 0 {
				if issue := connection.SetDeadline(time.Now().Add(setting.Timeout)); issue != nil {
					listenSignal.Push(issue.Error())

					return
				}
			}

			data = writeBuffer.String()

			if !strings.Contains(data, "\n") {
				continue
			}

			segments := strings.Split(data, "\n")

			remainder := segments[len(segments)-1]
			segments = segments[:len(segments)-1]

			writeBuffer.Reset()
			writeBuffer.WriteString(remainder)

			for _, message := range segments {
				if message == "" {
					continue
				}

				if issue := handler(message); issue != nil {
					listenSignal.Push(issue.Error())

					return
				}
			}
		}
	}()

	return listenSignal
}
