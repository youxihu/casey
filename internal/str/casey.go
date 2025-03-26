package str

// Host 表示主机配置
type Host struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	User    string `yaml:"user"`
	Passwd  string `yaml:"passwd"`
}

// System 表示系统主机配置
type System struct {
	Hosts map[string]Host `yaml:"hosts"`
}

// MySQL 表示 MySQL 实例配置
type MySQL struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	User    string `yaml:"user"`
	Passwd  string `yaml:"passwd"`
}

// Process 表示进程配置
type Process struct {
	MySQL map[string]MySQL `yaml:"mysql"`
}

// Config 表示整个配置
type Config struct {
	System  []System  `yaml:"system"`
	Process []Process `yaml:"process"`
}
