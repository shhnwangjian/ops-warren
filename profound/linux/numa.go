package linux

import (
	"fmt"

	"github.com/shhnwangjian/ops-warren/profound/linux/proc"
	"github.com/shhnwangjian/ops-warren/profound/linux/sys"
)

/*
NUMA相关的几个概念有Node、Socket、Core 以及 LogicalProcessor。
Socket是一个物理上的概念，指的是主板上的cpu插槽。
Node是一个逻辑上的概念，对应于socket。
Core就是一个物理cpu,一个独立的硬件执行单元。
LogicalProcessor就是超线程的概念，是一个逻辑cpu，共享core上的执行单元。
*/

type INuma interface {
	GetNuma() error
	CpuTopology() string
}

type Numa struct {
	Numa map[int]NumaNode
}

type NumaNode struct {
	Node   NodeInfo
	Socket SocketInfo
}

type NodeInfo struct {
	Name     string
	MemTotal int // MB
	MemFree  int // MB
	CpuList  string
	Distance string
}

type SocketInfo struct {
	Core map[int]CoreInfo
}

type CoreInfo struct {
	Processor map[int]LogicalProcessor
}

type LogicalProcessor struct {
	ModelName string
	Cache     sys.Cache
}

func NumaInstance() INuma {
	return &Numa{Numa: make(map[int]NumaNode)}
}

func (n *Numa) getNumaInfo() (err error) {
	if err = n.getNode(); err != nil {
		return
	}
	return n.getCpuInfo()
}

func (n *Numa) GetNuma() error {
	return n.getNumaInfo()
}

func (n *Numa) CpuTopology() (msg string) {
	if err := n.getNumaInfo(); err != nil {
		msg = err.Error()
		return
	}
	for nodeId, numa := range n.Numa {
		msg += fmt.Sprintf("NUMA node %d \n", nodeId)
		msg += fmt.Sprintf("node %d size: %d MB \n", nodeId, numa.Node.MemTotal)
		msg += fmt.Sprintf("node %d free: %d MB \n", nodeId, numa.Node.MemFree)
		msg += fmt.Sprintf("node %d cpus: %s \n", nodeId, numa.Node.CpuList)
		msg += fmt.Sprintf("node %d distances: %s \n", nodeId, numa.Node.Distance)
		msg += fmt.Sprintf("----socket %d \n", nodeId)
		for coreId, core := range numa.Socket.Core {
			msg += fmt.Sprintf("--------Physical core %d \n", coreId)
			for processorId, processor := range core.Processor {
				msg += fmt.Sprintf("------------Logical Processor %d \n", processorId)
				msg += fmt.Sprintf("------------model name: %s \n", processor.ModelName)
				for id, cache := range processor.Cache.CacheIndex {
					msg += fmt.Sprintf("------------cache index%d, id: %d, size: %d KB, type: %s \n",
						id, cache.Id, cache.Size, cache.Type)
				}
			}
		}
		msg += "---------------------------------------- \n"
		msg += "---------------------------------------- \n"
	}
	return
}

func (n *Numa) getNode() error {
	nodes, err := sys.NodeInfoInstance()
	if err != nil {
		return err
	}

	nodes.GetAllNodeInfo()
	if nodes.Count == 0 {
		return fmt.Errorf("no node data")
	}

	for name, line := range nodes.Node {
		nodeInfo := NodeInfo{
			Name:     name,
			MemFree:  line.MemFree,
			MemTotal: line.MemTotal,
			Distance: line.DistanceStr,
			CpuList:  line.CpuListStr,
		}
		n.Numa[line.Id] = NumaNode{Node: nodeInfo, Socket: SocketInfo{Core: make(map[int]CoreInfo)}}
	}

	return nil
}

func (n *Numa) getCpuCache() (*sys.Cpus, error) {
	cache := sys.CacheInstance()
	if err := cache.GetIndexInfo(); err != nil {
		return nil, err
	}
	return cache, nil
}

func (n *Numa) getCpuInfo() error {
	cpus := proc.CpuInfoInstance()
	if err := cpus.GetAllNodeInfo(); err != nil {
		return err
	}
	if cpus.Count == 0 {
		return fmt.Errorf("no cpuinfo data")
	}

	caches, err := n.getCpuCache()
	if err != nil {
		return err
	}

	for _, cpu := range cpus.Cpu {
		processor := LogicalProcessor{ModelName: cpu.ModelName, Cache: caches.Cache[cpu.Processor]}
		if _, ok := n.Numa[cpu.PhysicalId].Socket.Core[cpu.CorelId]; !ok {
			n.Numa[cpu.PhysicalId].Socket.Core[cpu.CorelId] = CoreInfo{Processor: make(map[int]LogicalProcessor)}
		}
		n.Numa[cpu.PhysicalId].Socket.Core[cpu.CorelId].Processor[cpu.Processor] = processor
	}
	return nil
}
