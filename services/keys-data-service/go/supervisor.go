package main

import (
	"keys-data-service/libraries"
	"keys-data-service/settings"
	"keys-data-service/utilities"
	"net"
	"os"
	"os/exec"

	"keys-data-service/channels"
)

func SetupTCPListener() (*net.TCPListener, bool) {
	if utilities.CheckWindows() {
		return nil, false
	}

	listener, issue := net.Listen("tcp", ":"+utilities.GetPort())
	if issue != nil {
		libraries.LogIssue("keys/supervisor", "setup tcp listener with issue", map[string]any{"issue": issue})
		utilities.KillService()
	}

	tcpListener, castSucceeded := listener.(*net.TCPListener)
	if !castSucceeded {
		libraries.LogIssue("keys/supervisor", "invalid tcp listener")
		utilities.KillService()
	}

	return tcpListener, true
}

func GetListenerFile(tcpListener *net.TCPListener) (*os.File, bool) {
	if utilities.CheckWindows() {
		return nil, true
	}

	listenerFile, issue := tcpListener.File()
	if issue != nil {
		libraries.LogIssue("keys/supervisor", "get listener file with issue", map[string]any{"issue": issue})

		return nil, false
	}

	return listenerFile, true
}

func ExecuteChildCommand(serviceType string, listenerFile *os.File) *exec.Cmd {
	childCommand := exec.Command(os.Args[0], "child")

	childCommand.Stdout = os.Stdout
	childCommand.Stderr = os.Stderr

	if utilities.CheckWindows() {
		childCommand.Env = append(os.Environ(), settings.ServiceModeVariable+"="+serviceType)
	} else {
		childCommand.Env = append(os.Environ(),
			settings.FileDescriptorUseVariable+"=1",
			"USE_INHERITED_FD=1",
			"SOCKET_FD=3",
			settings.ServiceModeVariable+"="+serviceType)

		childCommand.ExtraFiles = []*os.File{listenerFile}
	}

	return childCommand
}

func RunParent() {
	tcpListener, _ := SetupTCPListener()

	for {
		listenerFile, status := GetListenerFile(tcpListener)
		if !status {
			libraries.Sleep(300)

			continue
		}

		dataCommand := ExecuteChildCommand("data", listenerFile)

		if issue := dataCommand.Start(); issue != nil {
			libraries.LogIssue("keys/supervisor", "child start with issue", map[string]any{"issue": issue})

			if listenerFile != nil {
				_ = listenerFile.Close()
			}
			libraries.Sleep(300)

			continue
		}
		if listenerFile != nil {
			_ = listenerFile.Close()
		}

		issue := dataCommand.Wait()
		exitCode := 0
		if dataCommand.ProcessState != nil {
			exitCode = dataCommand.ProcessState.ExitCode()
		}

		if exitCode == settings.AbortExitCode {
			libraries.LogData("keys/supervisor", "child aborted")
			libraries.LogData("keys/supervisor", "control initiated")

			for {
				controlListenerFile, ok := GetListenerFile(tcpListener)
				if !ok {
					libraries.Sleep(300)

					continue
				}

				controlCommand := ExecuteChildCommand("control", controlListenerFile)
				if issue := controlCommand.Start(); issue != nil {
					libraries.LogIssue("keys/supervisor", "start control with issue", map[string]any{"issue": issue})

					if controlListenerFile != nil {
						_ = controlListenerFile.Close()
					}
					libraries.Sleep(300)

					continue
				}
				if controlListenerFile != nil {
					_ = controlListenerFile.Close()
				}

				_ = controlCommand.Wait()

				controlExitCode := 0
				if controlCommand.ProcessState != nil {
					controlExitCode = controlCommand.ProcessState.ExitCode()
				}
				if controlExitCode == settings.KillExitCode {
					libraries.LogData("keys/supervisor", "control kill initiated")

					return
				}

				libraries.LogData("keys/supervisor", "control done")
				libraries.LogData("keys/supervisor", "child initiated")

				break
			}

			continue
		}

		if issue != nil {
			libraries.LogIssue("keys/supervisor", "child done with issue", map[string]any{"issue": issue})

		} else {
			libraries.LogIssue("keys/supervisor", "child done")
		}

		libraries.Sleep(300)
	}
}

func RunChild() {
	if os.Getenv(settings.ServiceModeVariable) == "control" {
		RunControl()

		if channels.CheckKillStatus() {
			os.Exit(settings.KillExitCode)
		}

		return
	}

	RunService()

	if channels.CheckAbortStatus() {
		os.Exit(settings.AbortExitCode)
	}
}
