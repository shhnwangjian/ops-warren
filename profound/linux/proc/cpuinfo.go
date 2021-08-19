package proc

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/shhnwangjian/ops-warren/lib"
)

type Cpus struct {
	Count int
	Cpu   map[int]*CpuInfo
}

type CpuInfo struct {
	Processor  int
	PhysicalId int
	CorelId    int
	CpuCores   int
	ModelName  string
	CacheSize  int64 // kb
	ErrMsg     error
}

func CpuInfoInstance() *Cpus {
	return &Cpus{Cpu: make(map[int]*CpuInfo)}
}

func (c *Cpus) GetAllNodeInfo() error {
	cpuStrList, err := c.getCpuInfoList()
	if err != nil {
		return err
	}

	for i := 0; i < len(cpuStrList); i++ {
		c.Cpu[i] = ParseCpuInfo(cpuStrList[i])
		c.Count++
	}
	return nil
}

func (c *Cpus) getCpuInfoList() (cpuStrList []string, err error) {
	f, err := os.Open(procFilePath("cpuinfo"))
	if err != nil {
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	if string(data) == "" {
		return nil, fmt.Errorf("cpuinfo file no data")
	}

	cpuStrList = c.parseContent(string(data))
	if len(cpuStrList) == 0 {
		return nil, fmt.Errorf("parse cpuinfo content error")
	}

	return
}

func (c *Cpus) parseContent(data string) []string {
	return lib.SplitNullString(data)
}

func ParseCpuInfo(data string) *CpuInfo {
	c := &CpuInfo{}
	reader := strings.NewReader("\n" + data)
	scanner := bufio.NewScanner(reader)
	if b := scanner.Scan(); !b {
		c.ErrMsg = scanner.Err()
		return c
	}
	for scanner.Scan() {
		line := scanner.Text()
		lineList := strings.Split(line, ":")
		if len(lineList) < 2 {
			continue
		}
		val := strings.TrimSpace(lineList[1])
		switch strings.TrimSpace(lineList[0]) {
		case "processor":
			c.Processor, _ = strconv.Atoi(val)
		case "core id":
			c.CorelId, _ = strconv.Atoi(val)
		case "physical id":
			c.PhysicalId, _ = strconv.Atoi(val)
		case "cpu cores":
			c.CpuCores, _ = strconv.Atoi(val)
		case "model name":
			c.ModelName = val
		case "cache size":
			c.CacheSize, _ = strconv.ParseInt(strings.TrimSpace(lib.RemoveSuffix(val, "KB")), 10, 64)
		default:
			continue
		}
	}
	return c
}
