package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
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

func Read(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetUidGid(userName, groupName string) (uid, gid uint64, err error) {
	u, err := user.Lookup(userName)
	if err != nil {
		return 0, 0, err
	}
	uid, err = strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return 0, 0, err
	}
	gid, err = strconv.ParseUint(u.Gid, 10, 32)
	if err != nil && groupName == "" {
		return 0, 0, err
	}
	if groupName != "" {
		g, err := user.LookupGroup(groupName)
		if err != nil {
			return 0, 0, err
		}
		gid, err = strconv.ParseUint(g.Gid, 10, 32)
		if err != nil {
			return 0, 0, err
		}
	}
	return uid, gid, err
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
