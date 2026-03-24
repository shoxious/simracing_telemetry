//go:build windows

package main

import (
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32main              = windows.NewLazyDLL("kernel32.dll")
	procSetConsoleCtrlHandler = kernel32main.NewProc("SetConsoleCtrlHandler")
	procSetConsoleTitle       = kernel32main.NewProc("SetConsoleTitleW")
)

func initConsole() {
	// Set the terminal window title
	titlePtr, _ := windows.UTF16PtrFromString("SimRacing Dashboard")
	procSetConsoleTitle.Call(uintptr(unsafe.Pointer(titlePtr)))

	// Register a handler so closing the terminal window triggers graceful shutdown.
	// CTRL_CLOSE_EVENT = 2  (user clicked X on the console window)
	cb := syscall.NewCallback(func(ctrlType uint32) uintptr {
		if ctrlType == 2 {
			select {
			case shutdownCh <- struct{}{}:
			default:
			}
			// Windows gives us ~5 s before force-killing the process.
			// Sleep so main() has time to clean up before the OS kills us.
			time.Sleep(4 * time.Second)
			return 1 // tell Windows we handled it
		}
		return 0 // let default handler run for Ctrl+C / Ctrl+Break
	})
	procSetConsoleCtrlHandler.Call(cb, 1)
}
