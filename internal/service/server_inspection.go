package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/youxihu/casey/internal/str"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
	"time"
)

// 运行远程命令
func runRemoteCommand(session *ssh.Session, command string) (string, error) {
	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("命令 %s 执行失败: %v", command, err)
	}
	return strings.TrimSpace(string(output)), nil
}

// 创建新的 SSH 会话
func newSession(client *ssh.Client) (*ssh.Session, error) {
	if client == nil {
		return nil, fmt.Errorf("SSH 客户端为空")
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	return session, nil
}

// 运行命令并返回默认值（若失败）
func runRemoteCommandOrDefault(client *ssh.Client, command, defaultValue string) string {
	session, err := newSession(client)
	if err != nil {
		return defaultValue
	}
	defer session.Close()
	output, err := runRemoteCommand(session, command)
	if err != nil {
		return defaultValue
	}
	return output
}

// SSH 连接主函数
func sshConnect(ip string, port int, user, password string) (*str.Inspection, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	address := fmt.Sprintf("%s:%d", ip, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("无法连接到 %s: %v", address, err)
	}
	defer client.Close()

	inspection := &str.Inspection{
		Timestamp: time.Now(),
		Ip:        ip,
	}

	// 逐个采集信息
	if err := collectBasicInfo(client, inspection); err != nil {
		fmt.Printf("采集基本信息失败: %v\n", err)
	}
	if err := collectCpuInfo(client, inspection); err != nil {
		fmt.Printf("采集 CPU 信息失败: %v\n", err)
	}
	if err := collectMemoryInfo(client, inspection); err != nil {
		fmt.Printf("采集内存信息失败: %v\n", err)
	}
	if err := collectDiskInfo(client, inspection); err != nil {
		fmt.Printf("采集磁盘信息失败: %v\n", err)
	}
	if err := collectNetworkInfo(client, inspection); err != nil {
		fmt.Printf("采集网络信息失败: %v\n", err)
	}
	if err := collectProcessInfo(client, inspection); err != nil {
		fmt.Printf("采集进程信息失败: %v\n", err)
	}
	if err := collectLoadAverage(client, inspection); err != nil {
		fmt.Printf("采集负载平均值失败: %v\n", err)
	}

	return inspection, nil
}

// 采集基本信息（主机名、操作系统）
func collectBasicInfo(client *ssh.Client, insp *str.Inspection) error {
	session, err := newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	insp.Hostname, err = runRemoteCommand(session, "hostname")
	if err != nil {
		return err
	}

	session, err = newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	insp.Os, err = runRemoteCommand(session, "uname -s")
	return err
}

// 采集 CPU 信息
func collectCpuInfo(client *ssh.Client, insp *str.Inspection) error {
	session, err := newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	output, err := runRemoteCommand(session, "cat /proc/stat | grep '^cpu '")
	if err != nil {
		return err
	}
	fields := strings.Fields(output)
	if len(fields) < 5 {
		return fmt.Errorf("CPU 数据格式错误: %s", output)
	}
	totalCores, _ := strconv.ParseUint(runRemoteCommandOrDefault(client, "nproc", "1"), 10, 64)
	insp.Cpu = str.CpuInfo{
		Total:  totalCores,
		User:   parseFloat(fields[1]),
		System: parseFloat(fields[3]),
		Idle:   parseFloat(fields[4]),
	}
	return nil
}

// 采集内存信息
func collectMemoryInfo(client *ssh.Client, insp *str.Inspection) error {
	session, err := newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	output, err := runRemoteCommand(session, "cat /proc/meminfo")
	if err != nil {
		return err
	}
	lines := strings.Split(output, "\n")
	var memTotal, memFree, swapTotal, swapFree uint64
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			memTotal = parseUint(fields[1]) * 1024 // 转换为字节
		case "MemFree:":
			memFree = parseUint(fields[1]) * 1024
		case "SwapTotal:":
			swapTotal = parseUint(fields[1]) * 1024
		case "SwapFree:":
			swapFree = parseUint(fields[1]) * 1024
		}
	}
	insp.Memory = str.MemInfo{
		Total:     memTotal,
		Used:      memTotal - memFree,
		Free:      memFree,
		SwapTotal: swapTotal,
		SwapUsed:  swapTotal - swapFree,
	}
	return nil
}

// 采集磁盘信息（只收集根分区 / 的信息）
func collectDiskInfo(client *ssh.Client, insp *str.Inspection) error {
	// 获取磁盘空间
	session, err := newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	output, err := runRemoteCommand(session, "df -B1")
	if err != nil {
		return err
	}
	lines := strings.Split(output, "\n")
	var rootDisk str.DiskInfo
	for _, line := range lines[1:] { // 跳过表头
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		if fields[5] == "/" { // 只处理根分区
			total := parseUint(fields[1])
			used := parseUint(fields[2])
			free := parseUint(fields[3])
			rootDisk = str.DiskInfo{
				Path:  fields[5],
				Total: total,
				Used:  used,
				Free:  free,
			}
			break // 找到根分区后退出循环
		}
	}

	// 获取根分区的 IO 统计
	session, err = newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	ioOutput, err := runRemoteCommand(session, "cat /proc/diskstats")
	if err != nil {
		fmt.Printf("无法获取磁盘 IO 统计: %v\n", err)
	} else {
		ioLines := strings.Split(ioOutput, "\n")
		for _, ioLine := range ioLines {
			ioFields := strings.Fields(ioLine)
			if len(ioFields) < 14 {
				continue
			}
			// 假设根分区的设备名可以通过 /proc/diskstats 匹配（简化处理）
			diskName := ioFields[2] // 设备名，如 sda
			if strings.Contains(rootDisk.Path, diskName) || rootDisk.Path == "/" {
				rootDisk.IoRead = parseUint(ioFields[5]) * 512  // 读取扇区数 * 512 字节
				rootDisk.IoWrite = parseUint(ioFields[9]) * 512 // 写入扇区数 * 512 字节
				break
			}
		}
	}

	// 只将根分区信息存入 insp.Disk
	insp.Disk = []str.DiskInfo{rootDisk}
	return nil
}

