package str

type Inspection struct {
	Timestamp    string        `json:"timestamp"` // 采集时间
	Hostname     string        `json:"hostname"`  // 主机名
	Os           string        `json:"os"`        // 操作系统
	Ip           string        `json:"ip"`        // IP 地址
	Uptime       string        `json:"uptime"`
	Cpu          CpuInfo       `json:"cpu"`     // CPU 信息
	CpuLoad      [3]float64    `json:"cpuLoad"` // 负载平均值
	Memory       MemInfo       `json:"memory"`  // 内存信息
	Disk         []DiskInfo    `json:"disk"`    // 磁盘信息（多磁盘）
	DiskIO       DiskIO        `json:"diskIO"`
	Router       string        `json:"router"` // casey结构体中网卡接口
	NetStats     NetStats      `json:"netStats"`
	Processes    uint32        `json:"processes"`    // 进程总数
	ZombieProcs  uint32        `json:"zombieProcs"`  // 僵尸进程数
	TopProcesses []ProcessInfo `json:"topProcesses"` // 前 N 个高占用进程
	Env          []EnvInfo     `json:"env"`
}

// CPU 信息
type CpuInfo struct {
	Total uint64  `json:"total"` // 总核心数
	Usage float64 `json:"usage"` // 总使用率
}

// 内存信息
type MemInfo struct {
	Total     float64 `json:"total"`     // 总内存（GB）
	Used      float64 `json:"used"`      // 已用内存（GB）
	Free      float64 `json:"free"`      // 空闲内存（GB）
	SwapTotal float64 `json:"swapTotal"` // 交换空间总量（GB）
	SwapUsed  float64 `json:"swapUsed"`  // 交换空间已用量（GB）
}

// 磁盘信息
type DiskInfo struct {
	Path  string  `json:"path"`  // 挂载点
	Total float64 `json:"total"` // 总空间（GB）
	Used  float64 `json:"used"`  // 已用空间（GB）
	Free  float64 `json:"free"`  // 空闲空间（GB）
}
type DiskIO struct {
	WritePerSec float64 `json:"writePerSec"` // 单位：KB/s
	ReadPerSec  float64 `json:"readPerSec"`  // 单位：KB/s
}

// 网络信息
type NetInfo struct {
	Interface   string  `json:"interface"`   // 网卡名
	Recv        float64 `json:"recv"`        // 接收字节数
	Sent        float64 `json:"sent"`        // 发送字节数
	TcpEstab    uint32  `json:"tcpEstab"`    // ESTABLISHED 连接数
	TcpTimeWait uint32  `json:"tcpTimeWait"` // TIME_WAIT 连接数
}

type NetStats struct {
	Download float64 `json:"download"` // 网络速率，单位由前端动态显示
	Upload   float64 `json:"upload"`   // 网络速率，单位由前端动态显示
}

// 进程信息
type ProcessInfo struct {
	Pid        int     `json:"pid"`        // 进程 ID
	Name       string  `json:"name"`       // 进程名称
	CpuPercent float64 `json:"cpuPercent"` // CPU 使用率
	MemUsage   float64 `json:"memUsage"`   // 内存使用量（GB）
}

type EnvInfo struct {
	NodeName   string  `json:"nodeName"`
	Status     string  `json:"status"`
	Time       string  `json:"time"`
	Version    string  `json:"version"`
	Url        string  `json:"url"`
	CpuPercent float64 `json:"cpuPercent"`
	MemUsed    uint64  `json:"memUsed"`
	MemTotal   uint64  `json:"memTotal"`
	DiskUsed   uint64  `json:"diskUsed"`
	DiskTotal  uint64  `json:"diskTotal"`
}
