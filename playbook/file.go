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

type FileState int

const (
	Src FileState = iota
	Dest
	Owner
	Group
	State
	Mode
	Path
	Name
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

func readYamlConfig(s string) ([]*FileBook, error) {
	conf := make([]*FileBook, 0)
	err := yaml.Unmarshal([]byte(s), &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func parse(s string) ([]*FileBook, error) {
	return readYamlConfig(s)
}

func Parse(s string) ([]*FileBook, error) {
	return parse(s)
}

func (f *FileInfo) Res(ind int, r ResPlayBook) {
	f.Msg = append(f.Msg, fmt.Sprintf("执行步骤%d(%s),执行模块:%s,执行结果:%s,%s\n", ind+1, r.Name, r.Model, getStatus(r.Status), r.Msg))
}

func (f *FileInfo) Do(s string) error {
	l, err := Parse(s)
	if err != err {
		return err
	}
	if len(l) == 0 {
		return errors.New("file book no data")
	}
	for ind, line := range l {
		r := ResPlayBook{
			Name:  line.Name,
			Model: "file",
		}
		err := line.Op.StateExec()
		if err != nil {
			r.Status = -1
			r.Msg = err.Error()
			f.Res(ind, r)
			continue
		}
		err = line.Op.mode()
		if err != nil {
			r.Status = -1
			r.Msg = err.Error()
			f.Res(ind, r)
			continue
		}
		f.Res(ind, r)
	}
	return nil
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
func (f FileOp) isOp(s FileState) bool {
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
	uid, gid, err := lib.GetUidGid(f.Owner)
	if err != nil {
		return err
	}
	return os.Chown(f.Src, uid, gid)
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

func getNewMode(s string) string {
	var res string
	if strings.Contains(s, "r") {
		res = "r"
	} else {
		res = "-"
	}
	if strings.Contains(s, "w") {
		res += "w"
	} else {
		res += "-"
	}
	if strings.Contains(s, "w") {
		res += "x"
	} else {
		res += "-"
	}
	return res
}

func getModifyMode(s, old string) string {
	var res string
	if strings.Contains(s, "r") {
		res = "-"
	} else {
		res = old[:1]
	}
	if strings.Contains(s, "w") {
		res += "-"
	} else {
		res += old[1:2]
	}
	if strings.Contains(s, "x") {
		res += "-"
	} else {
		res += old[2:3]
	}
	return res
}

func getAddMode(s, old string) string {
	var res string
	if strings.Contains(s, "r") {
		res = "r"
	} else {
		res = old[:1]
	}
	if strings.Contains(s, "w") {
		res += "w"
	} else {
		res += old[1:2]
	}
	if strings.Contains(s, "x") {
		res += "x"
	} else {
		res += old[2:3]
	}
	return res
}

func getFileMode(str, fileModeStr string) (string, string, error) {
	if strings.Contains(str, "=") {
		strList := strings.Split(str, "=")
		if len(strList) != 2 {
			return "", "", errors.New("file mode content error")
		}
		res := getNewMode(strList[1])
		return strList[0], res, nil
	}
	if strings.Contains(str, "-") {
		strList := strings.Split(str, "-")
		if len(strList) != 2 {
			return "", "", errors.New("file mode content error")
		}
		if strList[0] == "u" {
			res := getModifyMode(strList[1], fileModeStr[1:4])
			return strList[0], res, nil
		}
		if strList[0] == "g" {
			res := getModifyMode(strList[1], fileModeStr[4:7])
			return strList[0], res, nil
		}
		if strList[0] == "o" {
			res := getModifyMode(strList[1], fileModeStr[7:10])
			return strList[0], res, nil
		}
	}
	if strings.Contains(str, "+") {
		strList := strings.Split(str, "+")
		if len(strList) != 2 {
			return "", "", errors.New("file mode content error")
		}
		if strList[0] == "u" {
			res := getAddMode(strList[1], fileModeStr[1:4])
			return strList[0], res, nil
		}
		if strList[0] == "g" {
			res := getAddMode(strList[1], fileModeStr[4:7])
			return strList[0], res, nil
		}
		if strList[0] == "x" {
			res := getAddMode(strList[1], fileModeStr[7:10])
			return strList[0], res, nil
		}
	}
	return "", "", errors.New("file mode content error")
}

func checkStr(name, list string) bool {
	nameList := strings.Split(name, "")
	for _, val := range nameList {
		if !strings.Contains(list, val) {
			return false
		}
	}
	return true
}

func init() {
	Register("file", &FileInfo{})
}
