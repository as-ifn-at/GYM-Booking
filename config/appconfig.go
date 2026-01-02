package config

type Config struct {
	Port            string
	CacheConfig     CacheConfig
	DBConfigOptions DBConfigOptions
}

func Load() *Config {
	config := &Config{
		Port:        "8080",
		CacheConfig: CacheConfig{},
		DBConfigOptions: DBConfigOptions{
			Host:         "localhost",
			Username:     "user",
			Password:     "",
			DBName:       "gym",
			DBToUse:      "mysql",
			IsLogEnabled: true,
		},
	}

	return config
}
