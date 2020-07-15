package lib

import (
	"os"
	"path/filepath"
)

// GetEnvPath 根据环境变量获取内容
func GetEnvPath(key string, silent string, path ...string) string {
	value := os.Getenv(key)
	if value == "" {
		value = silent
	}

	switch len(path) {
	case 0:
		return value
	case 1:
		return filepath.Join(value, path[0])
	default:
		all := make([]string, len(path)+1)
		all[0] = value
		copy(all[1:], path)
		return filepath.Join(all...)
	}
}

// HostProc proc路径
func HostProc(path ...string) string {
	return GetEnvPath("", "/proc", path...)
}

// HostSys sys路径
func HostSys(path ...string) string {
	return GetEnvPath("", "/sys", path...)
}

// HostEtc etc路径
func HostEtc(path ...string) string {
	return GetEnvPath("", "/etc", path...)
}

// HostVar var路径
func HostVar(path ...string) string {
	return GetEnvPath("", "/var", path...)
}

// HostHome home路径
func HostHome(path ...string) string {
	return GetEnvPath("", "/home", path...)
}

// HostDev dev路径
func HostDev(path ...string) string {
	return GetEnvPath("", "/dev", path...)
}

// HostLib lib路径
func HostLib(path ...string) string {
	return GetEnvPath("", "/lib", path...)
}

// GetPath 生成路径
func GetPath(silent string, path ...string) string {
	return GetEnvPath("", silent, path...)
}
