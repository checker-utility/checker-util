package checkerutil

import (
	"fmt"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// ConsoleTitle referrs to configurations pertaining to the console title bar
type ConsoleTitle struct {
	Activated bool
	// Delay referrs to how often the title bar will refresh
	Delay time.Duration
	/*
		Format can contain a few things. Here's an example for Format:
			Made by postuwu#0123 | CPM {CPM} | {ID:good}/{ID:bad} | Checked {ID:}

		ID referrs to *Output.ID, if ID is empty then it will add up all the IDs together.
	*/
	Format string
}

// setConsoleTitle sets to console title based off the c.ConsoleTitle.Format
// Only called by the initiator
// TODO: support windows, linux, and MACOS because atm it only supports windows
func (c *Checker) setConsoleTitle() {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return
	}
	defer syscall.FreeLibrary(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return
	}
	for {
		if c.Console.Activated {
			time.Sleep(c.Console.Delay)
			f := c.formatConsole()
			syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(f))), 0, 0)
		}
	}
}

var (
	titleIDRegex = regexp.MustCompile(`\{ID:(|.+?)\}`)
)

// formatConsole converts *ConsleTitle.Format to a string with {ID:} replaced and {CPM}
func (c *Checker) formatConsole() string {
	end := c.Console.Format
	IDs := titleIDRegex.FindAllString(c.Console.Format, -1)
	for i, ID := range IDs {
		ID = strings.Split(ID, ":")[1]
		ID = ID[:len(ID)-1]
		total := 0
		if a, ok := c.Outputs[ID]; ok {
			end = strings.ReplaceAll(end, IDs[i], fmt.Sprint(a.InputNum.Num))
			continue
		}
		if ID == "" {
			for _, a := range c.Outputs {
				a.InputNum.Mutex.Lock()
				total += a.InputNum.Num
				a.InputNum.Mutex.Unlock()
			}
			end = strings.ReplaceAll(end, IDs[i], fmt.Sprint(total))
			continue
		}
		end = strings.ReplaceAll(end, IDs[i], "N/A")
	}
	end = strings.ReplaceAll(end, "{CPM}", fmt.Sprint(c.CalculateCPM()))
	return end
}
