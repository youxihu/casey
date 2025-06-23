package str

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// LoadAllConfigsInDir 加载目录下所有文件（无论文件名），尝试解析为 YAML 配置
func LoadAllConfigsInDir(dirPath string) ([]*Config, error) {
	var configs []*Config

	// 读取目录下所有文件（不递归子目录）
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v (目录路径: %s)", err, dirPath)
	}

	// 遍历所有文件（包括无扩展名的文件）
	for _, entry := range entries {
		if entry.IsDir() {
			continue // 跳过子目录
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		var config Config
		if err := yaml.Unmarshal(data, &config); err != nil {
			continue
		}
		configs = append(configs, &config)
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("目录 %s 下没有有效的配置文件", dirPath)
	}

	return configs, nil
} 