package attrs

import (
	"testing"
)

func TestMapAttribute(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		ldapAttr string
		wantCol  string
		wantOK   bool
	}{
		// OpenLDAP mappings
		{name: "openldap uid", mode: ModeOpenLDAP, ldapAttr: "uid", wantCol: "username", wantOK: true},
		{name: "openldap cn", mode: ModeOpenLDAP, ldapAttr: "cn", wantCol: "display_name", wantOK: true},
		{name: "openldap displayName", mode: ModeOpenLDAP, ldapAttr: "displayName", wantCol: "display_name", wantOK: true},
		{name: "openldap mail", mode: ModeOpenLDAP, ldapAttr: "mail", wantCol: "email", wantOK: true},
		{name: "openldap telephoneNumber", mode: ModeOpenLDAP, ldapAttr: "telephoneNumber", wantCol: "phone", wantOK: true},
		{name: "openldap status", mode: ModeOpenLDAP, ldapAttr: "status", wantCol: "status", wantOK: true},
		{name: "openldap unknown", mode: ModeOpenLDAP, ldapAttr: "foobar", wantCol: "", wantOK: false},

		// AD mappings
		{name: "ad sAMAccountName", mode: ModeActiveDirectory, ldapAttr: "sAMAccountName", wantCol: "username", wantOK: true},
		{name: "ad cn", mode: ModeActiveDirectory, ldapAttr: "cn", wantCol: "display_name", wantOK: true},
		{name: "ad displayName", mode: ModeActiveDirectory, ldapAttr: "displayName", wantCol: "display_name", wantOK: true},
		{name: "ad mail", mode: ModeActiveDirectory, ldapAttr: "mail", wantCol: "email", wantOK: true},
		{name: "ad telephoneNumber", mode: ModeActiveDirectory, ldapAttr: "telephoneNumber", wantCol: "phone", wantOK: true},
		{name: "ad userAccountControl", mode: ModeActiveDirectory, ldapAttr: "userAccountControl", wantCol: "status", wantOK: true},
		{name: "ad unknown", mode: ModeActiveDirectory, ldapAttr: "foobar", wantCol: "", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMapper(tt.mode)
			col, ok := m.MapAttribute(tt.ldapAttr)
			if ok != tt.wantOK {
				t.Errorf("MapAttribute(%q) ok = %v, want %v", tt.ldapAttr, ok, tt.wantOK)
			}
			if col != tt.wantCol {
				t.Errorf("MapAttribute(%q) = %q, want %q", tt.ldapAttr, col, tt.wantCol)
			}
		})
	}
}

func TestUserObjectClasses(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want []string
	}{
		{
			name: "openldap",
			mode: ModeOpenLDAP,
			want: []string{"top", "person", "organizationalPerson", "inetOrgPerson"},
		},
		{
			name: "active directory",
			mode: ModeActiveDirectory,
			want: []string{"top", "person", "organizationalPerson", "user"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMapper(tt.mode)
			got := m.UserObjectClasses()
			assertStringSliceEqual(t, got, tt.want)
		})
	}
}

func TestGroupObjectClasses(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want []string
	}{
		{
			name: "openldap",
			mode: ModeOpenLDAP,
			want: []string{"top", "groupOfNames"},
		},
		{
			name: "active directory",
			mode: ModeActiveDirectory,
			want: []string{"top", "group"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMapper(tt.mode)
			got := m.GroupObjectClasses()
			assertStringSliceEqual(t, got, tt.want)
		})
	}
}

