package nacos

import (
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/youxihu/casey/internal/str"
)

// CreateNacosClient 创建 Nacos 配置客户端并获取远程配置
func CreateNacosClient(config *str.AuthConfig) (string, error) {
	if config.Host == "" || config.Port == 0 {
		return "", errors.New("host or port is empty")
	}

	// 创建服务端配置
	sc := []constant.ServerConfig{
		{
			IpAddr:      config.Host,
			Port:        config.Port,
			ContextPath: "/nacos",
			Scheme:      "http",
		},
	}

	// 创建客户端配置
	cc := constant.ClientConfig{
		NamespaceId:         config.NamespaceID,
		TimeoutMs:           config.TimeoutMS,
		LogDir:              config.LogDir,
		CacheDir:            config.CacheDir,
		LogLevel:            config.LogLevel,
		Username:            config.Username,
		Password:            config.Password,
		NotLoadCacheAtStart: true,
	}

	// 创建配置客户端
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create config client: %v", err)
	}

	// 获取远程配置
	content, err := getConfig(client, config.DataID, config.Group)
	if err != nil {
		return "", fmt.Errorf("failed to get config: %v", err)
	}

	return content, nil
}

// getConfig 获取远程配置的具体实现
func getConfig(client config_client.IConfigClient, dataID, group string) (string, error) {
	if dataID == "" || group == "" {
		return "", errors.New("dataID or group cannot be empty")
	}

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
	})
	if err != nil {
		return "", err
	}
	return content, nil
}
