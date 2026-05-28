package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration values read from config.yaml.
type Config struct {
	ServerPort string `yaml:"-"`

	DBHost     string `yaml:"-"`
	DBPort     string `yaml:"-"`
	DBUser     string `yaml:"-"`
	DBPassword string `yaml:"-"`
	DBName     string `yaml:"-"`

	RedisAddr     string `yaml:"-"`
	RedisPassword string `yaml:"-"`
	RedisDB       int    `yaml:"-"`

	AWSRegion          string `yaml:"-"`
	AWSAccessKeyID     string `yaml:"-"`
	AWSSecretAccessKey string `yaml:"-"`
	S3Bucket           string `yaml:"-"`

	ImageMaxWidth    int `yaml:"-"`
	ImageJPEGQuality int `yaml:"-"`
}

// rawConfig mirrors the YAML structure for parsing.
type rawConfig struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	AWS struct {
		Region          string `yaml:"region"`
		AccessKeyID     string `yaml:"access_key_id"`
		SecretAccessKey string `yaml:"secret_access_key"`
		S3Bucket        string `yaml:"s3_bucket"`
	} `yaml:"aws"`
	Image struct {
		MaxWidth    int `yaml:"max_width"`
		JPEGQuality int `yaml:"jpeg_quality"`
	} `yaml:"image"`
}

// LoadConfig reads configuration from config.yaml.
func LoadConfig(configPath string) Config {
	raw, err := loadRawConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Config] Failed to load %s: %v\n", configPath, err)
		fmt.Fprintln(os.Stderr, "[Config] Falling back to defaults / environment variables")
		return fallbackConfig()
	}

	return Config{
		ServerPort:         fmt.Sprintf("%d", raw.Server.Port),
		DBHost:             raw.Database.Host,
		DBPort:             fmt.Sprintf("%d", raw.Database.Port),
		DBUser:             raw.Database.User,
		DBPassword:         raw.Database.Password,
		DBName:             raw.Database.Name,
		RedisAddr:          raw.Redis.Addr,
		RedisPassword:      raw.Redis.Password,
		RedisDB:            raw.Redis.DB,
		AWSRegion:          raw.AWS.Region,
		AWSAccessKeyID:     raw.AWS.AccessKeyID,
		AWSSecretAccessKey: raw.AWS.SecretAccessKey,
		S3Bucket:           raw.AWS.S3Bucket,
		ImageMaxWidth:      defaultInt(raw.Image.MaxWidth, 1920),
		ImageJPEGQuality:   defaultInt(raw.Image.JPEGQuality, 85),
	}
}

func loadRawConfig(path string) (rawConfig, error) {
	var raw rawConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return raw, err
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return raw, err
	}
	return raw, nil
}

// fallbackConfig reads from environment variables for backward compatibility.
func fallbackConfig() Config {
	return Config{
		ServerPort:         envOrDefault("SERVER_PORT", "18080"),
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		RedisAddr:          envOrDefault("REDIS_ADDR", "localhost:6379"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		RedisDB:            envIntOrDefault("REDIS_DB", 0),
		AWSRegion:          os.Getenv("AWS_REGION"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3Bucket:           os.Getenv("S3_BUCKET"),
		ImageMaxWidth:      1920,
		ImageJPEGQuality:   85,
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	var i int
	if _, err := fmt.Sscanf(v, "%d", &i); err != nil {
		return fallback
	}
	return i
}

func defaultInt(v, fallback int) int {
	if v == 0 {
		return fallback
	}
	return v
}
