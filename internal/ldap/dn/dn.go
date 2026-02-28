// Package dn provides building and parsing of LDAP Distinguished Names
// for both OpenLDAP and Active Directory modes.
package dn

import (
	"fmt"
	"strings"
)

const (
	// ModeOpenLDAP indicates OpenLDAP-style DN construction.
	ModeOpenLDAP = "openldap"
	// ModeActiveDirectory indicates Active Directory-style DN construction.
	ModeActiveDirectory = "activedirectory"
)

// RDN represents a Relative Distinguished Name component.
type RDN struct {
	Type  string
	Value string
}

// String returns the RDN in "type=value" form.
func (r RDN) String() string {
	return r.Type + "=" + escapeRDNValue(r.Value)
}

// BuildUserDN builds a user DN based on LDAP mode.
// OpenLDAP: uid=<username>,ou=users,<baseDN>
// AD: cn=<displayName>,cn=Users,<baseDN>
func BuildUserDN(username, displayName, baseDN, mode string) string {
	if mode == ModeActiveDirectory {
		return "cn=" + escapeRDNValue(displayName) + "," + userContainer(mode) + "," + baseDN
	}
	return "uid=" + escapeRDNValue(username) + "," + userContainer(mode) + "," + baseDN
}

// BuildGroupDN builds a group DN based on LDAP mode.
// OpenLDAP: cn=<groupName>,ou=groups,<baseDN>
// AD: cn=<groupName>,cn=Groups,<baseDN>
func BuildGroupDN(groupName, baseDN, mode string) string {
	return "cn=" + escapeRDNValue(groupName) + "," + groupContainer(mode) + "," + baseDN
}

// ParseDN parses a DN string into RDN components.
func ParseDN(dn string) ([]RDN, error) {
	if dn == "" {
		return nil, fmt.Errorf("dn: %w", errEmptyDN)
	}

	parts := splitDN(dn)
	rdns := make([]RDN, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		idx := strings.IndexByte(part, '=')
		if idx < 0 {
			return nil, fmt.Errorf("dn: invalid RDN component %q: %w", part, errMalformedRDN)
		}

		rdns = append(rdns, RDN{
			Type:  strings.TrimSpace(part[:idx]),
			Value: unescapeRDNValue(strings.TrimSpace(part[idx+1:])),
		})
	}

	if len(rdns) == 0 {
		return nil, fmt.Errorf("dn: %w", errEmptyDN)
	}

	return rdns, nil
}

// ExtractUsername extracts the username from a user DN.
// OpenLDAP: extracts uid= value.
// AD: extracts cn= value.
func ExtractUsername(dn, baseDN, mode string) (string, error) {
	if !IsUserDN(dn, baseDN, mode) {
		return "", fmt.Errorf("dn: %q is not a valid user DN for mode %q: %w", dn, mode, errNotUserDN)
	}

	rdns, err := ParseDN(dn)
	if err != nil {
		return "", fmt.Errorf("dn: extracting username: %w", err)
	}

	if len(rdns) == 0 {
		return "", fmt.Errorf("dn: no RDN components found: %w", errNotUserDN)
	}

	first := rdns[0]
	if mode == ModeActiveDirectory {
		if !strings.EqualFold(first.Type, "cn") {
			return "", fmt.Errorf("dn: expected cn= in AD user DN, got %q: %w", first.Type, errNotUserDN)
		}
		return first.Value, nil
	}

	if !strings.EqualFold(first.Type, "uid") {
		return "", fmt.Errorf("dn: expected uid= in OpenLDAP user DN, got %q: %w", first.Type, errNotUserDN)
	}
	return first.Value, nil
}

// IsUserDN checks if a DN is a user DN.
func IsUserDN(dn, baseDN, mode string) bool {
	suffix := "," + UserBaseDN(baseDN, mode)
	return strings.HasSuffix(strings.ToLower(dn), strings.ToLower(suffix))
}

// IsGroupDN checks if a DN is a group DN.
func IsGroupDN(dn, baseDN, mode string) bool {
	suffix := "," + GroupBaseDN(baseDN, mode)
	return strings.HasSuffix(strings.ToLower(dn), strings.ToLower(suffix))
}

// UserBaseDN returns the base DN for users.
// OpenLDAP: ou=users,<baseDN>
// AD: cn=Users,<baseDN>
func UserBaseDN(baseDN, mode string) string {
	return userContainer(mode) + "," + baseDN
}

// GroupBaseDN returns the base DN for groups.
// OpenLDAP: ou=groups,<baseDN>
// AD: cn=Groups,<baseDN>
func GroupBaseDN(baseDN, mode string) string {
	return groupContainer(mode) + "," + baseDN
}

// --- internal helpers ---

var (
	errEmptyDN      = fmt.Errorf("empty DN")
	errMalformedRDN = fmt.Errorf("malformed RDN")
	errNotUserDN    = fmt.Errorf("not a user DN")
)

func userContainer(mode string) string {
	if mode == ModeActiveDirectory {
		return "cn=Users"
	}
	return "ou=users"
}

func groupContainer(mode string) string {
	if mode == ModeActiveDirectory {
		return "cn=Groups"
	}
	return "ou=groups"
}

// splitDN splits a DN string on commas while respecting escaped commas.
func splitDN(dn string) []string {
	var parts []string
	var buf strings.Builder
	escaped := false

	for i := 0; i < len(dn); i++ {
		ch := dn[i]
		if escaped {
			buf.WriteByte(ch)
			escaped = false
			continue
		}
		if ch == '\\' {
			buf.WriteByte(ch)
			escaped = true
			continue
		}
		if ch == ',' {
			parts = append(parts, buf.String())
			buf.Reset()
			continue
		}
		buf.WriteByte(ch)
	}

	if buf.Len() > 0 {
		parts = append(parts, buf.String())
	}

	return parts
}

// escapeRDNValue escapes special characters in an RDN value per RFC 4514.
func escapeRDNValue(val string) string {
	var buf strings.Builder
	for i, ch := range val {
		switch {
		case ch == ',' || ch == '+' || ch == '"' || ch == '\\' || ch == '<' || ch == '>' || ch == ';':
			buf.WriteByte('\\')
			buf.WriteRune(ch)
		case ch == '#' && i == 0:
			buf.WriteByte('\\')
			buf.WriteRune(ch)
		case ch == ' ' && (i == 0 || i == len(val)-1):
			buf.WriteByte('\\')
			buf.WriteRune(ch)
		default:
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

// unescapeRDNValue removes backslash escapes from an RDN value.
func unescapeRDNValue(val string) string {
	var buf strings.Builder
	escaped := false
	for _, ch := range val {
		if escaped {
			buf.WriteRune(ch)
			escaped = false
			continue
		}
		if ch == '\\' {
			escaped = true
			continue
		}
		buf.WriteRune(ch)
	}
	return buf.String()
}
