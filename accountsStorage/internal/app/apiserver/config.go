package apiserver

type Config struct {
	ShutdownTimeout int    `toml:"shutdown_timeout"`
	BindAddres      string `toml:"bind_addres"`
	LogLevel        string `toml:"log_level"`
	DatabaseType    string `toml:"database_type"`
	DatabaseURL     string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
