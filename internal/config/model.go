package config

type Config struct {
	App struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
		Env  string `toml:"env"`
	} `toml:"app"`

	Etcd struct {
		Endpoints []string `toml:"endpoints"`
		Username  string   `toml:"username"`
		Password  string   `toml:"password"`
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
