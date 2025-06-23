package service

import (
	"fmt"
	"github.com/youxihu/casey/internal/str"
	"golang.org/x/crypto/ssh"
)

// RunShellOnHost 根据主机名执行命令
func RunShellOnHost(configs []*str.Config, hostname, cmd string) (string, error) {
	return ConnectToRunShell(configs, hostname, cmd)
}

// ConnectToRunShell 根据主机名执行命令
func ConnectToRunShell(configs []*str.Config, hostname, cmd string) (string, error) {
	for _, config := range configs {
		for _, system := range config.System {
			if host, exists := system.Hosts[hostname]; exists {
				return runShell(host.Address, host.Port, host.User, host.Passwd, cmd)
			}
		}
	}
	return "", fmt.Errorf("主机 %s 未找到", hostname)
}

// runShell 远程执行命令（简化版）
func runShell(ip string, port int, user, password, cmd string) (string, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	address := fmt.Sprintf("%s:%d", ip, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return "", err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	output, err := session.CombinedOutput(cmd)
	return string(output), err
} 