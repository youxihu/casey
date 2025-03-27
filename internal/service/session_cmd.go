package service

import (
	"fmt"
	"golang.org/x/crypto/ssh"
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

// 创建 SSH 连接
func createSSHClient(ip string, port int, user, password string) (*ssh.Client, error) {
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

	return client, nil
}