func TestUserToLDAPAttrs(t *testing.T) {
	tests := []struct {
		name        string
		mode        string
		username    string
		displayName string
		email       string
		phone       string
		status      string
		wantKeys    []string
	}{
		{
			name:        "openldap full user",
			mode:        ModeOpenLDAP,
			username:    "jdoe",
			displayName: "John Doe",
			email:       "jdoe@example.com",
			phone:       "+1-555-0100",
			status:      "active",
			wantKeys:    []string{"objectClass", "cn", "displayName", "uid", "sn", "status", "mail", "telephoneNumber"},
		},
		{
			name:        "ad full user",
			mode:        ModeActiveDirectory,
			username:    "jdoe",
			displayName: "John Doe",
			email:       "jdoe@example.com",
			phone:       "+1-555-0100",
			status:      "active",
			wantKeys:    []string{"objectClass", "cn", "displayName", "sAMAccountName", "userAccountControl", "mail", "telephoneNumber"},
		},
		{
			name:        "openldap user without optional fields",
			mode:        ModeOpenLDAP,
			username:    "jdoe",
			displayName: "John Doe",
			email:       "",
			phone:       "",
			status:      "active",
			wantKeys:    []string{"objectClass", "cn", "displayName", "uid", "sn", "status"},
		},
		{
			name:        "ad disabled user",
			mode:        ModeActiveDirectory,
			username:    "jdoe",
			displayName: "John Doe",
			email:       "",
			phone:       "",
			status:      "disabled",
			wantKeys:    []string{"objectClass", "cn", "displayName", "sAMAccountName", "userAccountControl"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMapper(tt.mode)
			got := m.UserToLDAPAttrs(tt.username, tt.displayName, tt.email, tt.phone, tt.status)

			for _, key := range tt.wantKeys {
				if _, ok := got[key]; !ok {
					t.Errorf("UserToLDAPAttrs() missing key %q", key)
				}
			}

			// Verify specific values
			if got["cn"][0] != tt.displayName {
				t.Errorf("cn = %q, want %q", got["cn"][0], tt.displayName)
			}
			if got["displayName"][0] != tt.displayName {
				t.Errorf("displayName = %q, want %q", got["displayName"][0], tt.displayName)
			}

			if tt.mode == ModeOpenLDAP {
				if got["uid"][0] != tt.username {
					t.Errorf("uid = %q, want %q", got["uid"][0], tt.username)
				}
				if got["status"][0] != tt.status {
					t.Errorf("status = %q, want %q", got["status"][0], tt.status)
				}
			} else {
				if got["sAMAccountName"][0] != tt.username {
					t.Errorf("sAMAccountName = %q, want %q", got["sAMAccountName"][0], tt.username)
				}
				if tt.status == "active" && got["userAccountControl"][0] != "512" {
					t.Errorf("userAccountControl = %q, want %q", got["userAccountControl"][0], "512")
				}
				if tt.status == "disabled" && got["userAccountControl"][0] != "514" {
					t.Errorf("userAccountControl = %q, want %q", got["userAccountControl"][0], "514")
				}
			}

			// Verify optional fields not present when empty
			if tt.email == "" {
				if _, ok := got["mail"]; ok {
					t.Error("mail should not be present when email is empty")
				}
			}
			if tt.phone == "" {
				if _, ok := got["telephoneNumber"]; ok {
					t.Error("telephoneNumber should not be present when phone is empty")
				}
			}
		})
	}
}

func TestGroupToLDAPAttrs(t *testing.T) {
	tests := []struct {
		name        string
		mode        string
		groupName   string
		description string
		memberDNs   []string
		wantKeys    []string
	}{
		{
			name:        "openldap full group",
			mode:        ModeOpenLDAP,
			groupName:   "developers",
			description: "Development team",
			memberDNs:   []string{"uid=jdoe,ou=users,dc=example,dc=com"},
			wantKeys:    []string{"objectClass", "cn", "description", "member"},
		},
		{
			name:        "ad full group",
			mode:        ModeActiveDirectory,
			groupName:   "developers",
			description: "Development team",
			memberDNs:   []string{"cn=John Doe,cn=Users,dc=example,dc=com"},
			wantKeys:    []string{"objectClass", "cn", "description", "member"},
		},
		{
			name:        "openldap group without optional fields",
			mode:        ModeOpenLDAP,
			groupName:   "empty-group",
			description: "",
			memberDNs:   nil,
			wantKeys:    []string{"objectClass", "cn"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMapper(tt.mode)
			got := m.GroupToLDAPAttrs(tt.groupName, tt.description, tt.memberDNs)

			for _, key := range tt.wantKeys {
				if _, ok := got[key]; !ok {
					t.Errorf("GroupToLDAPAttrs() missing key %q", key)
				}
			}

			if got["cn"][0] != tt.groupName {
				t.Errorf("cn = %q, want %q", got["cn"][0], tt.groupName)
			}

			if tt.description == "" {
				if _, ok := got["description"]; ok {
					t.Error("description should not be present when empty")
				}
			} else if got["description"][0] != tt.description {
				t.Errorf("description = %q, want %q", got["description"][0], tt.description)
			}

			if len(tt.memberDNs) == 0 {
				if _, ok := got["member"]; ok {
					t.Error("member should not be present when memberDNs is empty")
				}
			} else {
				assertStringSliceEqual(t, got["member"], tt.memberDNs)
			}
		})
	}
}

// assertStringSliceEqual is a test helper that compares two string slices.
func assertStringSliceEqual(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("slice length = %d, want %d\n  got:  %v\n  want: %v", len(got), len(want), got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("slice[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
