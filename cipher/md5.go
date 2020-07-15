package cipher

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// FileMD5 文件md5计算
func FileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("os.Open Err:%s", err.Error())
	}
	defer f.Close()
	return BytesMD5(f)
}

func ContentMD5(content string) (string, error) {
	md5hash := md5.New()
	_, err := md5hash.Write([]byte(content))
	return fmt.Sprintf("%x", md5hash.Sum(nil)), err
}

// 字节流md5计算
func BytesMD5(br io.Reader) (string, error) {
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, br); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}
