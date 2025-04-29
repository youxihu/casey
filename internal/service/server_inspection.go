package service

import (
	"context"
	"fmt"
	"github.com/youxihu/casey/internal/str"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
	"math"
	"strconv"
	"strings"
	"time"
)

// SSH 连接主函数
func sshConnect(ip string, port int, user, password string, router string) (*str.Inspection, error) {
	client, err := createSSHClient(ip, port, user, password)
	if err != nil {
		return nil, fmt.Errorf("无法连接到 %s: %v", ip, err)
	}
	defer client.Close()

	inspection := &str.Inspection{
		Timestamp: time.Now().Format(time.DateTime),
		Ip:        ip,
		Router:    router,
	}

	// 使用通用函数执行采集任务
	tasks := []func(*ssh.Client, *str.Inspection) error{
		collectBasicInfo,
		collectCpuInfo,
		collectMemoryInfo,
		collectDiskInfo,
		collectNetworkInfo,
		collectProcessInfo,
		collectLoadAverage,
	}

	for _, task := range tasks {
		if err := task(client, inspection); err != nil {
			fmt.Printf("采集信息失败: %v\n", err)
		}
	}

	return inspection, nil
}

// 通用命令执行和解析函数
func runCommandAndParse(client *ssh.Client, command string, parser func(string) error) error {
	session, err := newSession(client)
	if err != nil {
		return err
	}
	defer session.Close()

	output, err := runRemoteCommand(session, command)
	if err != nil {
		return err
	}
	return parser(output)
}

// 采集基本信息（主机名、操作系统）
func collectBasicInfo(client *ssh.Client, insp *str.Inspection) error {
	err := runCommandAndParse(client, "hostname", func(output string) error {
		insp.Hostname = output
		return nil
	})
	if err != nil {
		return err
	}

	return runCommandAndParse(client, "uname -s", func(output string) error {
		insp.Os = output
		return nil
	})
}

func collectCpuInfo(client *ssh.Client, insp *str.Inspection) error {
	// Helper function to parse /proc/stat and extract CPU fields
	parseCpuStats := func() (map[string]float64, error) {
		var fields []string
		err := runCommandAndParse(client, "cat /proc/stat | grep '^cpu '", func(output string) error {
			fields = strings.Fields(output)
			return nil
		})
		if err != nil || len(fields) < 5 {
			return nil, fmt.Errorf("CPU 数据格式错误")
		}

		stats := map[string]float64{
			"user":   parseFloat(fields[1]),
			"nice":   parseFloat(fields[2]),
			"system": parseFloat(fields[3]),
			"idle":   parseFloat(fields[4]),
			"iowait": parseFloat(fields[5]),
		}
		return stats, nil
	}

	// First sample
	prevStats, err := parseCpuStats()
	if err != nil {
		return err
	}

	// Wait for 1 second
	time.Sleep(1 * time.Second)

	// Second sample
	currStats, err := parseCpuStats()
	if err != nil {
		return err
	}

	// Calculate diffs
	diff := map[string]float64{}
	for key := range prevStats {
		diff[key] = currStats[key] - prevStats[key]
	}

	// Total CPU time diff
	totalDiff := diff["user"] + diff["nice"] + diff["system"] + diff["idle"] + diff["iowait"]

	// Calculate CPU usage
	cpuUsage := (totalDiff - diff["idle"] - diff["iowait"]) / totalDiff * 100

	// Get CPU core count
	totalCores, _ := strconv.ParseUint(runRemoteCommandOrDefault(client, "nproc", "1"), 10, 64)

	// Convert jiffies to minutes (HZ=1000, 1 jiffy = 1ms)
	const hzToMinutes = 1000 * 60 // HZ=1000, convert to minutes

	// Populate Inspection struct
	insp.Cpu = str.CpuInfo{
		Total:  totalCores,
		User:   roundToThreeDecimalPlaces(diff["user"] / hzToMinutes),
		System: roundToThreeDecimalPlaces(diff["system"] / hzToMinutes),
		Idle:   roundToThreeDecimalPlaces(diff["idle"] / hzToMinutes),
		Usage:  roundToThreeDecimalPlaces(cpuUsage), // CPU 使用率不需要单位转换
	}
	return nil
}

