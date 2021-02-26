package playbook

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
	"unicode"

	"github.com/shhnwangjian/ops-warren/lib"
)

type Command struct {
	cmd       *exec.Cmd
	bufOut    bytes.Buffer
	bufErr    bytes.Buffer
	startTime time.Time
	stopTime  time.Time
	lock      sync.RWMutex
	env       []string
	dir       string
	user      string
	group     string
	content   string
	timeout   int
}

func NewCommand() *Command {
	return &Command{
		cmd:       nil,
		startTime: time.Unix(0, 0),
		stopTime:  time.Unix(0, 0),
		env:       make([]string, 0),
		timeout:   60,
		user:      "root",
		group:     "root",
	}
}

func (p *Command) setEnv() {
	if len(p.env) != 0 {
		p.cmd.Env = append(os.Environ(), p.env...)
	} else {
		p.cmd.Env = os.Environ()
	}
}

func (p *Command) setDir() {
	if p.dir != "" {
		p.cmd.Dir = p.dir
	}
}

func (p *Command) waitForExit() error {
	err := p.cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (p *Command) createCommand(ctxt context.Context) error {
	args, err := ParseCommand(p.content)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return fmt.Errorf("not command")
	}
	p.cmd = exec.CommandContext(ctxt, args[0])
	if len(args) > 1 {
		p.cmd.Args = args
	}
	p.cmd.SysProcAttr = &syscall.SysProcAttr{}
	if p.setUser() != nil {
		return fmt.Errorf("fail to set user")
	}
	setDeathsig(p.cmd.SysProcAttr)
	p.setEnv()
	p.setDir()
	p.setStdout()
	p.setStderr()
	return nil
}

func (p *Command) Run() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	ctxt, cancel := context.WithTimeout(context.Background(), time.Duration(p.timeout)*time.Second)
	defer cancel()
	p.startTime = time.Now()
	err := p.createCommand(ctxt)
	if err != nil {
		return err
	}
	err = p.cmd.Start()
	if err != nil {
		return err
	}
	err = p.waitForExit()
	if err != nil {
		return err
	}
	p.stopTime = time.Now()
	return nil
}

func (p *Command) setUser() error {
	uid, gid, err := lib.GetUidGid(p.user, p.group)
	if err != nil {
		return err
	}
	setUserID(p.cmd.SysProcAttr, uint32(uid), uint32(gid))
	return nil
}

func (p *Command) setStdout() {
	p.cmd.Stdout = &p.bufOut
}

func (p *Command) setStderr() {
	p.cmd.Stderr = &p.bufErr
}

func (p *Command) GetStdout() string {
	return p.bufOut.String()
}

func (p *Command) GetStderr() string {
	return p.bufErr.String()
}

func setDeathsig(sysProcAttr *syscall.SysProcAttr) {
	sysProcAttr.Setpgid = true
	//sysProcAttr.Pdeathsig = syscall.SIGKILL
}

func setUserID(procAttr *syscall.SysProcAttr, uid uint32, gid uint32) {
	procAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid, NoSetGroups: true}
}

func skipSpace(s string, offset int) int {
	for i := offset; i < len(s); i++ {
		if !unicode.IsSpace(rune(s[i])) {
			return i
		}
	}
	return -1
}

func findChar(s string, offset int, ch byte) int {
	for i := offset; i < len(s); i++ {
		if s[i] == '\\' {
			i++
		} else if s[i] == ch {
			return i
		}
	}
	return -1
}

func appendArgument(arg string, args []string) []string {
	if arg[0] == '"' || arg[0] == '\'' {
		return append(args, arg[1:len(arg)-1])
	}
	return append(args, arg)
}

func ParseCommand(command string) ([]string, error) {
	args := make([]string, 0)
	cmdLen := len(command)
	for i := 0; i < cmdLen; {
		//find the first non-space char
		j := skipSpace(command, i)
		if j == -1 {
			break
		}
		i = j
		for ; j < cmdLen; j++ {
			if unicode.IsSpace(rune(command[j])) {
				args = appendArgument(command[i:j], args)
				i = j + 1
				break
			} else if command[j] == '\\' {
				j++
			} else if command[j] == '"' || command[j] == '\'' {
				k := findChar(command, j+1, command[j])
				if k == -1 {
					args = appendArgument(command[i:], args)
					i = cmdLen
				} else {
					args = appendArgument(command[i:k+1], args)
					i = k + 1
				}
				break
			}
		}
		if j >= cmdLen {
			args = appendArgument(command[i:], args)
			i = cmdLen
		}
	}
	if len(args) <= 0 {
		return nil, fmt.Errorf("no command from empty string")
	}
	return args, nil
}
