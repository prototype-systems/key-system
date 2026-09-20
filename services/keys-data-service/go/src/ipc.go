package main

import (
	"keys-data-service/channels"
	"keys-data-service/libraries"
	"keys-data-service/routes"
	"keys-data-service/utilities"
)

func RunIpc() {
	ipcContract := libraries.SetupIpc(libraries.IpcSetting{
		Address: utilities.GetIpcSocketAddress(),
		Invoke:  routes.Invoke,
	})

	trigger := func(condition func(ipc *libraries.Ipc, breach func(issue string)) <-chan any) {
		if _, issue := ipcContract.Trigger(condition); issue != nil {
			panic("run ipc contract trigger with unrelated issue, " + issue.Error())
		}
	}

	trigger(func(ipc *libraries.Ipc, breach func(issue string)) <-chan any {
		done := make(chan any, 1)

		go func() {
			<-channels.GetKillChannel()

			libraries.StopIpc(ipc)

			breach("ipc done")

			done <- nil
		}()

		return done
	})

	ipcPipe, issue := ipcContract.Data()
	if issue != nil {
		libraries.LogIssue("keys/ipc", "run ipc contract data with issue", map[string]any{"issue": issue})

		utilities.KillService()
	}

	if issue := ipcContract.Initiate(); issue != nil {
		libraries.LogIssue("keys/ipc", "run ipc contract initiate with issue", map[string]any{"issue": issue})

		utilities.KillService()
	}

	libraries.SetIpc(<-ipcPipe)
}
