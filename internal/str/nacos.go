package str

// AuthConfigWrapper 用于解析带有 auth 顶层键的 YAML
type AuthConfigWrapper struct {
	Auth AuthConfig `yaml:"auth"`
}

// AuthConfig 表示 Nacos 的认证和配置信息
type AuthConfig struct {
	Host        string `yaml:"host"`
	Port        uint64 `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	NamespaceID string `yaml:"namespace_id"`
	TimeoutMS   uint64 `yaml:"timeout_ms"`
	LogDir      string `yaml:"log_dir"`
	CacheDir    string `yaml:"cache_dir"`
	LogLevel    string `yaml:"log_level"`
	DataID      string `yaml:"data_id"`
	Group       string `yaml:"group"`
}
