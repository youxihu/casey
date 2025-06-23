package service

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/youxihu/casey/internal/data/ent"
	"github.com/youxihu/casey/internal/repository/rds"
	"github.com/youxihu/casey/internal/str"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

// GetAllInspectionsFromCache 获取所有巡检缓存
func GetAllInspectionsFromCache(redisClient *redis.Client) ([]*str.Inspection, error) {
	return rds.GetAllInspectionsFromRedis(redisClient)
}

// GetLatestInspectionFromCache 获取最新巡检缓存（批量）
func GetLatestInspectionFromCache(redisClient *redis.Client) ([]*str.Inspection, error) {
	return rds.GetLatestInspectionFromRedis(redisClient)
}

// StartPeriodicInspection 定时巡检并存储到 Redis,Mysql
func StartPeriodicInspection(configs []*str.Config, redisClient *redis.Client, entClient *ent.Client) {
	if entClient == nil {
		log.Println("[定时巡检] entClient 为 nil，无法写入 MySQL")
	} else {
		log.Println("[定时巡检] entClient 初始化正常")
	}

	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("[定时巡检] 开始巡检...")
		inspections := ConnectToServers(configs)
		err := SaveAllInspectionsToCache(redisClient, entClient, inspections, false)
		if err != nil {
			fmt.Println("[定时巡检] 存储失败:", err)
		} else {
			fmt.Println("[定时巡检] 已写入 Redis")
		}
	}
}

// ConnectToServers 批量SSH采集主机信息
func ConnectToServers(configs []*str.Config) []*str.Inspection {
	var results []*str.Inspection
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, config := range configs {
		for _, system := range config.System {
			for _, host := range system.Hosts {
				h := host // 避免闭包问题
				wg.Add(1)
				go func() {
					defer wg.Done()
					inspection, err := sshConnect(h.Address, h.Port, h.User, h.Passwd, h.Router, h.Device)
					if err == nil {
						mu.Lock()
						results = append(results, roundInspection(inspection))
						mu.Unlock()
					}
				}()
			}
		}
	}
	wg.Wait()
	return results
}

// sshConnect 采集单台主机信息（自动识别主用网卡）
func sshConnect(ip string, port int, user, password string, router string, device string) (*str.Inspection, error) {
	getCmdOut := func(cmd string) string {
		out, err := runShell(ip, port, user, password, cmd)
		if err != nil {
			return ""
		}
		return strings.TrimSpace(out)
	}
	// 自动识别主用网卡（排除lo、docker、veth，优先有IP的网卡，兼容无ip命令主机）
	cfgRouter := router
	now := time.Now().Format(time.DateTime)
	inspection := &str.Inspection{
		Timestamp: now,
		Hostname:  getCmdOut("hostname"),
		Os:        getCmdOut("uname"),
		Ip:        ip,
		Uptime:    getCmdOut("uptime -p"),
		Router:    cfgRouter,
	}
	inspection.Cpu = collectCPUInfo(getCmdOut)
	inspection.CpuLoad = collectCPULoad(getCmdOut)
	inspection.Memory = collectMemInfo(getCmdOut)
	inspection.Disk = collectDiskInfo(getCmdOut)
	inspection.DiskIO = collectDiskIO(getCmdOut, device)
	inspection.NetStats = collectNetStats(getCmdOut, router)
	inspection.Processes, inspection.ZombieProcs = collectProcesses(getCmdOut)
	inspection.TopProcesses = collectTopProcesses(getCmdOut)
	return inspection, nil
}

func collectCPUInfo(getCmdOut func(string) string) str.CpuInfo {
	total := parseUint(getCmdOut("nproc"))
	cpuLine := getCmdOut(`top -bn1 | grep 'Cpu(s)' | sed 's/.*, *\([0-9.]*\)%* id.*/\1/'`)
	idle, _ := strconv.ParseFloat(cpuLine, 64)
	usage := 100 - idle
	usage = round2(usage)
	return str.CpuInfo{Total: total, Usage: usage}
}

func collectCPULoad(getCmdOut func(string) string) [3]float64 {
	var load [3]float64
	loadStr := getCmdOut("cat /proc/loadavg | awk '{print $1,$2,$3}'")
	fmt.Sscanf(loadStr, "%f %f %f", &load[0], &load[1], &load[2])
	return load
}

func collectMemInfo(getCmdOut func(string) string) str.MemInfo {
	memTotal := parseUint(getCmdOut("awk '/MemTotal/ {print $2}' /proc/meminfo"))
	memFree := parseUint(getCmdOut("awk '/MemFree/ {print $2}' /proc/meminfo"))
	swapTotal := parseUint(getCmdOut("awk '/SwapTotal/ {print $2}' /proc/meminfo"))
	swapFree := parseUint(getCmdOut("awk '/SwapFree/ {print $2}' /proc/meminfo"))
	toGB := func(kb uint64) float64 { return float64(kb) / 1024 / 1024 }
	return str.MemInfo{
		Total:     toGB(memTotal),
		Free:      toGB(memFree),
		Used:      toGB(memTotal - memFree),
		SwapTotal: toGB(swapTotal),
		SwapUsed:  toGB(swapTotal - swapFree),
	}
}

