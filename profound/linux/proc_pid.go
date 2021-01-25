package linux

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

type FdFilesStat struct {
	Path string `json:"path"`
	Fd   uint64 `json:"fd"`
}

// getFdList /proc/(pid)/fd
func getFdList(pid string) ([]*FdFilesStat, error) {
	var fdFiles []*FdFilesStat
	statPath := filepath.Join(DefaultProcMountPoint, pid, "fd")
	d, err := os.Open(statPath)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	fNames, err := d.Readdirnames(-1)
	if len(fNames) == 0 {
		return nil, errors.New("no data")
	}
	for _, fd := range fNames {
		fPath := filepath.Join(statPath, fd)
		filePath, err := os.Readlink(fPath)
		if err != nil {
			continue
		}
		t, err := strconv.ParseUint(fd, 10, 64)
		if err != nil {
			continue
		}
		fdFiles = append(fdFiles, &FdFilesStat{
			Path: filePath,
			Fd:   t,
		})
	}
	return fdFiles, err
}
