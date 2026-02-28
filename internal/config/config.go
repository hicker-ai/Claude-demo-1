package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	LDAP     LDAPConfig     `mapstructure:"ldap"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	HTTPPort int `mapstructure:"http_port"`
}

// DatabaseConfig holds PostgreSQL connection configuration.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// DSN returns the PostgreSQL connection string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

// LDAPConfig holds LDAP server configuration.
type LDAPConfig struct {
	Port   int    `mapstructure:"port"`
	BaseDN string `mapstructure:"base_dn"`
	Mode   string `mapstructure:"mode"` // "openldap" or "activedirectory"
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// JWTConfig holds JWT authentication configuration.
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// Load reads configuration from the specified YAML file.
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}
