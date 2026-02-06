package config

type Config struct {
	App struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
		Env  string `toml:"env"`
	} `toml:"app"`

	Etcd struct {
		Endpoints            []string `toml:"endpoints"`
		Username             string   `toml:"username"`
		Password             string   `toml:"password"`
		ReconnectInterval    int      `toml:"reconnect_interval"`     // 初始重连间隔(秒)
		MaxReconnectInterval int      `toml:"max_reconnect_interval"` // 最大重连间隔(秒)
		HealthCheckInterval  int      `toml:"health_check_interval"`  // 健康检查间隔(秒)
		DialTimeout          int      `toml:"dial_timeout"`           // 连接超时(秒)
		CorednsPrefix        string   `toml:"coredns_prefix"`         // CoreDNS etcd key 前缀, 默认 /skydns
	} `toml:"etcd"`

	JWT struct {
		Secret string `toml:"secret"`
		Expiry int64  `toml:"expiry"`
	} `toml:"jwt"`

	Logger struct {
		Level     string `toml:"level"`
		FilePath  string `toml:"file_path"`
		MaxSize   int    `toml:"max_size"`
		MaxBackup int    `toml:"max_backup"`
		MaxAge    int    `toml:"max_age"`
	} `toml:"logger"`
}

var GlobalConfig *Config

func GetConfig() *Config {
	return GlobalConfig
}
