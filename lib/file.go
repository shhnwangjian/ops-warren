package lib

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// CopyFile 覆盖文件
func CopyFile(src, des string, perm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("os.Open Error, %s", err.Error())
	}
	defer srcFile.Close()

	desFile, err := os.OpenFile(des, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("os.OpenFile Error, %s", err.Error())
	}
	defer desFile.Close()

	_, err = io.Copy(desFile, srcFile)
	if err != nil {
		return fmt.Errorf("io.Copy Error, %s", err.Error())
	}
	return nil
}

// PathExists 判断文件或文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

// GetOneWalk 获取指定目录下一层的所有文件和目录，不包含子目录
func GetOneWalk(path string) (files []string, dirs []string, err error) {
	dirOrFile, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, line := range dirOrFile {
		if line.IsDir() {
			dirs = append(dirs, line.Name())
		} else {
			files = append(files, line.Name())
		}
	}
	return
}

func Read(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer fi.Close()

	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return ""
	}
	return string(fd)
}

func GetUidGid(user string) (uid, gid int, err error) {
	con := strings.TrimSpace(Read(HostEtc("passwd")))
	conList := strings.Split(con, "\n")
	if len(conList) == 0 {
		return -1, -1, errors.New("user no exist")
	}
	for _, v := range conList {
		if !strings.EqualFold(user, strings.Split(v, ":")[0]) {
			continue
		}
		uid, err := strconv.Atoi(strings.Split(v, ":")[2])
		if err != nil {
			return -1, -1, err
		}
		gid, err := strconv.Atoi(strings.Split(v, ":")[3])
		if err != nil {
			return -1, -1, err
		}
		return uid, gid, nil
	}
	return -1, -1, errors.New("user no exist")
}

// http://www.filepermissions.com/file-permissions-index
func GetChmodPermissions(s string) int {
	switch s {
	case "---": // No permission
		return 0
	case "--x": // Execute permission
		return 1
	case "-w-": // Write permission
		return 2
	case "-wx": // Execute and write permission: 1 (execute) + 2 (write) = 3
		return 3
	case "r--": // Read permission
		return 4
	case "r-x": // Read and execute permission: 4 (read) + 1 (execute) = 5
		return 5
	case "rw-": // Read and write permission: 4 (read) + 2 (write) = 6
		return 6
	case "rwx": // All permissions: 4 (read) + 2 (write) + 1 (execute) = 7
		return 7
	default:
		return 4
	}
}
