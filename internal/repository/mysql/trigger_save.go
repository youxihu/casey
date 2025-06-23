package mysql

import (
	"context"
	"time"

	"github.com/youxihu/casey/internal/data/ent"
)

// InsertTriggerToMySQL 插入一条 shell 执行记录到 MySQL 中
func InsertTriggerToMySQL(
	client *ent.Client,
	executor string,
	hostname string,
	command string,
	response string,
) error {
	ctx := context.Background()

	_, err := client.CaseyTrigger.Create().
		SetExecutor(executor).
		SetHostname(hostname).
		SetCommand(command).
		SetResponse(response).
		SetExecutedAt(time.Now()). // 使用当前时间
		Save(ctx)

	if err != nil {
		return err
	}

	return nil
}