// 采集网络信息
func collectNetworkInfo(client *ssh.Client, insp *str.Inspection) error {
	// 获取网卡流量
	session, err := newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	netOutput, err := runRemoteCommand(session, "cat /proc/net/dev")
	if err != nil {
		return err
	}
	lines := strings.Split(netOutput, "\n")
	var netInfos []str.NetInfo
	for _, line := range lines[2:] {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		iface := strings.TrimRight(fields[0], ":")
		recv := parseUint(fields[1])
		sent := parseUint(fields[9])
		netInfos = append(netInfos, str.NetInfo{
			Interface: iface,
			Recv:      recv,
			Sent:      sent,
		})
	}
	insp.Network = netInfos

	// 获取 TCP 状态
	session, err = newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	tcpOutput, err := runRemoteCommand(session, "cat /proc/net/tcp")
	if err != nil {
		return err
	}
	for i := range insp.Network {
		var tcpEstab, tcpTimeWait uint64
		tcpLines := strings.Split(tcpOutput, "\n")
		for _, line := range tcpLines[1:] {
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}
			switch fields[3] {
			case "01": // ESTABLISHED
				tcpEstab++
			case "06": // TIME_WAIT
				tcpTimeWait++
			}
		}
		insp.Network[i].TcpEstab = tcpEstab
		insp.Network[i].TcpTimeWait = tcpTimeWait
	}
	return nil
}

// 采集进程信息（包括 TopProcesses）
func collectProcessInfo(client *ssh.Client, insp *str.Inspection) error {
	// 进程总数
	session, err := newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	procCount, err := runRemoteCommand(session, "ps -e | wc -l")
	if err == nil {
		insp.Processes = uint32(parseUint(procCount))
	}

	// 僵尸进程数
	session, err = newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	zombieCount, err := runRemoteCommand(session, "ps -eo state | grep -c '^Z'")
	if err == nil {
		insp.ZombieProcs = uint32(parseUint(zombieCount))
	}

	// 前 5 个高 CPU 进程
	session, err = newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	topOutput, err := runRemoteCommand(session, "ps -eo pid,comm,pcpu,rss --sort=-pcpu | head -n 6") // 跳过表头取前5
	if err != nil {
		fmt.Printf("无法获取高 CPU 进程: %v\n", err)
	} else {
		lines := strings.Split(topOutput, "\n")
		var topProcs []str.ProcessInfo
		for _, line := range lines[1:] { // 跳过表头
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}
			pid, _ := strconv.Atoi(fields[0])
			cpuPercent, _ := strconv.ParseFloat(fields[2], 64)
			memUsage := parseUint(fields[3]) * 1024 // RSS 以 KB 为单位，转换为字节
			topProcs = append(topProcs, str.ProcessInfo{
				Pid:        pid,
				Name:       fields[1],
				CpuPercent: cpuPercent,
				MemUsage:   memUsage,
			})
		}
		insp.TopProcesses = topProcs
	}
	return nil
}

// 采集负载平均值
func collectLoadAverage(client *ssh.Client, insp *str.Inspection) error {
	session, err := newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	output, err := runRemoteCommand(session, "uptime")
	if err != nil {
		return err
	}
	fields := strings.Fields(output)
	if len(fields) >= 10 {
		insp.CpuLoad[0], _ = strconv.ParseFloat(strings.TrimRight(fields[len(fields)-3], ","), 64)
		insp.CpuLoad[1], _ = strconv.ParseFloat(strings.TrimRight(fields[len(fields)-2], ","), 64)
		insp.CpuLoad[2], _ = strconv.ParseFloat(fields[len(fields)-1], 64)
	}
	return nil
}

// 工具函数
func parseFloat(s string) float64 {
	value, _ := strconv.ParseFloat(s, 64)
	return value
}

func parseUint(s string) uint64 {
	value, _ := strconv.ParseUint(s, 10, 64)
	return value
}

func ConnectToServers(configs []*str.Config) []*str.Inspection {
	var results []*str.Inspection
	for _, config := range configs {
		for _, system := range config.System {
			for name, host := range system.Hosts {
				fmt.Printf("尝试 SSH 连接 %s (%s:%d)\n", name, host.Address, host.Port)
				inspection, err := sshConnect(host.Address, host.Port, host.User, host.Passwd)
				if err != nil {
					fmt.Printf("连接 %s 失败: %v\n", name, err)
					continue
				}
				results = append(results, inspection)
			}
		}
	}
	return results
}

// SetupHTTP 使用 Gin 启动 HTTP 服务
func SetupHTTP(configs []*str.Config, addr string) error {
	router := gin.Default()

	// 定义 GET /inspect 接口
	router.GET("/inspect", func(c *gin.Context) {
		// 调用原有函数获取检查结果
		inspections := ConnectToServers(configs)

		// 返回 JSON 响应
		if len(inspections) == 0 {
			c.JSON(200, gin.H{
				"message": "没有成功采集到任何数据",
				"data":    []interface{}{},
			})
		} else {
			c.JSON(200, gin.H{
				"message": "成功采集数据",
				"data":    inspections,
			})
		}
	})

	fmt.Printf("HTTP 服务启动在 %s\n", addr)
	return router.Run(addr)
}
