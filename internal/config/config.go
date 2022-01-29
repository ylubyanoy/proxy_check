package config

// Config struct
type Config struct {
	DB_NAME     string `yaml:"DB_NAME"`
	DB_USER     string `yaml:"DB_USER"`
	DB_PASSWORD string `yaml:"DB_PASSWORD"`
	DB_HOST     string `yaml:"DB_HOST"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		DB_NAME:     "",
		DB_USER:     "",
		DB_PASSWORD: "",
		DB_HOST:     "localhost",
	}
}
