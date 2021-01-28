package playbook

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shhnwangjian/ops-warren/lib"
	"gopkg.in/yaml.v3"
)

// https://docs.ansible.com/ansible/latest/modules/file_module.html#examples

type fileState int

const (
	Src fileState = iota
	Dest
	Owner
	Group
	State
	Mode
	Path
	Name
)

func (p fileState) String() string {
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
	case Name:
		return "Name"
	default:
		return "Unknown"
	}
}

type FileInfo struct {
	Msg []string
}

type FileBook struct {
	Name string  `yaml:"name"`
	Op   *FileOp `yaml:"file"`
}

type FileOp struct {
	Src   string `yaml:"src"`
	Dest  string `yaml:"dest"`
	Name  string `yaml:"name"`
	Owner string `yaml:"owner"`
	Group string `yaml:"group"`
	State string `yaml:"state"`
	Mode  string `yaml:"mode"`
	Path  string `yaml:"path"`
}

func (f *FileInfo) readYamlConfigList(s string) ([]*FileBook, error) {
	conf := make([]*FileBook, 0)
	err := yaml.Unmarshal([]byte(s), &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (f *FileInfo) parseList(s string) ([]*FileBook, error) {
	return f.readYamlConfigList(s)
}

func (f *FileInfo) readYamlConfig(s string) (*FileBook, error) {
	conf := &FileBook{}
	err := yaml.Unmarshal([]byte(s), conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (f *FileInfo) parse(s string) (*FileBook, error) {
	return f.readYamlConfig(s)
}

func (f *FileInfo) Res(r ResPlayBook) {
	f.Msg = append(f.Msg, fmt.Sprintf("执行步骤名称:%s,执行模块:%s,执行结果:%s,%s\n", r.Name, r.Model, getStatus(r.Status), r.Msg))
}

func (f *FileInfo) DoOne(s string) {
	r := ResPlayBook{
		Model: "file",
	}
	l, err := f.parse(s)
	defer f.Res(r)
	if err != nil {
		r.Status = -1
		r.Msg = err.Error()
		return
	}
	err = l.Op.StateExec()
	if err != nil {
		r.Status = -1
		r.Msg = err.Error()
		return
	}
	err = l.Op.mode()
	if err != nil {
		r.Status = -1
		r.Msg = err.Error()
	}
}

func (f *FileInfo) DoAll(s string) {
	r := ResPlayBook{
		Model: "file",
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
		r.Msg = "file book no data"
		f.Res(r)
		return
	}
	for _, line := range l {
		r.Name = line.Name
		err := line.Op.StateExec()
		if err != nil {
			r.Status = -1
			r.Msg = err.Error()
			f.Res(r)
			continue
		}
		err = line.Op.mode()
		if err != nil {
			r.Status = -1
			r.Msg = err.Error()
			f.Res(r)
			continue
		}
		f.Res(r)
	}
}

// Path to the file being managed.  aliases: dest, name
func (f *FileOp) getName() (string, error) {
	if !f.isOp(Name) && !f.isOp(Path) && !f.isOp(Dest) {
		return "", errors.New("no file name")
	}
	if f.isOp(Name) {
		return f.Name, nil
	}
	if f.isOp(Path) {
		return f.Path, nil
	}
	if f.isOp(Dest) {
		return f.Dest, nil
	}
	return "", errors.New(" no file name")
}

func (f *FileOp) StateExec() error {
	if !f.isOp(State) {
		return nil
	}
	switch f.State {
	case "absent":
		return f.absent()
	case "touch":
		return f.touch()
	case "link":
		return f.link()
	case "hard":
		return f.hard()
	case "directory":
		return f.directory()
	default:
		return errors.New(fmt.Sprintf("no state(%s)", f.State))
	}
}

// isOp 检测操作key的值是否存在
func (f *FileOp) isOp(s fileState) bool {
	ref := reflect.ValueOf(*f)
	k, b := ref.Type().FieldByName(s.String())
	if !b {
		return false
	}
	if ref.FieldByName(k.Name).String() == "" {
		return false
	}
	return true
}

func (f *FileOp) hard() error {
	if !f.isOp(Src) {
		return errors.New(" no src path")
	}
	name, err := f.getName()
	if err != nil {
		return err
	}
	if f.isChown() {
		err := f.chown()
		if err != nil {
			return err
		}
	}
	return os.Link(f.Src, name)
}

func (f *FileOp) link() error {
	if !f.isOp(Src) {
		return errors.New(" no file src")
	}
	name, err := f.getName()
	if err != nil {
		return err
	}
	if f.isChown() {
		err := f.chown()
		if err != nil {
			return err
		}
	}
	return os.Symlink(f.Src, name)
}

func (f *FileOp) isChown() bool {
	if !f.isOp(Owner) {
		return false
	}
	return true
}

func (f *FileOp) chown() error {
	uid, gid, err := lib.GetUidGid(f.Owner, f.Group)
	if err != nil {
		return err
	}
	return os.Chown(f.Src, int(uid), int(gid))
}

func (f *FileOp) touch() error {
	name, err := f.getName()
	if err != nil {
		return err
	}
	_, err = os.Stat(name)
	if os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(f.Path, currentTime, currentTime)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileOp) directory() error {
	name, err := f.getName()
	if err != nil {
		return err
	}
	err = os.MkdirAll(name, 0744)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileOp) absent() error {
	name, err := f.getName()
	if err != nil {
		return err
	}
	err = os.RemoveAll(name)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileOp) mode() error {
	if !f.isOp(Mode) {
		return nil
	}
	path, err := f.getName()
	if err != nil {
		return err
	}
	s, err := lib.PathExists(path)
	if !s {
		return err
	}
	numList := "1234567890"
	numStatus := checkStr(f.Mode, numList)
	strList := "+-=ugorwx,"
	strStatus := checkStr(f.Mode, strList)
	if !numStatus && !strStatus {
		return errors.New("file mode content error")
	}
	if numStatus {
		if len(f.Mode) < 4 && len(f.Mode) > 5 {
			return errors.New("file mode content error")
		}
		num, err := strconv.ParseInt(f.Mode, 8, 0)
		if err != nil {
			return err
		}
		err = os.Chmod(f.Path, os.FileMode(num))
		if err != nil {
			return err
		}
		return nil
	}
	if strStatus {
		err := f.getFileMode()
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (f *FileOp) getFileMode() error {
	var u, g, o string
	name, err := f.getName()
	if err != nil {
		return err
	}
	fInfo, err := os.Stat(name)
	if err != nil {
		return err
	}
	fileModeStr := fmt.Sprintf("%v", fInfo.Mode()) // -rw-r--r--
	u = fileModeStr[1:4]
	g = fileModeStr[4:7]
	o = fileModeStr[7:10]
	modeList := strings.Split(f.Mode, ",")
	if len(modeList) == 0 {
		return errors.New("file mode content error")
	}
	for _, line := range modeList {
		role, modeStr, err := getFileMode(line, fileModeStr)
		if err != nil {
			continue
		}
		switch role {
		case "u":
			u = modeStr
		case "g":
			g = modeStr
		case "o":
			o = modeStr
		}
	}
	uNum := lib.GetChmodPermissions(u)
	gNum := lib.GetChmodPermissions(g)
	oNum := lib.GetChmodPermissions(o)
	fMode, err := strconv.ParseInt(fmt.Sprintf("0%d%d%d", uNum, gNum, oNum), 8, 0)
	if err != nil {
		return err
	}
	err = os.Chmod(f.Path, os.FileMode(fMode))
	if err != nil {
		return err
	}
	return nil
}

func init() {
	Register("file", &FileInfo{})
}
