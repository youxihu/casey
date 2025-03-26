package service

import (
	"fmt"
	"github.com/youxihu/casey/internal/str"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// LoadAllConfigsInDir 加载目录下所有文件（无论文件名），尝试解析为 YAML 配置
func LoadAllConfigsInDir(dirPath string) ([]*str.Config, error) {
	var configs []*str.Config

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
			fmt.Printf("警告: 无法读取文件 %s (跳过): %v\n", fullPath, err)
			continue
		}

		var config str.Config
		if err := yaml.Unmarshal(data, &config); err != nil {
			fmt.Printf("警告: 文件 %s 不是有效的 YAML (跳过): %v\n", fullPath, err)
			continue
		}
		configs = append(configs, &config)
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("目录 %s 下没有有效的配置文件", dirPath)
	}

	return configs, nil

}
