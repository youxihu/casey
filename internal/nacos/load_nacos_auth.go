package nacos

import (
	"github.com/youxihu/casey/internal/str"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// LoadNacosAuth 从本地文件加载 Nacos 认证配置
func LoadNacosAuth(filePath string) (*str.AuthConfig, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	configFile, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var wrapper str.AuthConfigWrapper
	err = yaml.Unmarshal(configFile, &wrapper)
	if err != nil {
		return nil, err
	}
	return &wrapper.Auth, nil
}