func collectDiskInfo(getCmdOut func(string) string) []str.DiskInfo {
	disk := getCmdOut("df -B1 / | awk 'NR==2 {print $2, $3, $4, $6}'")
	var dtotal, dused, dfree uint64
	var dpath string
	fmt.Sscanf(disk, "%d %d %d %s", &dtotal, &dused, &dfree, &dpath)
	toGB := func(b uint64) float64 { return float64(b) / 1024 / 1024 / 1024 }
	return []str.DiskInfo{{Path: dpath, Total: toGB(dtotal), Used: toGB(dused), Free: toGB(dfree)}}
}

func collectDiskIO(getCmdOut func(string) string, deviceName string) str.DiskIO {
	if deviceName == "" {
		return str.DiskIO{}
	}
	// 获取 /proc/diskstats 内容
	diskStats := getCmdOut("cat /proc/diskstats")
	lines := strings.Split(diskStats, "\n")

	var readMB, writeMB float64 = 0, 0
	fmt.Println("device", deviceName)
	for _, line := range lines {
		// 判断是否包含目标设备名（兼容 LVM、NVMe、普通设备）
		if strings.Contains(line, " "+deviceName) || strings.HasPrefix(line, deviceName+" ") {
			fields := strings.Fields(line)
			if len(fields) >= 14 {
				sectorsRead, _ := strconv.ParseFloat(fields[5], 64)  // 读扇区数
				sectorsWrite, _ := strconv.ParseFloat(fields[9], 64) // 写扇区数

				readMB = sectorsRead * 512 / 1024 / 1024 // 转为 MB
				writeMB = sectorsWrite * 512 / 1024 / 1024
			}
			break
		}
	}

	return str.DiskIO{
		ReadPerSec:  round2(readMB),
		WritePerSec: round2(writeMB),
	}
}

func collectNetStats(getCmdOut func(string) string, router string) str.NetStats {
	line := getCmdOut(fmt.Sprintf(`sar -n DEV 1 1 | grep -w %s | tail -1`, router))
	fields := strings.Fields(line)

	var rxKB, txKB float64 = 0, 0

	if len(fields) >= 6 {
		rxKB, _ = strconv.ParseFloat(fields[4], 64) // rxkB/s
		txKB, _ = strconv.ParseFloat(fields[5], 64) // txkB/s
	} else {
		fmt.Println("[网络IO采集][sar] 字段不足，内容：", line)
	}

	return str.NetStats{
		Download: round2(rxKB / 1024), // MB/s
		Upload:   round2(txKB / 1024), // MB/s
	}
}

func collectProcesses(getCmdOut func(string) string) (uint32, uint32) {
	procs := uint32(parseUint(getCmdOut("ps -e | wc -l")))
	zombies := uint32(parseUint(getCmdOut("ps -eo stat | grep -c Z")))
	return procs, zombies
}

func collectTopProcesses(getCmdOut func(string) string) []str.ProcessInfo {
	topProcs := getCmdOut("ps -eo pid,comm,pcpu,rss --sort=-pcpu | head -n 11")
	lines := strings.Split(topProcs, "\n")[1:]
	var result []str.ProcessInfo
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		pid, _ := strconv.Atoi(fields[0])
		cpuPercent, _ := strconv.ParseFloat(fields[2], 64)
		memUsageKB := parseUint(fields[3])
		memUsageGB := float64(memUsageKB) / 1024 / 1024
		memUsageGB = round2(memUsageGB)
		result = append(result, str.ProcessInfo{
			Pid:        pid,
			Name:       fields[1],
			CpuPercent: cpuPercent,
			MemUsage:   memUsageGB,
		})
	}
	return result
}

func parseUint(s string) uint64 {
	n, _ := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	return n
}

// SaveAllInspectionsToCache 存储所有巡检结果到Redis，isManual为true用cache前缀，否则用inspect
func SaveAllInspectionsToCache(
	redisClient *redis.Client,
	entClient *ent.Client,
	ins []*str.Inspection,
	isManual bool,
) error {
	prefix := "inspect"
	if isManual {
		prefix = "cache"
	}
	return rds.SaveAllInspectionsToRedis(redisClient, entClient, ins, prefix)
}

func round2(x float64) float64 {
	return math.Round(x*100) / 100
}

func roundInspection(ins *str.Inspection) *str.Inspection {
	ins.Memory.Total = round2(ins.Memory.Total)
	ins.Memory.Used = round2(ins.Memory.Used)
	ins.Memory.Free = round2(ins.Memory.Free)
	ins.Memory.SwapTotal = round2(ins.Memory.SwapTotal)
	ins.Memory.SwapUsed = round2(ins.Memory.SwapUsed)
	for i := range ins.Disk {
		ins.Disk[i].Total = round2(ins.Disk[i].Total)
		ins.Disk[i].Used = round2(ins.Disk[i].Used)
		ins.Disk[i].Free = round2(ins.Disk[i].Free)
	}
	ins.DiskIO.ReadPerSec = round2(ins.DiskIO.ReadPerSec)
	ins.DiskIO.WritePerSec = round2(ins.DiskIO.WritePerSec)
	ins.NetStats.Download = round2(ins.NetStats.Download)
	ins.NetStats.Upload = round2(ins.NetStats.Upload)
	for i := range ins.CpuLoad {
		ins.CpuLoad[i] = round2(ins.CpuLoad[i])
	}
	for i := range ins.TopProcesses {
		ins.TopProcesses[i].CpuPercent = round2(ins.TopProcesses[i].CpuPercent)
	}
	return ins
}
