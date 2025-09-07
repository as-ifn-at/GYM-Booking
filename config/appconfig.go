package config

type Config struct {
	Port            string
	CacheConfig     CacheConfig
	DBConfigOptions DBConfigOptions
}

func Load() *Config {
	config := &Config{}
	
	return config
}
