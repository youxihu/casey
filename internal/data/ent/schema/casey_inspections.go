package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"time"
)

// CaseyInspections holds the schema definition for the CaseyInspections entity.
type CaseyInspections struct {
	ent.Schema
}

// Fields of the CaseyInspections.
func (CaseyInspections) Fields() []ent.Field {

	return []ent.Field{

		field.Int64("id").SchemaType(map[string]string{
			dialect.MySQL: "bigintunsigned", // Override MySQL.
		}).Comment("主键").Unique(),

		field.String("timestamp").SchemaType(map[string]string{
			dialect.MySQL: "varchar(32)", // Override MySQL.
		}).Comment("采集时间，格式：YYYY-MM-DD HH:mm:ss"),

		field.String("hostname").SchemaType(map[string]string{
			dialect.MySQL: "varchar(128)", // Override MySQL.
		}).Comment("主机名"),

		field.String("ip").SchemaType(map[string]string{
			dialect.MySQL: "varchar(64)", // Override MySQL.
		}).Comment("IP地址"),

		field.String("os").SchemaType(map[string]string{
			dialect.MySQL: "varchar(64)", // Override MySQL.
		}).Comment("操作系统名称及版本"),

		field.String("uptime").SchemaType(map[string]string{
			dialect.MySQL: "varchar(128)", // Override MySQL.
		}).Comment("系统运行时间，例如 up 1 hour, 49 minutes"),

		field.String("cpu").SchemaType(map[string]string{
			dialect.MySQL: "json",
		}).Comment("CPU信息，包含 total 核心数和 usage 使用率百分比"),

		field.String("cpu_load").SchemaType(map[string]string{
			dialect.MySQL: "json",
		}).Comment("负载平均值，三个 float 组成的数组 [1m, 5m, 15m]"),

		field.String("memory").SchemaType(map[string]string{
			dialect.MySQL: "json",
		}).Comment("内存信息，包含 total/used/free/swapTotal/swapUsed (单位 GB)"),

		field.String("disk").SchemaType(map[string]string{
			dialect.MySQL: "json",
		}).Comment("磁盘信息，数组，每项含 path 挂载点、total/used/free (单位 GB)"),

		field.String("disk_io").SchemaType(map[string]string{
			dialect.MySQL: "json",
		}).Comment("磁盘 IO，包含 readPerSec/writePerSec 单位 KB/s"),

		field.String("router").SchemaType(map[string]string{
			dialect.MySQL: "varchar(64)", // Override MySQL.
		}).Comment("网卡接口名称，如 ens33, em1"),

		field.String("net_stats").SchemaType(map[string]string{
			dialect.MySQL: "json",
		}).Comment("网络统计，包含 download/upload 流量 (单位 MB/s 或 GB/s)"),

		field.Int32("processes").SchemaType(map[string]string{
			dialect.MySQL: "intunsigned", // Override MySQL.
		}).Comment("进程总数"),

		field.Int32("zombie_procs").SchemaType(map[string]string{
			dialect.MySQL: "intunsigned", // Override MySQL.
		}).Comment("僵尸进程数"),

		field.String("top_processes").SchemaType(map[string]string{
			dialect.MySQL: "json",
		}).Comment("占用资源最高的前 N 个进程，包含 pid/name/cpuPercent/memUsage"),

		field.String("env").SchemaType(map[string]string{
			dialect.MySQL: "json",
		}).Optional().Comment("环境变量或其他附加信息（可为空）"),

		field.Time("created_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Default(time.Now).Comment("记录创建时间"),
	}

}

// Edges of the CaseyInspections.
func (CaseyInspections) Edges() []ent.Edge {
	return nil
}

func (CaseyInspections) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "casey_inspections"},
	}
}
