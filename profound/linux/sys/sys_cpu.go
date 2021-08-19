package sys

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/shhnwangjian/ops-warren/lib"
)

type Cpus struct {
	Path  string
	Cache map[int]Cache
}

type Cache struct {
	Name       string
	CacheIndex map[int]CacheIndex
}

type CacheIndex struct {
	Id                  int
	Level               int
	Type                string
	Size                int // kb
	CoherencyLineSize   int
	SharedCpuList       []string
	WaysOfAssociativity int
	NumberOfSets        int
}

func CacheInstance() *Cpus {
	return &Cpus{Cache: make(map[int]Cache), Path: DefaultSysDeviceCpu}
}

func (c *Cpus) getCpuMap() error {
	f, err := os.Open(c.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	list, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("no cpu dir")
	}

	for i := 0; i < len(list); i++ {
		if !list[i].IsDir() {
			continue
		}
		if !strings.HasPrefix(list[i].Name(), "cpu") {
			continue
		}

		var num int
		numStr := lib.RemovePrefix(list[i].Name(), "cpu")
		if numStr != "" {
			num, _ = strconv.Atoi(numStr)
		}
		indexPath := path.Join(c.Path, list[i].Name(), "cache")
		c.Cache[num] = Cache{list[i].Name(), c.getAllIndexInfo(indexPath)}
	}

	return nil
}

func (c *Cpus) GetIndexInfo() (err error) {
	return c.getCpuMap()
}

func (c *Cpus) getAllIndexInfo(indexPath string) map[int]CacheIndex {
	indexMap := make(map[int]CacheIndex)
	f, err := os.Open(indexPath)
	if err != nil {
		return indexMap
	}
	list, err := f.Readdir(-1)
	if err != nil {
		return indexMap
	}
	if len(list) == 0 {
		return indexMap
	}

	for i, line := range list {
		if !line.IsDir() {
			continue
		}
		indexMap[i], _ = c.getIndexInfo(path.Join(indexPath, line.Name()))
	}
	return indexMap
}

func (c *Cpus) getIndexInfo(dir string) (index CacheIndex, err error) {
	return GetIndexInfo(dir)
}

func GetIndexInfo(dir string) (index CacheIndex, err error) {
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	defer f.Close()

	list, err := f.Readdir(-1)
	if err != nil {
		return
	}
	if len(list) == 0 {
		return CacheIndex{}, fmt.Errorf("%s no file", dir)
	}

	for i := 0; i < len(list); i++ {
		switch list[i].Name() {
		case "id":
			if data, err := lib.Read(path.Join(dir, "id")); err == nil {
				index.Id, _ = strconv.Atoi(strings.TrimSpace(string(data)))
			}
		case "level":
			if data, err := lib.Read(path.Join(dir, "level")); err == nil {
				index.Level, _ = strconv.Atoi(strings.TrimSpace(string(data)))
			}
		case "type":
			if data, err := lib.Read(path.Join(dir, "type")); err == nil {
				index.Type = strings.TrimSpace(string(data))
			}
		case "size":
			if data, err := lib.Read(path.Join(dir, "size")); err == nil {
				index.Size, _ = strconv.Atoi(lib.RemoveSuffix(strings.TrimSpace(string(data)), "K"))
			}
		case "shared_cpu_list":
			if data, err := lib.Read(path.Join(dir, "shared_cpu_list")); err == nil {
				index.SharedCpuList = strings.Split(strings.TrimSpace(string(data)), ",")
			}
		case "coherency_line_size":
			if data, err := lib.Read(path.Join(dir, "coherency_line_size")); err == nil {
				index.CoherencyLineSize, _ = strconv.Atoi(strings.TrimSpace(string(data)))
			}
		case "ways_of_associativity":
			if data, err := lib.Read(path.Join(dir, "ways_of_associativity")); err == nil {
				index.WaysOfAssociativity, _ = strconv.Atoi(strings.TrimSpace(string(data)))
			}
		case "number_of_sets":
			if data, err := lib.Read(path.Join(dir, "number_of_sets")); err == nil {
				index.NumberOfSets, _ = strconv.Atoi(strings.TrimSpace(string(data)))
			}
		}
	}
	return
}
