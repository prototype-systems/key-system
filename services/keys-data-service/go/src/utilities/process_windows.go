//go:build windows

package utilities

// # Windows Firewall Auto-Allow — How It Works
//
// ## The Problem
//
// When a Go service opens a network port (e.g. starts a TCP server), Windows
// immediately checks whether an inbound firewall rule exists for that program.
// If no rule is found, Windows pops up a security alert dialog asking the user
// to allow or block the connection — every single time the service starts.
//
// ## The Goal
//
// Register a permanent "allow inbound" firewall rule for this service's
// executable the first time it runs, so Windows never shows the popup again.
//
// ## Why It's Complicated
//
// Writing a firewall rule requires Administrator privileges. The service itself
// does not run as Administrator. We need to ask Windows to temporarily elevate
// a separate helper program (netsh) just to write that one rule.
//
// Additionally, writing the rule must be fully finished BEFORE the service
// opens its port — otherwise the race condition brings the popup back.
//
// ## Step-by-Step Workflow
//
// 1. CHECK — silently ask Windows "does a firewall rule with our service name
//    already exist?" using netsh. If yes, skip everything and return immediately.
//    This is the fast path on every run after the first.
//
// 2. BUILD THE COMMAND — assemble the netsh arguments that will add an inbound
//    allow rule for this specific executable file.
//
// 3. LOAD THE WINDOWS API — load ShellExecuteExW from shell32.dll. This is the
//    Windows function that can launch a program with elevated (Admin) privileges
//    by triggering a UAC prompt.
//
// 4. CONFIGURE THE LAUNCH — fill in a SHELLEXECUTEINFO structure that tells
//    Windows: run "netsh" with our arguments, use the "runas" verb (which means
//    "run as Administrator"), hide the console window, and crucially — set the
//    SEE_MASK_NOCLOSEPROCESS flag so Windows hands us back a process handle
//    instead of discarding it.
//
// 5. LAUNCH AND WAIT — call ShellExecuteExW, which shows the UAC prompt to the
//    user (one time only). Once approved, netsh runs elevated and writes the
//    firewall rule. We call WaitForSingleObject on the process handle, which
//    BLOCKS this function until netsh finishes. Only then does execution continue
//    and the service proceeds to open its port — no more race condition.
//
// 6. CLEANUP — close the process handle to release the Windows resource.
//
// ## After the First Run
//
// The firewall rule is now permanent in Windows Firewall. Every subsequent
// start goes: check rule → found → return immediately. No UAC, no popup, no delay.

import (
	"keys-data-service/settings"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	shell32Library      = windows.NewLazySystemDLL("shell32.dll")
	shellExecuteExWProc = shell32Library.NewProc("ShellExecuteExW")
)

// shellExecuteInfo mirrors the Win32 SHELLEXECUTEINFOW struct.
// Windows requires this exact memory layout to be passed to ShellExecuteExW.
// The field names use Win32 Hungarian notation and must not be renamed.
type shellExecuteInfo struct {
	cbSize         uint32
	fMask          uint32
	hwnd           uintptr
	lpVerb         *uint16
	lpFile         *uint16
	lpParameters   *uint16
	lpDirectory    *uint16
	nShow          int32
	hInstApp       uintptr
	lpIDList       uintptr
	lpClass        *uint16
	hkeyClass      uintptr
	dwHotKey       uint32
	hIconOrMonitor uintptr

	// populated by Windows when SEE_MASK_NOCLOSEPROCESS is set
	hProcess windows.Handle
}

// keepProcessHandleOpenFlag tells ShellExecuteExW to keep the launched process
// handle open and store it in shellExecuteInfo.hProcess so we can wait on it.
const keepProcessHandleOpenFlag = 0x00000040

func AllowProcess() {
	executablePath := GetExecutablePath()
	firewallRuleName := settings.ServiceName

	// Step 1: Check — silently run "netsh show rule". Exit code 0 means the rule
	// already exists, so we skip everything and return. HideWindow suppresses the
	// console flash.
	ruleCheckCommand := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", "name="+firewallRuleName)
	ruleCheckCommand.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if ruleCheckCommand.Run() == nil {
		return
	}

	// Step 2: Build the netsh command arguments that will add the inbound allow
	// rule for this executable. ShellExecuteExW takes a single parameter string,
	// not a slice, so we join the arguments with spaces.
	netshArguments := strings.Join([]string{
		"advfirewall", "firewall", "add", "rule",
		"name=" + firewallRuleName,
		"dir=in",
		"action=allow",
		"program=" + executablePath,
		"enable=yes",
		"profile=any",
	}, " ")

	// Steps 3 & 4: Configure the elevated launch. Convert strings to UTF-16
	// pointers because that is what the Windows API expects.
	elevationVerb, _ := syscall.UTF16PtrFromString("runas") // "runas" = run as Administrator
	executablePointer, _ := syscall.UTF16PtrFromString("netsh")
	argumentsPointer, _ := syscall.UTF16PtrFromString(netshArguments)

	shellExecuteParameters := &shellExecuteInfo{
		fMask:        keepProcessHandleOpenFlag, // keep hProcess open so we can wait
		lpVerb:       elevationVerb,
		lpFile:       executablePointer,
		lpParameters: argumentsPointer,
		nShow:        windows.SW_HIDE, // hide the netsh console window
	}
	shellExecuteParameters.cbSize = uint32(unsafe.Sizeof(*shellExecuteParameters)) // Windows requires the struct size

	// Step 5: Launch netsh elevated. Windows shows the UAC prompt here (once).
	// After approval, netsh writes the firewall rule with Admin rights.
	shellExecuteExWProc.Call(uintptr(unsafe.Pointer(shellExecuteParameters)))

	if shellExecuteParameters.hProcess != 0 {
		// Block until netsh finishes — the rule must be written before the
		// service opens its port, otherwise the security dialog races back.
		_, _ = windows.WaitForSingleObject(shellExecuteParameters.hProcess, windows.INFINITE)

		// Step 6: Release the process handle.
		_ = windows.CloseHandle(shellExecuteParameters.hProcess)
	}
}
