package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Type aliases for compatibility with existing code
type PGConnectionConfig = struct {
	Host                  string `mapstructure:"host"`
	Port                  string `mapstructure:"port"`
	Name                  string `mapstructure:"name"`
	User                  string `mapstructure:"user"`
	Password              string `mapstructure:"password"`
	Timezone              string `mapstructure:"timezone"`
	SSLMode               string `mapstructure:"sslmode"`
	MaxConnectionLifetime string `mapstructure:"max_connection_lifetime"`
	MaxIdleConnection     int    `mapstructure:"max_idle_connection"`
	MaxOpenConnection     int    `mapstructure:"max_open_connection"`
}

type JWTConfig = struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"`
}

type CORSConfig = struct {
	Enabled          bool   `mapstructure:"enabled"`
	AllowCredentials bool   `mapstructure:"allow_credentials"`
	AllowedHeaders   string `mapstructure:"allowed_headers"`
	AllowedMethods   string `mapstructure:"allowed_methods"`
	AllowedOrigins   string `mapstructure:"allowed_origins"`
	MaxAgeSeconds    int    `mapstructure:"max_age_seconds"`
}

// Config holds all configuration for our application
type Config struct {
	Server struct {
		Host                         string `mapstructure:"host"`
		Port                         string `mapstructure:"port"`
		Env                          string `mapstructure:"env"`
		LogLevel                     string `mapstructure:"log_level"`
		ShutdownCleanupPeriodSeconds int    `mapstructure:"shutdown_cleanup_period_seconds"`
		ShutdownGracePeriodSeconds   int    `mapstructure:"shutdown_grace_period_seconds"`
	} `mapstructure:"server"`

	Database struct {
		PG struct {
			Read  PGConnectionConfig `mapstructure:"read"`
			Write PGConnectionConfig `mapstructure:"write"`
		} `mapstructure:"pg"`
	} `mapstructure:"db"`

	JWT JWTConfig `mapstructure:"jwt"`

	Logger struct {
		Level  string `mapstructure:"log_level"`
		Format string `mapstructure:"format"` // json, console
	} `mapstructure:"server"`

	App struct {
		Name        string `mapstructure:"name"`
		Version     string `mapstructure:"version"`
		Description string `mapstructure:"description"`
		URL         string `mapstructure:"url"`
		Timezone    string `mapstructure:"timezone"`
		Locale      string `mapstructure:"locale"`
		Pagination  struct {
			DefaultLimit int  `mapstructure:"default_limit"`
			MaxLimit     int  `mapstructure:"max_limit"`
			Enabled      bool `mapstructure:"enabled"`
		} `mapstructure:"pagination"`
		CORS CORSConfig `mapstructure:"cors"`
	} `mapstructure:"app"`

	Mail struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Workers  int    `mapstructure:"workers"`
	} `mapstructure:"mail"`

	Cache struct {
		Redis struct {
			Primary struct {
				Host     string `mapstructure:"host"`
				Port     int    `mapstructure:"port"`
				Password string `mapstructure:"password"`
				DB       int    `mapstructure:"db"`
			} `mapstructure:"primary"`
		} `mapstructure:"redis"`
	} `mapstructure:"cache"`

	External struct {
		S3 struct {
			BucketName string `mapstructure:"bucket_name"`
			Region     string `mapstructure:"region"`
			BaseURL    string `mapstructure:"base_url"`
			AccessKey  string `mapstructure:"access_key"`
			Secret     string `mapstructure:"secret"`
		} `mapstructure:"s3"`
	} `mapstructure:"external"`

	Encryption struct {
		Key struct {
			Users string `mapstructure:"users"`
		} `mapstructure:"key"`
	} `mapstructure:"encryption"`
}

