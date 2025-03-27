package service

import (
	"fmt"
	"github.com/youxihu/casey/internal/str"
)

// SSH 连接主函数
func runShell(ip string, port int, user, password, cmd string) (string, error) {
	client, err := createSSHClient(ip, port, user, password)
	if err != nil {
		return "", fmt.Errorf("无法连接到 %s: %v", ip, err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建session失败: %v", err)
	}
	defer session.Close()

	output, err := runRemoteCommand(session, cmd)
	if err != nil {
		return "", err // 将详细错误传递给上层
	}
	return output, nil
}

// 根据主机名执行命令
func ConnectToRunShell(configs []*str.Config, hostname, cmd string) (string, error) {
	for _, config := range configs {
		for _, system := range config.System {
			if host, exists := system.Hosts[hostname]; exists {
				output, err := runShell(host.Address, host.Port, host.User, host.Passwd, cmd)
				if err != nil {
					return "", fmt.Errorf("执行 %s 失败: %v", hostname, err)
				}
				return output, nil
			}
		}
	}
	return "", fmt.Errorf("主机 %s 未找到", hostname)
}
