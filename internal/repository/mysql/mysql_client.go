package mysql

import (
	stdsql "database/sql"
	entSql "entgo.io/ent/dialect/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/youxihu/casey/internal/data/ent"
	"github.com/youxihu/casey/internal/str"
	"log"
)

func NewMySQLClient(dirPath string) (*ent.Client, error) {
	configs, err := str.LoadAllConfigsInDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %v", err)
	}

	for _, config := range configs {
		for _, process := range config.Process {
			if len(process.MySQL) == 0 {
				continue
			}

			// 只处理第一个 MySQL 实例
			for name, mysqlCfg := range process.MySQL {
				dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/aiops?parseTime=True",
					mysqlCfg.User,
					mysqlCfg.Passwd,
					mysqlCfg.Address,
					mysqlCfg.Port,
				)

				log.Printf("尝试连接 MySQL 实例: %s (%s)", name, mysqlCfg.Address)

				sqlDB, err := stdsql.Open("mysql", dsn)
				if err != nil {
					return nil, fmt.Errorf("sql.Open 失败: %v", err)
				}

				if err := sqlDB.Ping(); err != nil {
					_ = sqlDB.Close()
					return nil, fmt.Errorf("Ping 失败: %v", err)
				}

				driver := entSql.OpenDB("mysql", sqlDB)
				client := ent.NewClient(ent.Driver(driver))
				return client, nil
			}

			break
		}
		break
	}

	return nil, fmt.Errorf("没有找到有效的 MySQL 配置")
}
