package sys

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/shhnwangjian/ops-warren/lib"
)

// 在NUMA架构下，CPU的概念从大到小依次是：Node、Socket、Core、Logical Processor（Node->Socket->Core->Logical Processor）

type Nodes struct {
	Count int
	Node  map[string]*NodeInfo
}

type NodeInfo struct {
	Id          int
	MemTotal    int // MB
	MemFree     int // MB
	CpuListStr  string
	DistanceStr string
	ErrMsg      error
}

func NodeInfoInstance() (*Nodes, error) {
	n := &Nodes{Node: make(map[string]*NodeInfo)}
	if err := n.getNodeNum(); err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Nodes) getNodeNum() (err error) {
	f, err := os.Open(DefaultSysDeviceNode)
	if err != nil {
		return
	}
	defer f.Close()

	list, err := f.Readdir(-1)
	if err != nil {
		return
	}
	if len(list) == 0 {
		return fmt.Errorf("no node dir")
	}

	for i := 0; i < len(list); i++ {
		if !list[i].IsDir() {
			continue
		}
		if !strings.HasPrefix(list[i].Name(), "node") {
			continue
		}
		var num int
		numStr := lib.RemovePrefix(list[i].Name(), "node")
		if numStr != "" {
			num, _ = strconv.Atoi(numStr)
		}
		n.Node[list[i].Name()] = &NodeInfo{Id: num}
		n.Count++
	}
	return
}

func (n *Nodes) getNodeCpuList(nodeName string) error {
	cpuListPath := fmt.Sprintf("/sys/devices/system/node/%s/cpulist", nodeName)
	data, err := ioutil.ReadFile(cpuListPath)
	if err != nil {
		return err
	}
	n.Node[nodeName].CpuListStr = strings.TrimSpace(string(data))
	return nil
}

func (n *Nodes) getNodeDistance(nodeName string) error {
	cpuListPath := fmt.Sprintf("/sys/devices/system/node/%s/distance", nodeName)
	data, err := ioutil.ReadFile(cpuListPath)
	if err != nil {
		return err
	}
	n.Node[nodeName].DistanceStr = strings.TrimSpace(string(data))
	return nil
}

func (n *Nodes) GetAllNodeInfo() {
	for k, _ := range n.Node {
		if err := n.getNodeMemoryInfo(k); err != nil {
			n.Node[k].ErrMsg = err
			continue
		}
		if err := n.getNodeCpuList(k); err != nil {
			n.Node[k].ErrMsg = err
			continue
		}
		if err := n.getNodeDistance(k); err != nil {
			n.Node[k].ErrMsg = err
			continue
		}
	}
}

func (n *Nodes) getNodeMemoryInfo(nodeName string) error {
	memInfoPath := fmt.Sprintf("/sys/devices/system/node/%s/meminfo", nodeName)
	data, err := ioutil.ReadFile(memInfoPath)
	if err != nil {
		return err
	}

	total, free, err := ParseNodeMemInfo(string(data))
	if err != nil {
		return err
	}

	n.Node[nodeName].MemTotal = int(total)
	n.Node[nodeName].MemFree = int(free)
	return nil
}

func ParseNodeMemInfo(data string) (total, free int64, err error) {
	reader := strings.NewReader("\n" + data)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	if b := scanner.Scan(); !b {
		err = scanner.Err()
		return
	}
	for scanner.Scan() {
		line := scanner.Text()
		name := strings.Fields(line)
		if len(name) < 3 {
			continue
		}
		switch name[2] {
		case "MemTotal:":
			total, _ = strconv.ParseInt(name[3], 10, 64)
			total /= 1024
		case "MemFree:":
			free, _ = strconv.ParseInt(name[3], 10, 64)
			free /= 1024
		default:
			continue
		}
	}
	return
}
