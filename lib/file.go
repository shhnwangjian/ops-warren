package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