func collectMemoryInfo(client *ssh.Client, insp *str.Inspection) error {
	memTotal, memFree, swapTotal, swapFree := uint64(0), uint64(0), uint64(0), uint64(0)
	err := runCommandAndParse(client, "cat /proc/meminfo", func(output string) error {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			switch fields[0] {
			case "MemTotal:":
				memTotal = parseUint(fields[1]) * 1024
			case "MemFree:":
				memFree = parseUint(fields[1]) * 1024
			case "SwapTotal:":
				swapTotal = parseUint(fields[1]) * 1024
			case "SwapFree:":
				swapFree = parseUint(fields[1]) * 1024
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 转换为 MB
	insp.Memory = str.MemInfo{
		Total:     bytesToMB(memTotal),
		Used:      bytesToMB(memTotal - memFree),
		Free:      bytesToMB(memFree),
		SwapTotal: bytesToMB(swapTotal),
		SwapUsed:  bytesToMB(swapTotal - swapFree),
	}
	return nil
}

func collectDiskInfo(client *ssh.Client, insp *str.Inspection) error {
	var rootDisk str.DiskInfo
	err := runCommandAndParse(client, "df -B1", func(output string) error {
		lines := strings.Split(output, "\n")
		for _, line := range lines[1:] { // 跳过表头
			fields := strings.Fields(line)
			if len(fields) < 6 || fields[5] != "/" {
				continue
			}
			rootDisk = str.DiskInfo{
				Path:  fields[5],
				Total: parseUint(fields[1]),
				Used:  parseUint(fields[2]),
				Free:  parseUint(fields[3]),
			}
			break
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = runCommandAndParse(client, "cat /proc/diskstats", func(output string) error {
		ioLines := strings.Split(output, "\n")
		for _, ioLine := range ioLines {
			ioFields := strings.Fields(ioLine)
			if len(ioFields) < 14 {
				continue
			}
			diskName := ioFields[2]
			if strings.Contains(rootDisk.Path, diskName) || rootDisk.Path == "/" {
				rootDisk.IoRead = bytesToMB(parseUint(ioFields[5]) * 512)  // 转换为 MB
				rootDisk.IoWrite = bytesToMB(parseUint(ioFields[9]) * 512) // 转换为 MB
				break
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 转换为 MB
	rootDisk.Total = bytesToMB(rootDisk.Total)
	rootDisk.Used = bytesToMB(rootDisk.Used)
	rootDisk.Free = bytesToMB(rootDisk.Free)

	insp.Disk = []str.DiskInfo{rootDisk}
	return nil
}

func collectNetworkInfo(client *ssh.Client, insp *str.Inspection) error {
	cmd := fmt.Sprintf("cat /proc/net/dev | grep -w %s", insp.Router)
	var netInfos []str.NetInfo
	err := runCommandAndParse(client, cmd, func(output string) error {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 10 {
				continue
			}
			iface := strings.TrimRight(fields[0], ":")
			netInfos = append(netInfos, str.NetInfo{
				Interface: iface,
				Recv:      bytesToMB(parseUint(fields[1])), // 转换为 MB
				Sent:      bytesToMB(parseUint(fields[9])), // 转换为 MB
			})
		}
		return nil
	})
	if err != nil {
		return err
	}
	insp.Network = netInfos

	err = runCommandAndParse(client, "cat /proc/net/tcp", func(output string) error {
		for i := range insp.Network {
			tcpLines := strings.Split(output, "\n")
			for _, line := range tcpLines[1:] {
				fields := strings.Fields(line)
				if len(fields) < 4 {
					continue
				}
				switch fields[3] {
				case "01": // ESTABLISHED
					insp.Network[i].TcpEstab++
				case "06": // TIME_WAIT
					insp.Network[i].TcpTimeWait++
				}
			}
		}
		return nil
	})
	return err
}

func collectProcessInfo(client *ssh.Client, insp *str.Inspection) error {
	// 获取总进程数
	err := runCommandAndParse(client, "ps -e | wc -l", func(output string) error {
		insp.Processes = uint32(parseUint(output))
		return nil
	})
	if err != nil {
		return err
	}

	// 获取僵尸进程数
	err = runCommandAndParse(client, "ps -eo state | grep -c '^Z'", func(output string) error {
		// 如果命令返回空字符串或错误，说明没有僵尸进程
		count, _ := strconv.ParseUint(strings.TrimSpace(output), 10, 64)
		insp.ZombieProcs = uint32(count)
		return nil
	})
	if err != nil && !strings.Contains(err.Error(), "status 1") {
		// 只有当错误不是 "status 1" 时才返回错误
		return err
	}

	// 获取高 CPU 进程
	return runCommandAndParse(client, "ps -eo pid,comm,pcpu,rss --sort=-pcpu | head -n 11", func(output string) error {
		lines := strings.Split(output, "\n")
		var topProcs []str.ProcessInfo
		for _, line := range lines[1:] { // 跳过表头
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}
			pid, _ := strconv.Atoi(fields[0])
			cpuPercent, _ := strconv.ParseFloat(fields[2], 64)
			memUsage := parseUint(fields[3]) * 1024
			topProcs = append(topProcs, str.ProcessInfo{
				Pid:        pid,
				Name:       fields[1],
				CpuPercent: cpuPercent,
				MemUsage:   memUsage,
			})
		}
		insp.TopProcesses = topProcs
		return nil
	})
}

func collectLoadAverage(client *ssh.Client, insp *str.Inspection) error {
	return runCommandAndParse(client, "uptime", func(output string) error {
		fields := strings.Fields(output)
		if len(fields) >= 10 {
			insp.CpuLoad[0], _ = strconv.ParseFloat(strings.TrimRight(fields[len(fields)-3], ","), 64)
			insp.CpuLoad[1], _ = strconv.ParseFloat(strings.TrimRight(fields[len(fields)-2], ","), 64)
			insp.CpuLoad[2], _ = strconv.ParseFloat(fields[len(fields)-1], 64)
		}
		return nil
	})
}

func ConnectToServers(configs []*str.Config) []*str.Inspection {
	var results []*str.Inspection
	var eg, _ = errgroup.WithContext(context.Background())

	for _, config := range configs {
		for _, system := range config.System {
			for name, host := range system.Hosts {
				eg.Go(func() error {
					inspection, err := sshConnect(host.Address, host.Port, host.User, host.Passwd, host.Router)
					if err != nil {
						fmt.Printf("连接 %s 失败: %v\n", name, err)
						return err
					}
					results = append(results, inspection)
					return nil
				})
			}
		}
	}
	err := eg.Wait()
	if err != nil {
		return nil
	}
	return results
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
func bytesToMB(bytes uint64) uint64 {
	return bytes / 1048576
}

func roundToThreeDecimalPlaces(value float64) float64 {
	return math.Round(value*1000) / 1000
}
