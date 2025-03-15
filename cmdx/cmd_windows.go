package cmdx

import (
	"os"
	"os/exec"
	"syscall"
)

func setProcessGroup(c *exec.Cmd) *exec.Cmd {
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.SysProcAttr.HideWindow = true
	return c
}

// Stop stops the command by sending its process group a SIGTERM signal.
// Stop is idempotent.
// An error should only be returned in the rare case that Stop is called immediately after the command ends but before Start can update its internal state.
func terminateProcess(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

// // terminate terminate the process and all its children in Windows
// func terminate(pid int) (err error) {
// 	// Open a handle to the process with PROCESS_TERMINATE access
// 	handle, err := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, uint32(pid))
// 	if err != nil {
// 		return
// 	}
// 	defer func() { _ = syscall.CloseHandle(handle) }()
//
// 	// Get the list of children process IDs
// 	children, err := getProcessChildren(pid)
// 	if err != nil {
// 		return
// 	}
//
// 	// Kill the child processes first
// 	for _, childPid := range children {
// 		if err = terminate(childPid); err != nil {
// 			return
// 		}
// 	}
//
// 	// Kill the process
// 	if err = syscall.TerminateProcess(handle, 0); err != nil {
// 		return
// 	}
//
// 	return
// }
//
// // getProcessChildren gets the list of child process IDs
// func getProcessChildren(pid int) (children []int, err error) {
// 	// Create a snapshot of the process list
// 	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
// 	if err != nil {
// 		return
// 	}
// 	defer func() { _ = syscall.CloseHandle(snapshot) }()
//
// 	// Get the first process in the list
// 	var procEntry syscall.ProcessEntry32
// 	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
// 	err = syscall.Process32First(snapshot, &procEntry)
// 	if err != nil {
// 		return
// 	}
//
// 	// Find the parent process and its children
// 	for {
// 		if procEntry.ProcessID == uint32(pid) {
// 			// Found the parent process, add its children to the list
// 			for {
// 				err := syscall.Process32Next(snapshot, &procEntry)
// 				if err != nil {
// 					break
// 				}
// 				if procEntry.ParentProcessID == uint32(pid) {
// 					children = append(children, int(procEntry.ProcessID))
// 				}
// 			}
// 			break
// 		}
// 		err = syscall.Process32Next(snapshot, &procEntry)
// 		if err != nil {
// 			return
// 		}
// 	}
//
// 	return
// }
