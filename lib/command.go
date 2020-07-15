package lib

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// ExecShell 执行命令返回结果
func ExecShell(s string) (string, error) {
	cmd := &exec.Cmd{}
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", s)
	} else {
		cmd = exec.Command("/bin/sh", "-c", s)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		str := fmt.Sprint(err)
		if strings.Contains(str, "exit") {
			err = nil
		}
	}
	return out.String(), err
}

func ExecShellByTimeout(timeout int, command string, uid, gid uint32) (string, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctxt, "/bin/sh", "-c", command)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	//cmd.SysProcAttr.Cloneflags = syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uid,
		Gid: gid,
	}

	cmd.Env = append(os.Environ(),
		"LANG=en_US.UTF-8",
	)

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("cmd.Start: %s", err.Error())
	}

	if err := cmd.Wait(); err != nil {
		return buf.String(), fmt.Errorf("cmd.Wait: %s", err)
	}

	return buf.String(), nil
}