// Load loads configuration from YAML file with fallback to environment variables
func Load() (*Config, error) {
	var cfg Config

	// Set defaults
	viper.SetDefault("server.host", "0.0.0.0") // Bind to all interfaces by default
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.env", "development")
	viper.SetDefault("server.log_level", "info")
	viper.SetDefault("server.shutdown.cleanup_period_seconds", 3)
	viper.SetDefault("server.shutdown.grace_period_seconds", 3)

	// App defaults
	viper.SetDefault("app.name", "KelasGo")
	viper.SetDefault("app.version", "0.1.0")
	viper.SetDefault("app.description", "KelasGo API Service")
	viper.SetDefault("app.url", "http://localhost:8080")
	viper.SetDefault("app.timezone", "UTC")
	viper.SetDefault("app.locale", "en_US")
	viper.SetDefault("app.pagination.default_limit", 10)
	viper.SetDefault("app.pagination.max_limit", 100)
	viper.SetDefault("app.pagination.enabled", true)

	viper.SetDefault("app.cors.enabled", true)
	viper.SetDefault("app.cors.allow_credentials", true)
	viper.SetDefault("app.cors.allowed_headers", "Accept,Authorization,Content-Type")
	viper.SetDefault("app.cors.allowed_methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	viper.SetDefault("app.cors.allowed_origins", "http://localhost:8080,http://127.0.0.1:8080")
	viper.SetDefault("app.cors.max_age_seconds", 300)

	viper.SetDefault("mail.host", "smtp.gmail.com")
	viper.SetDefault("mail.port", 587)
	viper.SetDefault("mail.workers", 5)

	viper.SetDefault("cache.redis.primary.host", "localhost")
	viper.SetDefault("cache.redis.primary.port", 6379)
	viper.SetDefault("cache.redis.primary.db", 1)

	viper.SetDefault("jwt.expire_time", 24)

	// Read from YAML config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")

	configLoaded := false
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("⚠️  config.yaml not found. Using environment variables and defaults only.")
			fmt.Println("💡 Create config.yaml from config.example.yaml for full configuration")
		} else {
			return nil, fmt.Errorf("error reading config.yaml file: %w", err)
		}
	} else {
		configLoaded = true
		fmt.Println("✅ Loaded configuration from config.yaml")
	}

	// Configure viper to handle environment variables with underscore replacement
	// This allows reading system environment variables like DB_PG_READ_HOST
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Unmarshal into struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Set JWT expire time manually if not set
	if cfg.JWT.ExpireTime == 0 {
		cfg.JWT.ExpireTime = 24
	}

	// Set logger format
	cfg.Logger.Format = "json" // Default to JSON format
	if cfg.Server.Env == "development" {
		cfg.Logger.Format = "console"
	}

	// Copy server log level to logger level for backward compatibility
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = cfg.Server.LogLevel
	}

	if configLoaded {
		fmt.Printf("Configuration loaded successfully (env: %s, port: %s)\n", cfg.Server.Env, cfg.Server.Port)
	}

	return &cfg, nil
}

// GetWriteDSN returns the database DSN string for write operations
func (c *Config) GetWriteDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.PG.Write.Host, c.Database.PG.Write.Port, c.Database.PG.Write.User, c.Database.PG.Write.Password, c.Database.PG.Write.Name, c.Database.PG.Write.SSLMode)
}

// GetReadDSN returns the database DSN string for read operations
func (c *Config) GetReadDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.PG.Read.Host, c.Database.PG.Read.Port, c.Database.PG.Read.User, c.Database.PG.Read.Password, c.Database.PG.Read.Name, c.Database.PG.Read.SSLMode)
}

// GetDSN returns the database DSN string (defaults to write DSN for backward compatibility)
func (c *Config) GetDSN() string {
	return c.GetWriteDSN()
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	if c.Server.Host == "" {
		return fmt.Sprintf(":%s", c.Server.Port)
	}
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// GetHost returns the server host
func (c *Config) GetHost() string {
	if c.Server.Host == "" {
		return "0.0.0.0"
	}
	return c.Server.Host
}
