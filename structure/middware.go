package structure

type TomcatData struct {
	CATALINA_BASE   string `json:"CATALINA_BASE"`
	CATALINA_HOME   string `json:"CATALINA_HOME"`
	CATALINA_TMPDIR string `json:"CATALINA_TMPDIR"`
	CLASSPATH       string `json:"CLASSPATH"`
	Version         string `json:"Version"`
	Built           string `json:"Built"`
	Architecture    string `json:"Architecture"`
	JvmVersion      string `json:"JvmVersion"`
}

type RedisData struct {
	Version         string `json:"redis_version"`     // Redis 服务器版本
	GccVersion      string `json:"gcc_version"`       // 编译 Redis 时所使用的 GCC 版本
	ProcessId       string `json:"process_id"`        // 服务器进程的 PID
	TcpPort         string `json:"tcp_port"`          // TCP/IP 监听端口
	ConfigFile      string `json:"config_file"`       // 配置文件路径
	RedisMode       string `json:"redis_mode"`        // 运行模式 "standalone", "sentinel" or "cluster" 模式：单例、哨兵或是集群
	UptimeInSeconds string `json:"uptime_in_seconds"` // 自 Redis 服务器启动以来，经过的秒数
	Role            string `json:"role"`              // 在主从复制中，充当的角色。如果没有主从复制，单点的，它充当的角色也是master
	Os              string `json:"os"`                // Redis 服务器的宿主操作系统
	Hz              string `json:"hz"`                // 时间事件，单位是赫兹,每10毫秒循环一次
	MasterHost      string `json:"master_host"`       // 主服务器的IP地址
	MasterPort      string `json:"master_port"`       // 主服务器监听的端口号
	RunId           string `json:"run_id"`            // Redis 服务器的随机标识符（用于 Sentinel 和集群）
	ArchBits        string `json:"arch_bits"`         // 架构（32 或 64 位）
	MultiplexingApi string `json:"multiplexing_api"`  // 使用的事件处理机制
	AofEnabled      string `json:"aof_enabled"`       // appenonly是否开启,appendonly为yes则为1,no则为0
	ClusterEnabled  string `json:"cluster_enabled"`   // 表示是否启用Redis集群功能的标志。
}

type RabbitMQInfo struct {
	ManagementVersion string `json:"management_version"`
	RabbitmqVersion   string `json:"rabbitmq_version"`
	ClusterName       string `json:"cluster_name"`
	RatesMode         string `json:"rates_mode"` // 消息速率，默认为basic。
	ErlangVersion     string `json:"erlang_version"`
	ErlangFullVersion string `json:"erlang_full_version"`
	Node              string `json:"node"`
}

type NginxData struct {
	Version            string `json:"Version"`
	Gcc                string `json:"Gcc"`
	OpenSSL            string `json:"OpenSSL"`
	TLSSNIEnabled      int    `json:"TLSSNIEnabled"`
	Configure          string `json:"Configure"`
	WorkerProcesses    string `json:"WorkerProcesses"`    // nginx进程数
	WorkerConnections  string `json:"WorkerConnections"`  // 单个进程最大连接数（最大连接数=连接数*进程数）
	WorkerRlimitNoFile string `json:"WorkerRlimitNoFile"` // nginx进程打开的最多文件描述符数目
}

type HttpdData struct {
	Version      string `json:"Version"` // 版本
	Built        string `json:"Built"`   // 编译时间
	MPM          string `json:"MPM"`     // 多进程处理模块
	MMN          string `json:"MMN"`
	Configure    string `json:"Configure"`
	Threaded     string `json:"Threaded"`
	Forked       string `json:"Forked"`
	Architecture string `json:"Architecture"`
	Module       string `json:"Module"`
}

type MongodData struct {
	Version           string            `json:"version"`           // 版本
	GitVersion        string            `json:"gitVersion"`        // git版本
	Modules           []string          `json:"modules"`           // modules
	Allocator         string            `json:"allocator"`         // 内存分配器
	JavascriptEngine  string            `json:"javascriptEngine"`  // JavaScript引擎
	MaxBsonObjectSize map[string]string `json:"maxBsonObjectSize"` // 最大BSON文档大小
	Openssl           MongodOpensslData `json:"openssl"`           // openssl
	BuildEnvironment  map[string]string `json:"buildEnvironment"`  // 构建环境
	StorageEngines    []string          `json:"storageEngines"`    // 存储引擎
}

type MongodOpensslData struct {
	Running  string `json:"running"`
	Compiled string `json:"compiled"`
}
