package playbook

import (
	"errors"
	"os"
	"reflect"

	"github.com/shhnwangjian/ops-warren/lib"

	"gopkg.in/yaml.v3"
)

// https://docs.ansible.com/ansible/latest/modules/file_module.html#examples

type FileState int

const (
	Src FileState = iota
	Dest
	Owner
	Group
	State
	Mode
	Path
)

func (p FileState) String() string {
	switch p {
	case Src:
		return "Src"
	case Dest:
		return "Dest"
	case Owner:
		return "Owner"
	case Group:
		return "Group"
	case State:
		return "State"
	case Mode:
		return "Mode"
	case Path:
		return "Path"
	default:
		return "Unknown"
	}
}

type FileInfo struct {
}

type FileBook struct {
	Name string `yaml:"name"`
	Op   FileOp `yaml:"file"`
}

type FileOp struct {
	Src   string `yaml:"src"`
	Dest  string `yaml:"dest"`
	Owner string `yaml:"owner"`
	Group string `yaml:"group"`
	State string `yaml:"state"`
	Mode  string `yaml:"mode"`
	Path  string `yaml:"path"`
	Op    FileState
}

var t = `
- name: Give insecure permissions to an existing file
  file:
    path: /work
    owner: root
    group: root
    mode: '1777'
`

func readYamlConfig(s string) ([]FileBook, error) {
	conf := make([]FileBook, 0)
	err := yaml.Unmarshal([]byte(s), &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func parse(s string) ([]FileBook, error) {
	return readYamlConfig(s)
}

func Parse(s string) ([]FileBook, error) {
	return parse(s)
}

func (f *FileInfo) Do(s string) error {
	l, err := Parse(s)
	if err != err {
		return err
	}

	if len(l) == 0 {
		return errors.New("file book no data")
	}

	return nil
}

func (f *FileOp) StateExec() error {
	if !f.isOp(State) {
		return errors.New(" no state")
	}
	switch f.State {
	case "absent":

	case "touch":

	case "link":
		return f.link()
	case "directory":

	default:
		return errors.New(" no state")
	}
	return nil
}

// isOp 检测操作key的值是否存在
func (f *FileOp) isOp(s FileState) bool {
	ref := reflect.ValueOf(f)
	k, b := ref.Type().FieldByName(s.String())
	if !b {
		return false
	}
	if ref.FieldByName(k.Name).String() == "" {
		return false
	}
	return true
}

func (f *FileOp) link() error {
	if !f.isOp(Src) {
		return errors.New(" no src path")
	}
	if !f.isOp(Dest) {
		return errors.New(" no dest path")
	}
	return os.Symlink(f.Src, f.Dest)
}

func (f *FileOp) chown() error {
	if !f.isOp(Owner) {
		return errors.New(" no owner path")
	}
	if !f.isOp(Group) {
		return errors.New(" no group path")
	}
	uid, gid, err := lib.GetUidGid(f.Owner)
	if err != nil {
		return err
	}
	return os.Chown(f.Src, uid, gid)
}

func init() {
	Register("file", &FileInfo{})
}
