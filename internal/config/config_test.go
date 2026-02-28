package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	content := `
server:
  http_port: 9090
database:
  driver: postgres
  host: testhost
  port: 5433
  user: testuser
  password: testpass
  dbname: testdb
  sslmode: require
ldap:
  port: 10389
  base_dn: "dc=test,dc=com"
  mode: "activedirectory"
log:
  console: true
  add_source: true
  level: "debug"
  format: "json"
  rotate:
    enabled: true
    output_path: "/tmp/logs"
    max_size: 50
    max_age: 14
    max_backups: 5
    compress: true
    local_time: true
jwt:
  secret: "testsecret"
  expire_hours: 12
`
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"HTTPPort", cfg.Server.HTTPPort, 9090},
		{"DBDriver", cfg.Database.Driver, "postgres"},
		{"DBHost", cfg.Database.Host, "testhost"},
		{"DBPort", cfg.Database.Port, 5433},
		{"DBUser", cfg.Database.User, "testuser"},
		{"DBPassword", cfg.Database.Password, "testpass"},
		{"DBName", cfg.Database.DBName, "testdb"},
		{"DBSSLMode", cfg.Database.SSLMode, "require"},
		{"LDAPPort", cfg.LDAP.Port, 10389},
		{"LDAPBaseDN", cfg.LDAP.BaseDN, "dc=test,dc=com"},
		{"LDAPMode", cfg.LDAP.Mode, "activedirectory"},
		{"LogConsole", cfg.Log.Console, true},
		{"LogAddSource", cfg.Log.AddSource, true},
		{"LogLevel", cfg.Log.Level, "debug"},
		{"LogFormat", cfg.Log.Format, "json"},
		{"LogRotateEnabled", cfg.Log.Rotate.Enabled, true},
		{"LogRotateOutputPath", cfg.Log.Rotate.OutputPath, "/tmp/logs"},
		{"LogRotateMaxSize", cfg.Log.Rotate.MaxSize, 50},
		{"LogRotateMaxAge", cfg.Log.Rotate.MaxAge, 14},
		{"LogRotateMaxBackups", cfg.Log.Rotate.MaxBackups, 5},
		{"LogRotateCompress", cfg.Log.Rotate.Compress, true},
		{"LogRotateLocalTime", cfg.Log.Rotate.LocalTime, true},
		{"JWTSecret", cfg.JWT.Secret, "testsecret"},
		{"JWTExpireHours", cfg.JWT.ExpireHours, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestLoadDSN_Postgres(t *testing.T) {
	content := `
database:
  driver: postgres
  host: localhost
  port: 5432
  user: postgres
  password: secret
  dbname: mydb
  sslmode: disable
`
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	want := "host=localhost port=5432 user=postgres password=secret dbname=mydb sslmode=disable"
	got := cfg.Database.DSN()
	if got != want {
		t.Errorf("DSN() = %q, want %q", got, want)
	}
}

func TestLoadDSN_SQLite(t *testing.T) {
	content := `
database:
  driver: sqlite3
  path: data/app.db
`
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	want := "file:data/app.db?_fk=1"
	got := cfg.Database.DSN()
	if got != want {
		t.Errorf("DSN() = %q, want %q", got, want)
	}
}

func TestEnsureDataDir(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "subdir", "test.db")

	cfg := DatabaseConfig{
		Driver: "sqlite3",
		Path:   dbPath,
	}

	if err := cfg.EnsureDataDir(); err != nil {
		t.Fatalf("EnsureDataDir() error: %v", err)
	}

	parentDir := filepath.Dir(dbPath)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		t.Errorf("EnsureDataDir() did not create directory %s", parentDir)
	}
}

func TestEnsureDataDir_Postgres(t *testing.T) {
	cfg := DatabaseConfig{
		Driver: "postgres",
		Host:   "localhost",
	}

	if err := cfg.EnsureDataDir(); err != nil {
		t.Fatalf("EnsureDataDir() should be no-op for postgres, got error: %v", err)
	}
}

func TestToLogsConfig(t *testing.T) {
	cfg := LogConfig{
		Console:   true,
		AddSource: true,
		Level:     "debug",
		Format:    "json",
		Rotate: RotateConfig{
			Enabled:    true,
			OutputPath: "/tmp/logs",
			MaxSize:    50,
			MaxAge:     14,
			MaxBackups: 5,
			Compress:   true,
			LocalTime:  true,
		},
	}

	lc := cfg.ToLogsConfig()

	if lc.Console != true {
		t.Errorf("Console = %v, want true", lc.Console)
	}
	if lc.AddSource != true {
		t.Errorf("AddSource = %v, want true", lc.AddSource)
	}
	if lc.Level != "debug" {
		t.Errorf("Level = %q, want %q", lc.Level, "debug")
	}
	if lc.Format != "json" {
		t.Errorf("Format = %q, want %q", lc.Format, "json")
	}
	if lc.Rotate.Enabled != true {
		t.Errorf("Rotate.Enabled = %v, want true", lc.Rotate.Enabled)
	}
	if lc.Rotate.OutputPath != "/tmp/logs" {
		t.Errorf("Rotate.OutputPath = %q, want %q", lc.Rotate.OutputPath, "/tmp/logs")
	}
	if lc.Rotate.MaxSize != 50 {
		t.Errorf("Rotate.MaxSize = %d, want %d", lc.Rotate.MaxSize, 50)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("Load() expected error for missing file, got nil")
	}
}
