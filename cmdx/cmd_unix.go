//go:build !windows

package cmdx

import (
	"os/exec"
	"syscall"
)

func terminateProcess(pid int) error {
	// Signal the process group (-pid), not just the process, so that the process
	// and all its children are signaled. Else, child procs can keep running and
	// keep the stdout/stderr fd open and cause cmd.Wait to hang.
	return syscall.Kill(-pid, syscall.SIGTERM)
}

func setProcessGroup(c *exec.Cmd) *exec.Cmd {
	// Set process group ID so the cmd and all its children become a new
	// process group. This allows Stop to SIGTERM the cmd's process group
	// without killing this process (i.e. this code here).
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return c
}
