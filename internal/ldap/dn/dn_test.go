package dn

import (
	"testing"
)

const testBaseDN = "dc=example,dc=com"

func TestBuildUserDN(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		displayName string
		baseDN      string
		mode        string
		want        string
	}{
		{
			name:        "openldap user dn",
			username:    "jdoe",
			displayName: "John Doe",
			baseDN:      testBaseDN,
			mode:        ModeOpenLDAP,
			want:        "uid=jdoe,ou=users,dc=example,dc=com",
		},
		{
			name:        "active directory user dn",
			username:    "jdoe",
			displayName: "John Doe",
			baseDN:      testBaseDN,
			mode:        ModeActiveDirectory,
			want:        "cn=John Doe,cn=Users,dc=example,dc=com",
		},
		{
			name:        "openldap user dn with special chars",
			username:    "j,doe",
			displayName: "John, Doe",
			baseDN:      testBaseDN,
			mode:        ModeOpenLDAP,
			want:        "uid=j\\,doe,ou=users,dc=example,dc=com",
		},
		{
			name:        "active directory user dn with special chars",
			username:    "jdoe",
			displayName: "Doe, John",
			baseDN:      testBaseDN,
			mode:        ModeActiveDirectory,
			want:        "cn=Doe\\, John,cn=Users,dc=example,dc=com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildUserDN(tt.username, tt.displayName, tt.baseDN, tt.mode)
			if got != tt.want {
				t.Errorf("BuildUserDN() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildGroupDN(t *testing.T) {
	tests := []struct {
		name      string
		groupName string
		baseDN    string
		mode      string
		want      string
	}{
		{
			name:      "openldap group dn",
			groupName: "developers",
			baseDN:    testBaseDN,
			mode:      ModeOpenLDAP,
			want:      "cn=developers,ou=groups,dc=example,dc=com",
		},
		{
			name:      "active directory group dn",
			groupName: "developers",
			baseDN:    testBaseDN,
			mode:      ModeActiveDirectory,
			want:      "cn=developers,cn=Groups,dc=example,dc=com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildGroupDN(tt.groupName, tt.baseDN, tt.mode)
			if got != tt.want {
				t.Errorf("BuildGroupDN() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseDN(t *testing.T) {
	tests := []struct {
		name    string
		dn      string
		want    []RDN
		wantErr bool
	}{
		{
			name: "simple user dn",
			dn:   "uid=jdoe,ou=users,dc=example,dc=com",
			want: []RDN{
				{Type: "uid", Value: "jdoe"},
				{Type: "ou", Value: "users"},
				{Type: "dc", Value: "example"},
				{Type: "dc", Value: "com"},
			},
		},
		{
			name: "ad user dn",
			dn:   "cn=John Doe,cn=Users,dc=example,dc=com",
			want: []RDN{
				{Type: "cn", Value: "John Doe"},
				{Type: "cn", Value: "Users"},
				{Type: "dc", Value: "example"},
				{Type: "dc", Value: "com"},
			},
		},
		{
			name: "escaped comma in value",
			dn:   "cn=Doe\\, John,cn=Users,dc=example,dc=com",
			want: []RDN{
				{Type: "cn", Value: "Doe, John"},
				{Type: "cn", Value: "Users"},
				{Type: "dc", Value: "example"},
				{Type: "dc", Value: "com"},
			},
		},
		{
			name:    "empty dn",
			dn:      "",
			wantErr: true,
		},
		{
			name:    "malformed rdn",
			dn:      "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDN(tt.dn)
			if tt.wantErr {
				if err == nil {
					t.Fatal("ParseDN() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseDN() unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("ParseDN() returned %d RDNs, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i].Type != tt.want[i].Type || got[i].Value != tt.want[i].Value {
					t.Errorf("RDN[%d] = {%q, %q}, want {%q, %q}",
						i, got[i].Type, got[i].Value, tt.want[i].Type, tt.want[i].Value)
				}
			}
		})
	}
}

func TestExtractUsername(t *testing.T) {
	tests := []struct {
		name    string
		dn      string
		baseDN  string
		mode    string
		want    string
		wantErr bool
	}{
		{
			name:   "openldap extract uid",
			dn:     "uid=jdoe,ou=users,dc=example,dc=com",
			baseDN: testBaseDN,
			mode:   ModeOpenLDAP,
			want:   "jdoe",
		},
		{
			name:   "ad extract cn",
			dn:     "cn=John Doe,cn=Users,dc=example,dc=com",
			baseDN: testBaseDN,
			mode:   ModeActiveDirectory,
			want:   "John Doe",
		},
		{
			name:    "not a user dn",
			dn:      "cn=developers,ou=groups,dc=example,dc=com",
			baseDN:  testBaseDN,
			mode:    ModeOpenLDAP,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractUsername(tt.dn, tt.baseDN, tt.mode)
			if tt.wantErr {
				if err == nil {
					t.Fatal("ExtractUsername() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ExtractUsername() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ExtractUsername() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsUserDN(t *testing.T) {
	tests := []struct {
		name   string
		dn     string
		baseDN string
		mode   string
		want   bool
	}{
		{
			name:   "openldap user dn",
			dn:     "uid=jdoe,ou=users,dc=example,dc=com",
			baseDN: testBaseDN,
			mode:   ModeOpenLDAP,
			want:   true,
		},
		{
			name:   "ad user dn",
			dn:     "cn=John Doe,cn=Users,dc=example,dc=com",
			baseDN: testBaseDN,
			mode:   ModeActiveDirectory,
			want:   true,
		},
		{
			name:   "openldap group dn is not user dn",
			dn:     "cn=developers,ou=groups,dc=example,dc=com",
			baseDN: testBaseDN,
			mode:   ModeOpenLDAP,
			want:   false,
		},
		{
			name:   "ad group dn is not user dn",
			dn:     "cn=developers,cn=Groups,dc=example,dc=com",
			baseDN: testBaseDN,
			mode:   ModeActiveDirectory,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUserDN(tt.dn, tt.baseDN, tt.mode)
			if got != tt.want {
				t.Errorf("IsUserDN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsGroupDN(t *testing.T) {
	tests := []struct {
		name   string
		dn     string
		baseDN string
		mode   string
		want   bool
	}{
		{
			name:   "openldap group dn",
			dn:     "cn=developers,ou=groups,dc=example,dc=com",
			baseDN: testBaseDN,
			mode:   ModeOpenLDAP,
			want:   true,
		},
		{
			name:   "ad group dn",
			dn:     "cn=developers,cn=Groups,dc=example,dc=com",
			baseDN: testBaseDN,
			mode:   ModeActiveDirectory,
			want:   true,
		},
		{
			name:   "openldap user dn is not group dn",
			dn:     "uid=jdoe,ou=users,dc=example,dc=com",
			baseDN: testBaseDN,
			mode:   ModeOpenLDAP,
			want:   false,
		},
		{
			name:   "ad user dn is not group dn",
			dn:     "cn=John Doe,cn=Users,dc=example,dc=com",
			baseDN: testBaseDN,
			mode:   ModeActiveDirectory,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsGroupDN(tt.dn, tt.baseDN, tt.mode)
			if got != tt.want {
				t.Errorf("IsGroupDN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserBaseDN(t *testing.T) {
	tests := []struct {
		name   string
		baseDN string
		mode   string
		want   string
	}{
		{
			name:   "openldap",
			baseDN: testBaseDN,
			mode:   ModeOpenLDAP,
			want:   "ou=users,dc=example,dc=com",
		},
		{
			name:   "active directory",
			baseDN: testBaseDN,
			mode:   ModeActiveDirectory,
			want:   "cn=Users,dc=example,dc=com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserBaseDN(tt.baseDN, tt.mode)
			if got != tt.want {
				t.Errorf("UserBaseDN() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGroupBaseDN(t *testing.T) {
	tests := []struct {
		name   string
		baseDN string
		mode   string
		want   string
	}{
		{
			name:   "openldap",
			baseDN: testBaseDN,
			mode:   ModeOpenLDAP,
			want:   "ou=groups,dc=example,dc=com",
		},
		{
			name:   "active directory",
			baseDN: testBaseDN,
			mode:   ModeActiveDirectory,
			want:   "cn=Groups,dc=example,dc=com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupBaseDN(tt.baseDN, tt.mode)
			if got != tt.want {
				t.Errorf("GroupBaseDN() = %q, want %q", got, tt.want)
			}
		})
	}
}
