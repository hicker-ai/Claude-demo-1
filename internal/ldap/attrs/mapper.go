// Package attrs provides LDAP attribute mapping for OpenLDAP and
// Active Directory modes.
package attrs

const (
	// ModeOpenLDAP indicates OpenLDAP attribute mapping.
	ModeOpenLDAP = "openldap"
	// ModeActiveDirectory indicates Active Directory attribute mapping.
	ModeActiveDirectory = "activedirectory"
)

// Mapper translates between database column names and LDAP attribute
// names for a given LDAP mode.
type Mapper struct {
	mode string
}

// NewMapper creates a new attribute Mapper for the specified mode.
func NewMapper(mode string) *Mapper {
	return &Mapper{mode: mode}
}

// MapAttribute maps an LDAP attribute name to the corresponding
// database column name. It returns the column name and true if a
// mapping exists, or an empty string and false otherwise.
func (m *Mapper) MapAttribute(ldapAttr string) (dbColumn string, ok bool) {
	table := openLDAPAttrMap
	if m.mode == ModeActiveDirectory {
		table = adAttrMap
	}
	col, found := table[ldapAttr]
	return col, found
}

// UserObjectClasses returns the objectClass values for user entries
// in the current LDAP mode.
func (m *Mapper) UserObjectClasses() []string {
	if m.mode == ModeActiveDirectory {
		return []string{"top", "person", "organizationalPerson", "user"}
	}
	return []string{"top", "person", "organizationalPerson", "inetOrgPerson"}
}

// GroupObjectClasses returns the objectClass values for group entries
// in the current LDAP mode.
func (m *Mapper) GroupObjectClasses() []string {
	if m.mode == ModeActiveDirectory {
		return []string{"top", "group"}
	}
	return []string{"top", "groupOfNames"}
}

// UserToLDAPAttrs converts a user's fields to LDAP attributes for the
// current mode.
func (m *Mapper) UserToLDAPAttrs(username, displayName, email, phone, status string) map[string][]string {
	attrs := map[string][]string{
		"objectClass": m.UserObjectClasses(),
		"cn":          {displayName},
		"displayName": {displayName},
	}

	if m.mode == ModeActiveDirectory {
		attrs["sAMAccountName"] = []string{username}
		attrs["userAccountControl"] = []string{adAccountControl(status)}
	} else {
		attrs["uid"] = []string{username}
		attrs["sn"] = []string{username} // inetOrgPerson requires sn
		attrs["status"] = []string{status}
	}

	if email != "" {
		attrs["mail"] = []string{email}
	}
	if phone != "" {
		attrs["telephoneNumber"] = []string{phone}
	}

	return attrs
}

// GroupToLDAPAttrs converts a group's fields to LDAP attributes for
// the current mode.
func (m *Mapper) GroupToLDAPAttrs(name, description string, memberDNs []string) map[string][]string {
	attrs := map[string][]string{
		"objectClass": m.GroupObjectClasses(),
		"cn":          {name},
	}

	if description != "" {
		attrs["description"] = []string{description}
	}
	if len(memberDNs) > 0 {
		attrs["member"] = memberDNs
	}

	return attrs
}

// --- internal helpers ---

// openLDAPAttrMap maps OpenLDAP attribute names to DB column names.
var openLDAPAttrMap = map[string]string{
	"uid":             "username",
	"cn":              "display_name",
	"displayName":     "display_name",
	"mail":            "email",
	"telephoneNumber": "phone",
	"status":          "status",
}

// adAttrMap maps Active Directory attribute names to DB column names.
var adAttrMap = map[string]string{
	"sAMAccountName":     "username",
	"cn":                 "display_name",
	"displayName":        "display_name",
	"mail":               "email",
	"telephoneNumber":    "phone",
	"userAccountControl": "status",
}

// adAccountControl converts a simple status string to an AD
// userAccountControl value. "active" yields 512 (normal account);
// anything else yields 514 (disabled).
func adAccountControl(status string) string {
	if status == "active" {
		return "512"
	}
	return "514"
}
