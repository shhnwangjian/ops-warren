package structure

type DownloadImageData struct {
	Url           string `json:"Url" valid:"required"`  // url
	Md5           string `json:"Md5" valid:"required"`  // md5
	Name          string `json:"Name" valid:"required"` // 名称
	AgentIp       string `json:"AgentIp" valid:"required"`
	ServerAddr    string `json:"ServerAddr" valid:"required"`
	RpcServerName string `json:"RpcServerName" valid:"required"`
	UserName      string `json:"UserName" valid:"required"`
}

type CreateVmData struct {
	DomainName    string   `json:"DomainName" valid:"required"`
	Uuid          string   `json:"Uuid" valid:"required"`
	Ip            string   `json:"Ip" valid:"required"`
	Template      string   `json:"Template" valid:"required"`
	TemplateMd5   string   `json:"TemplateMd5" valid:"required"`
	Cpu           int      `json:"Cpu" valid:"required"`
	MaxCpu        int      `json:"MaxCpu" valid:"required"`
	Memory        int      `json:"Memory" valid:"required"`
	MaxMemory     int      `json:"MaxMemory" valid:"required"`
	AgentIp       string   `json:"AgentIp" valid:"required"`
	ServerAddr    string   `json:"ServerAddr" valid:"required"`
	RpcServerName string   `json:"RpcServerName" valid:"required"`
	UserName      string   `json:"UserName" valid:"required"`
	Gateway       string   `json:"Gateway" valid:"required"`
	VlanId        int      `json:"VlanId" valid:"required"`
	Inet          string   `json:"Inet" valid:"required"`
	Brif          string   `json:"Brif" valid:"required"`
	Dns           []string `json:"Dns" valid:"required"`
}

type VmServiceData struct {
	DomainName    string `json:"DomainName" valid:"required"`
	Ops           int8   `json:"Ops" valid:"required"`
	AgentIp       string `json:"AgentIp" valid:"required"`
	ServerAddr    string `json:"ServerAddr" valid:"required"`
	RpcServerName string `json:"RpcServerName" valid:"required"`
	UserName      string `json:"UserName" valid:"required"`
	SetMemory     bool   `json:"SetMemory"`
	Memory        int    `json:"Memory"`
}

type VmCpuMemData struct {
	DomainName    string `json:"DomainName" valid:"required"`
	Cpu           int    `json:"Cpu"`
	Memory        int    `json:"Memory"`
	MaxMemory     int    `json:"MaxMemory"`
	AgentIp       string `json:"AgentIp" valid:"required"`
	ServerAddr    string `json:"ServerAddr" valid:"required"`
	RpcServerName string `json:"RpcServerName" valid:"required"`
	UserName      string `json:"UserName" valid:"required"`
}

type VmAttachDiskData struct {
	DomainName    string `json:"DomainName" valid:"required"`
	Disk          int    `json:"Disk" valid:"required"`
	AgentIp       string `json:"AgentIp" valid:"required"`
	ServerAddr    string `json:"ServerAddr" valid:"required"`
	RpcServerName string `json:"RpcServerName" valid:"required"`
	UserName      string `json:"UserName" valid:"required"`
}

type VmChangeUserPasswordData struct {
	DomainName    string `json:"DomainName" valid:"required"`
	User          string `json:"User" valid:"required"`
	Password      string `json:"Password" valid:"required"`
	AgentIp       string `json:"AgentIp" valid:"required"`
	ServerAddr    string `json:"ServerAddr" valid:"required"`
	RpcServerName string `json:"RpcServerName" valid:"required"`
	UserName      string `json:"UserName" valid:"required"`
}
