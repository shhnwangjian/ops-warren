package structure

type ProcessData struct {
	Pid                    int32    `json:"pid"`
	Name                   string   `json:"name"`  // 进程名
	State                  string   `json:"state"` // 状态
	PPid                   int32    `json:"ppid"`
	VoluntaryCtxSwitches   int64    `json:"voluntary_ctxt_switches"`
	InvoluntaryCtxSwitches int64    `json:"nonvoluntary_ctxt_switches"`
	Uids                   []string `json:"uids"`       // 用户
	Gids                   []string `json:"gids"`       // 用户组
	NumThreads             int32    `json:"numThreads"` // 线程数
	MemInfo                *MemoryInfoStat
	SigInfo                *SignalInfoStat
	Tgid                   int32  `json:"tgid"`
	CmdLine                string `json:"cmdLine"` // 启动命令
	Exe                    string `json:"exe"`     // 程序路径
	Statm                  *MemoryInfoStatm
	Environ                string  `json:"environ"`
	NumFds                 int     `json:"num_fds"` // FD
	Limits                 string  `json:"limits"`  // Limit
	Ports                  string  `json:"ports"`
	Listen                 string  `json:"listen"`
	CpuPer                 float64 `json:"cpu_percent"` // cpu百分比
}

type OpenFilesStat struct {
	Path string `json:"path"`
	Fd   uint64 `json:"fd"`
}

type MemoryInfoStat struct {
	VMS    uint64 `json:"vms"`    // bytes
	RSS    uint64 `json:"rss"`    // bytes 应用程序正在使用的物理内存的大小
	Data   uint64 `json:"data"`   // bytes
	Stack  uint64 `json:"stack"`  // bytes
	Locked uint64 `json:"locked"` // bytes
	Swap   uint64 `json:"swap"`   // bytes
}

type SignalInfoStat struct {
	PendingProcess uint64 `json:"pending_process"`
	PendingThread  uint64 `json:"pending_thread"`
	Blocked        uint64 `json:"blocked"`
	Ignored        uint64 `json:"ignored"`
	Caught         uint64 `json:"caught"`
}

type MemoryInfoStatm struct {
	Size     uint64 `json:"size"`     // bytes
	Resident uint64 `json:"resident"` // bytes
	Shared   uint64 `json:"shared"`   // bytes
	Text     uint64 `json:"text"`     // bytes
	Lib      uint64 `json:"lib"`      // bytes
	Data     uint64 `json:"data"`     // bytes
	Dirty    uint64 `json:"dirty"`    // bytes
}
