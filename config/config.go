package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		APP  `yaml:"app"`
		HTTP `yaml:"http"`
		Test `yaml:"test"`
		Log  `yaml:"logger"`
		DB   `yaml:"db"`
	}

	APP struct {
		Name    string `yaml:"name" env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT"`
	}

	Test struct {
		Timeout        time.Duration `yaml:"timeout" env:"TEST_TIMEOUT"`
		BuildDir       string        `yaml:"buildDir" env:"TEST_BUILD_DIR"`
		DockerFilesDir string        `yaml:"dockerFilesDir" env:"TEST_DOCKER_FILES_DIR"`
		RepeatFor      int           `yaml:"repeatFor" env:"TEST_REPEAT_FOR"`
		ExpectedOutput string        `yaml:"expectedOutput" env:"TEST_EXPECTED_OUTPUT"`
	}

	DB struct {
		Driver string `yaml:"driver" env:"DB_DRIVER"`
		Path   string `yaml:"path" env:"DB_PATH"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
		Type  string `env-required:"true" yaml:"type"  env:"LOG_TYPE"`
	}
)

func NewConfig(configPath ...string) (*Config, error) {
	cfg := new(Config)

	path := "./"
	if len(configPath) > 0 {
		path = configPath[0]
	}

	if err := cleanenv.ReadConfig(path+"config.yml", cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	if cfg.DB.Driver == "" {
		cfg.DB.Driver = "sqlite3"
	}
	if cfg.DB.Path == "" {
		cfg.DB.Path = "./tester.db"
	}

	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Type == "" {
		cfg.Log.Type = "text"
	}

	return cfg, nil
}
