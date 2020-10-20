package playbook

import (
	"errors"
	"strings"
)

var (
	PlayInfo = make(map[string]PlayBook)
)

type PlayBook interface {
	DoAll(string)
	DoOne(string)
	Res(ResPlayBook)
}

type ResPlayBook struct {
	Name   string
	Model  string
	Msg    string
	Status int8
}

func Register(name string, collect PlayBook) {
	if collect == nil {
		panic("config: Register PlayBook is nil")
	}
	if _, ok := PlayInfo[name]; ok {
		panic("config: Register PlayBook twice for adapter " + name)
	}
	PlayInfo[name] = collect
}

func getStatus(i int8) string {
	switch i {
	case 0:
		return "success"
	case -1:
		return "fail"
	default:
		return "Unknown"
	}
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

// getFileMode 获取文件权限
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

// getAddMode 获取添加文件权限
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
