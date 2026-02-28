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
  level: "debug"
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
		{"DBHost", cfg.Database.Host, "testhost"},
		{"DBPort", cfg.Database.Port, 5433},
		{"DBUser", cfg.Database.User, "testuser"},
		{"DBPassword", cfg.Database.Password, "testpass"},
		{"DBName", cfg.Database.DBName, "testdb"},
		{"DBSSLMode", cfg.Database.SSLMode, "require"},
		{"LDAPPort", cfg.LDAP.Port, 10389},
		{"LDAPBaseDN", cfg.LDAP.BaseDN, "dc=test,dc=com"},
		{"LDAPMode", cfg.LDAP.Mode, "activedirectory"},
		{"LogLevel", cfg.Log.Level, "debug"},
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

func TestLoadDSN(t *testing.T) {
	content := `
database:
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

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("Load() expected error for missing file, got nil")
	}
}
