package main

import (
	"flag"
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                       = windows.NewLazyDLL("user32.dll")
	kernel32                     = windows.NewLazyDLL("kernel32.dll")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procFindWindowW              = user32.NewProc("FindWindowW")
	procSetForegroundWindow      = user32.NewProc("SetForegroundWindow")
	procSendInput                = user32.NewProc("SendInput")
	procGetCurrentThreadId       = kernel32.NewProc("GetCurrentThreadId")
)

type KEYBDINPUT struct {
	Type  uint32
	Wvk   uint16
	WScan uint16
	Flags uint32
	Time  uint32
	Extra uintptr
}

const (
	VK_RETURN  = 0x0D
	VK_CONTROL = 0x11
	VK_ALT     = 0x12
	VK_SHIFT   = 0x10
)

// go run macro.go -pid YOUR_PID -short-delay 50 -long-delay 3000
func main() {
	pidFlag := flag.Int("pid", 0, "Target process ID")
	shortDelay := flag.Int("short-delay", 50, "Short delay in milliseconds")
	longDelay := flag.Int("long-delay", 3000, "Long delay in milliseconds")
	flag.Parse()

	if *pidFlag == 0 {
		fmt.Println("Please provide a PID using -pid flag")
		return
	}

	// Find window by PID
	hwnd := findWindowByPID(*pidFlag)
	if hwnd == 0 {
		fmt.Printf("No window found for PID %d\n", *pidFlag)
		return
	}

	// Run the macro in an infinite loop
	fmt.Println("Macro started. Press Ctrl+C to stop...")
	for {
		sendSpecialKey(0x09) // VK_TAB
		time.Sleep(time.Duration(*shortDelay) * time.Millisecond)

		sendSpecialKey(0x20) // VK_SPACE
		time.Sleep(time.Duration(*longDelay) * time.Millisecond)
	}
}

func findWindowByPID(pid int) uintptr {
	var hwnd uintptr
	cb := windows.NewCallback(func(h windows.Handle, p uintptr) uintptr {
		var currentPID uint32
		procGetWindowThreadProcessId.Call(uintptr(h), uintptr(unsafe.Pointer(&currentPID)))

		if uint32(pid) == currentPID {
			hwnd = uintptr(h)
			return 0 // Stop enumeration
		}
		return 1 // Continue enumeration
	})

	user32.NewProc("EnumWindows").Call(cb, 0)
	return hwnd
}

func sendSpecialKey(vk uint16) {
	// Get the foreground window's thread
	foregroundWindow, _, _ := user32.NewProc("GetForegroundWindow").Call()
	foregroundThreadId, _, _ := user32.NewProc("GetWindowThreadProcessId").Call(foregroundWindow, 0)

	// Use the correct procedure from kernel32
	currentThreadId, _, _ := procGetCurrentThreadId.Call()

	// Attach the threads
	user32.NewProc("AttachThreadInput").Call(currentThreadId, foregroundThreadId, 1)

	input := KEYBDINPUT{
		Type:  1, // INPUT_KEYBOARD
		Wvk:   vk,
		WScan: 0,
		Flags: 0,
		Time:  0,
		Extra: 0,
	}

	size := unsafe.Sizeof(input)
	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), size)

	// Key up event
	input.Flags = 2 // KEYEVENTF_KEYUP
	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), size)

	// Detach the threads
	user32.NewProc("AttachThreadInput").Call(currentThreadId, foregroundThreadId, 0)
}
