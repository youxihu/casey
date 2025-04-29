package str

type Inspection struct {
	Timestamp    string        `json:"timestamp"`    // 采集时间
	Hostname     string        `json:"hostname"`     // 主机名
	Os           string        `json:"os"`           // 操作系统
	Ip           string        `json:"ip"`           // IP 地址
	Cpu          CpuInfo       `json:"cpu"`          // CPU 信息
	CpuLoad      [3]float64    `json:"cpuLoad"`      // 负载平均值
	Memory       MemInfo       `json:"memory"`       // 内存信息
	Disk         []DiskInfo    `json:"disk"`         // 磁盘信息（多磁盘）
	Router       string        `json:"router"`       // casey结构体中网卡接口
	Network      []NetInfo     `json:"network"`      // 网络信息（多网卡）
	Processes    uint32        `json:"processes"`    // 进程总数
	ZombieProcs  uint32        `json:"zombieProcs"`  // 僵尸进程数
	TopProcesses []ProcessInfo `json:"topProcesses"` // 前 N 个高占用进程
}

// CPU 信息
type CpuInfo struct {
	Total  uint64  `json:"total"`  // 总核心数
	User   float64 `json:"user"`   // 用户态使用率
	System float64 `json:"system"` // 系统态使用率
	Idle   float64 `json:"idle"`   // 空闲率
	Usage  float64 `json:"usage"`  // 总使用率
}

// 内存信息
type MemInfo struct {
	Total     uint64 `json:"total"`     // 总内存（字节）
	Used      uint64 `json:"used"`      // 已用内存（字节）
	Free      uint64 `json:"free"`      // 空闲内存（字节）
	SwapTotal uint64 `json:"swapTotal"` // 交换空间总量（字节）
	SwapUsed  uint64 `json:"swapUsed"`  // 交换空间已用量（字节）
}

// 磁盘信息
type DiskInfo struct {
	Path    string `json:"path"`    // 挂载点
	Total   uint64 `json:"total"`   // 总空间（字节）
	Used    uint64 `json:"used"`    // 已用空间（字节）
	Free    uint64 `json:"free"`    // 空闲空间（字节）
	IoRead  uint64 `json:"ioRead"`  // 累计读取（字节）
	IoWrite uint64 `json:"ioWrite"` // 累计写入（字节）
}

// 网络信息
type NetInfo struct {
	Interface   string `json:"interface"`   // 网卡名
	Sent        uint64 `json:"sent"`        // 发送字节数
	Recv        uint64 `json:"recv"`        // 接收字节数
	TcpEstab    uint64 `json:"tcpEstab"`    // ESTABLISHED 连接数
	TcpTimeWait uint64 `json:"tcpTimeWait"` // TIME_WAIT 连接数
}

// 进程信息
type ProcessInfo struct {
	Pid        int     `json:"pid"`        // 进程 ID
	Name       string  `json:"name"`       // 进程名称
	CpuPercent float64 `json:"cpuPercent"` // CPU 使用率
	MemUsage   uint64  `json:"memUsage"`   // 内存使用量（字节）
}
