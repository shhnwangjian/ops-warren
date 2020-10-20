package playbook

import (
	"errors"
	"fmt"
	"path"

	"github.com/shhnwangjian/ops-warren/lib"
	"gopkg.in/yaml.v3"
)

type ShellInfo struct {
	Msg []string
}

type ShellBook struct {
	Name         string    `yaml:"name"`
	Shell        string    `yaml:"shell"`
	Command      string    `yaml:"command"`
	Args         ShellArgs `yaml:"args"`
	IgnoreErrors bool      `yaml:"ignore_errors"`
}

type ShellArgs struct {
	Chdir      string `yaml:"chdir"`
	Creates    string `yaml:"creates"` // A filename, when it already exists, this step will not be run.
	Removes    string `yaml:"removes"` // A filename, when it does not exist, this step will not be run.
	Executable string `yaml:"executable"`
	User       string `yaml:"user"`
	Group      string `yaml:"group"`
	TimeOut    int    `yaml:"timeout"`
}

func (f *ShellInfo) readYamlConfigList(s string) ([]*ShellBook, error) {
	conf := make([]*ShellBook, 0)
	err := yaml.Unmarshal([]byte(s), &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (f *ShellInfo) parseList(s string) ([]*ShellBook, error) {
	return f.readYamlConfigList(s)
}

func (f *ShellInfo) readYamlConfig(s string) (*ShellBook, error) {
	conf := &ShellBook{}
	err := yaml.Unmarshal([]byte(s), conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (f *ShellInfo) parse(s string) (*ShellBook, error) {
	return f.readYamlConfig(s)
}

func (f *ShellInfo) Res(r ResPlayBook) {
	f.Msg = append(f.Msg, fmt.Sprintf("执行步骤名称:%s,执行模块:%s,执行结果:%s,%s\n", r.Name, r.Model, getStatus(r.Status), r.Msg))
}

func (f *ShellInfo) DoAll(s string) {
	r := ResPlayBook{
		Model: "shell",
	}
	l, err := f.parseList(s)
	if err != nil {
		r.Status = -1
		r.Msg = err.Error()
		f.Res(r)
		return
	}
	if len(l) == 0 {
		r.Status = -1
		r.Msg = "shell book no data"
		f.Res(r)
		return
	}
	for _, line := range l {
		r.Status, r.Msg = line.run()
		f.Res(r)
	}
}

func (f *ShellInfo) DoOne(s string) {
	r := ResPlayBook{
		Model: "shell",
	}
	defer f.Res(r)
	l, err := f.parse(s)
	if err != nil {
		r.Status = -1
		r.Msg = err.Error()
		return
	}
	r.Status, r.Msg = l.run()
}

func (f *ShellBook) getShell() (string, error) {
	if f.Shell != "" {
		return f.Shell, nil
	}
	if f.Command != "" {
		return f.Command, nil
	}
	return "", errors.New(" no shell content")
}

func (f *ShellArgs) isExecShell() bool {
	if f.Creates != "" {
		st, _ := lib.PathExists(path.Join(f.Chdir, f.Creates))
		if st {
			return true
		}
		return false
	}
	if f.Removes != "" {
		st, _ := lib.PathExists(path.Join(f.Chdir, f.Removes))
		if st {
			return true
		}
		return false
	}
	return false
}

func (f *ShellBook) run() (int8, string) {
	if f.Args.isExecShell() {
		return 0, "this step will not be run"
	}
	n := NewCommand()
	cmd, err := f.getShell()
	if err != nil {
		return -1, err.Error()
	}
	n.content = cmd
	if f.Args.Chdir != "" {
		n.dir = f.Args.Chdir
	}
	if f.Args.User != "" {
		n.user = f.Args.User
	}
	if f.Args.Group != "" {
		n.group = f.Args.Group
	}
	if f.Args.TimeOut != 0 {
		n.timeout = f.Args.TimeOut
	}
	err = n.Run()
	if err != nil {
		return -1, err.Error()
	}
	return 0, fmt.Sprintf("STDOUT:%s,STDERR:%s", n.GetStdout(), n.GetStderr())
}

func init() {
	Register("shell", &ShellInfo{})
}
