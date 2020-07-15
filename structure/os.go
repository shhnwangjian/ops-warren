package structure

type SysData struct {
	Os_type          string              `json:"os_type"`
	Asset_id         string              `json:"asset_id"`
	Hostname         string              `json:"hostname"`
	Basic            Basic               `json:"basic"`
	Ram              []Ram               `json:"ram"`
	Cpu              []Cpu               `json:"cpu"`
	Disk             []Disk              `json:"disk"`
	Nic              []Nic               `json:"nic"`
	Startup          []Startup           `json:"startup"`
	Route            string              `json:"route"`
	Crontab          []Crontab           `json:"crontab"`
	Etc_hosts        string              `json:"etc_hosts"`
	Timestamp        int64               `json:"timestamp"`
	Tag              string              `json:"tag"`
	Server_tag       []map[string]string `json:"server_tag"`
	Agent            Agent               `json:"agent"`
	ServerProperties []map[string]string `json:"server_properties"`
	Kubernetes       Kubernetes          `json:"kubernetes"`
	Kvm              Kvm                 `json:"kvm"`
	Process          []Process           `json:"process"`
}

type Basic struct {
	Attribute      string `json:"attribute"`
	Os_release     string `json:"os_release"`
	Os_name        string `json:"os_name"`
	Os_version     string `json:"os_version"`
	Os_units       string `json:"os_units"`
	Manufacturer   string `json:"manufacturer"`
	Model          string `json:"model"`
	Os_kernel      string `json:"os_kernel"`
	Os_sn          string `json:"os_sn"`
	Sn             string `json:"sn"`
	Ram_slot_total int    `json:"ram_slot_total"`
	Ram_slot_used  int    `json:"ram_slot_used"`
	Max_mem_size   uint32 `json:"max_mem_size"`
	Uptime         string `json:"uptime"`
	Language       string `json:"language"`
	LocalTime      string `json:"local_time"`
	Ipmi           string `json:"ipmi"`
	IpmiMac        string `json:"IpmiMac"`
	BiosVersion    string `json:"BiosVersion"`
}

type Cpu struct {
	Device_id  string `json:"device_id"`   // 系统中逻辑处理核的编号
	PhysicalId string `json:"physical_id"` // 单个CPU的标号
	Core_count uint32 `json:"core_count"`  // 该逻辑核所处CPU的物理核数
	Model      string `json:"model"`       // CPU属于的名字及其编号、标称主频
}

type Capacity struct {
	Volume int    `json:"volume"`
	Unit   string `json:"unit"`
}

type DiskCapacity struct {
	Volume float64 `json:"volume"`
	Unit   string  `json:"unit"`
}

type Ram struct {
	Manufactory   string `json:"manufacturer"`
	DeviceLocator string `json:"device_locator"`
	Capacity      `json:"capacity"`
}

type Disk struct {
	Iface_type   string `json:"iface_type"`
	Slot         uint32 `json:"slot"`
	Sn           string `json:"sn"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
	DiskCapacity `json:"real_capacity"`
}

type Nic struct {
	Macaddress string   `json:"macaddress"`
	Name       string   `json:"name"`
	Duplex     string   `json:"duplex"`
	Ipv4       string   `json:"ipv4"`
	Ipv6       string   `json:"ipv6"`
	Netmask    string   `json:"netmask"`
	GateWay    string   `json:"gateway"`
	Speed      string   `json:"speed"`
	Init       bool     `json:"init"`
	Dns        []string `json:"dns"`
	Vip        []string `json:"vip"`
}

type Startup struct {
	Caption string `json:"caption"`
	Command string `json:"command"`
}

type Crontab struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Agent struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildtime"`
}

type Kubernetes struct {
	IsKubelet    int8 `json:"IsKubelet"`
	IsKubeMaster int8 `json:"IsKubeMaster"`
}

type Kvm struct {
	Enable int8 `json:"Enable"`
}

type Process struct {
	Name  string `json:"Name"`
	Ports string `json:"ports"`
}
