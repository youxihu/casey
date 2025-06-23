package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/youxihu/casey/internal/data/ent"
	"github.com/youxihu/casey/internal/str"
)

// InsertInspectionsToMySQL 将 str.Inspection 数据转换为 ent 可接受格式并批量插入 MySQL
func InsertInspectionsToMySQL(client *ent.Client, insList []*str.Inspection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var bulk []*ent.CaseyInspectionsCreate

	for _, ins := range insList {

		cpuJSON, err := json.Marshal(ins.Cpu)
		if err != nil {
			return fmt.Errorf("failed to marshal cpu info: %v", err)
		}

		memoryJSON, err := json.Marshal(ins.Memory)
		if err != nil {
			return fmt.Errorf("failed to marshal memory info: %v", err)
		}

		diskJSON, err := json.Marshal(ins.Disk)
		if err != nil {
			return fmt.Errorf("failed to marshal disk info: %v", err)
		}

		networkJSON, err := json.Marshal(ins.NetStats)
		if err != nil {
			return fmt.Errorf("failed to marshal network info: %v", err)
		}

		topProcessesJSON, err := json.Marshal(ins.TopProcesses)
		if err != nil {
			return fmt.Errorf("failed to marshal top processes: %v", err)
		}

		cpuLoadJSON, err := json.Marshal(ins.CpuLoad)
		if err != nil {
			return fmt.Errorf("failed to marshal CPU load: %v", err)
		}

		diskIOJSON, err := json.Marshal(ins.DiskIO)
		if err != nil {
			return fmt.Errorf("failed to marshal disk IO: %v", err)
		}

		layout := "2006-01-02 15:04:05"
		ts, err := time.Parse(layout, ins.Timestamp)
		if err != nil {
			log.Printf("时间解析失败: %v", err)
			continue
		}
		bulk = append(bulk, client.CaseyInspections.Create().
			SetTimestamp(ins.Timestamp).
			SetHostname(ins.Hostname).
			SetIP(ins.Ip).
			SetOs(ins.Os).
			SetUptime(ins.Uptime).
			SetCPU(string(cpuJSON)).
			SetMemory(string(memoryJSON)).
			SetDisk(string(diskJSON)).
			SetNetStats(string(networkJSON)).
			SetTopProcesses(string(topProcessesJSON)).
			SetCPULoad(string(cpuLoadJSON)).
			SetDiskIo(string(diskIOJSON)).
			SetRouter(ins.Router).
			SetProcesses(int32(ins.Processes)).
			SetZombieProcs(int32(ins.ZombieProcs)).
			SetCreatedAt(ts))
	}

	if err := client.CaseyInspections.CreateBulk(bulk...).OnConflict().DoNothing().Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert inspections into MySQL: %v", err)
	}

	return nil
}
