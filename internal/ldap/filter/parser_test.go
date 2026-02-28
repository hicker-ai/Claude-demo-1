package filter

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantType  FilterType
		wantAttr  string
		wantValue string
		wantErr   bool
		check     func(t *testing.T, f *Filter)
	}{
		{
			name:      "Equal",
			input:     "(cn=John)",
			wantType:  FilterEqual,
			wantAttr:  "cn",
			wantValue: "John",
		},
		{
			name:     "Presence",
			input:    "(cn=*)",
			wantType: FilterPresent,
			wantAttr: "cn",
		},
		{
			name:     "Substring prefix",
			input:    "(cn=Jo*)",
			wantType: FilterSubstring,
			wantAttr: "cn",
			check: func(t *testing.T, f *Filter) {
				t.Helper()
				if f.Substr == nil {
					t.Fatal("Substr is nil")
				}
				if f.Substr.Initial != "Jo" {
					t.Errorf("Initial = %q, want %q", f.Substr.Initial, "Jo")
				}
				if f.Substr.Final != "" {
					t.Errorf("Final = %q, want empty", f.Substr.Final)
				}
				if len(f.Substr.Any) != 0 {
					t.Errorf("Any = %v, want empty", f.Substr.Any)
				}
			},
		},
		{
			name:     "Substring suffix",
			input:    "(cn=*ohn)",
			wantType: FilterSubstring,
			wantAttr: "cn",
			check: func(t *testing.T, f *Filter) {
				t.Helper()
				if f.Substr == nil {
					t.Fatal("Substr is nil")
				}
				if f.Substr.Initial != "" {
					t.Errorf("Initial = %q, want empty", f.Substr.Initial)
				}
				if f.Substr.Final != "ohn" {
					t.Errorf("Final = %q, want %q", f.Substr.Final, "ohn")
				}
			},
		},
		{
			name:     "Substring any",
			input:    "(cn=*oh*)",
			wantType: FilterSubstring,
			wantAttr: "cn",
			check: func(t *testing.T, f *Filter) {
				t.Helper()
				if f.Substr == nil {
					t.Fatal("Substr is nil")
				}
				if f.Substr.Initial != "" {
					t.Errorf("Initial = %q, want empty", f.Substr.Initial)
				}
				if f.Substr.Final != "" {
					t.Errorf("Final = %q, want empty", f.Substr.Final)
				}
				if len(f.Substr.Any) != 1 || f.Substr.Any[0] != "oh" {
					t.Errorf("Any = %v, want [oh]", f.Substr.Any)
				}
			},
		},
		{
			name:     "Substring complex",
			input:    "(cn=J*o*hn)",
			wantType: FilterSubstring,
			wantAttr: "cn",
			check: func(t *testing.T, f *Filter) {
				t.Helper()
				if f.Substr == nil {
					t.Fatal("Substr is nil")
				}
				if f.Substr.Initial != "J" {
					t.Errorf("Initial = %q, want %q", f.Substr.Initial, "J")
				}
				if len(f.Substr.Any) != 1 || f.Substr.Any[0] != "o" {
					t.Errorf("Any = %v, want [o]", f.Substr.Any)
				}
				if f.Substr.Final != "hn" {
					t.Errorf("Final = %q, want %q", f.Substr.Final, "hn")
				}
			},
		},
		{
			name:      "GreaterOrEqual",
			input:     "(age>=18)",
			wantType:  FilterGreaterOrEqual,
			wantAttr:  "age",
			wantValue: "18",
		},
		{
			name:      "LessOrEqual",
			input:     "(age<=65)",
			wantType:  FilterLessOrEqual,
			wantAttr:  "age",
			wantValue: "65",
		},
		{
			name:      "ApproxMatch",
			input:     "(cn~=Jon)",
			wantType:  FilterApproxMatch,
			wantAttr:  "cn",
			wantValue: "Jon",
		},
		{
			name:     "AND",
			input:    "(&(cn=John)(mail=j@e.com))",
			wantType: FilterAnd,
			check: func(t *testing.T, f *Filter) {
				t.Helper()
				if len(f.Children) != 2 {
					t.Fatalf("Children count = %d, want 2", len(f.Children))
				}
				c0 := f.Children[0]
				if c0.Type != FilterEqual || c0.Attr != "cn" || c0.Value != "John" {
					t.Errorf("child[0] = %v, want Equal cn=John", c0)
				}
				c1 := f.Children[1]
				if c1.Type != FilterEqual || c1.Attr != "mail" || c1.Value != "j@e.com" {
					t.Errorf("child[1] = %v, want Equal mail=j@e.com", c1)
				}
			},
		},
		{
			name:     "OR",
			input:    "(|(cn=John)(cn=Jane))",
			wantType: FilterOr,
			check: func(t *testing.T, f *Filter) {
				t.Helper()
				if len(f.Children) != 2 {
					t.Fatalf("Children count = %d, want 2", len(f.Children))
				}
				if f.Children[0].Value != "John" {
					t.Errorf("child[0].Value = %q, want %q", f.Children[0].Value, "John")
				}
				if f.Children[1].Value != "Jane" {
					t.Errorf("child[1].Value = %q, want %q", f.Children[1].Value, "Jane")
				}
			},
		},
		{
			name:     "NOT",
			input:    "(!(cn=John))",
			wantType: FilterNot,
			check: func(t *testing.T, f *Filter) {
				t.Helper()
				if len(f.Children) != 1 {
					t.Fatalf("Children count = %d, want 1", len(f.Children))
				}
				c := f.Children[0]
				if c.Type != FilterEqual || c.Attr != "cn" || c.Value != "John" {
					t.Errorf("child = %v, want Equal cn=John", c)
				}
			},
		},
		{
			name:     "Nested AND with OR and Substring",
			input:    "(&(|(cn=John)(cn=Jane))(mail=*@example.com))",
			wantType: FilterAnd,
			check: func(t *testing.T, f *Filter) {
				t.Helper()
				if len(f.Children) != 2 {
					t.Fatalf("Children count = %d, want 2", len(f.Children))
				}
				or := f.Children[0]
				if or.Type != FilterOr {
					t.Errorf("child[0].Type = %v, want OR", or.Type)
				}
				if len(or.Children) != 2 {
					t.Fatalf("OR children count = %d, want 2", len(or.Children))
				}
				sub := f.Children[1]
				if sub.Type != FilterSubstring {
					t.Errorf("child[1].Type = %v, want Substring", sub.Type)
				}
				if sub.Attr != "mail" {
					t.Errorf("child[1].Attr = %q, want %q", sub.Attr, "mail")
				}
				if sub.Substr == nil || sub.Substr.Final != "@example.com" {
					t.Errorf("child[1].Substr.Final = %q, want %q", sub.Substr.Final, "@example.com")
				}
			},
		},
		{
			name:      "objectClass equality",
			input:     "(objectClass=inetOrgPerson)",
			wantType:  FilterEqual,
			wantAttr:  "objectClass",
			wantValue: "inetOrgPerson",
		},
		{
			name:     "Deep nested",
			input:    "(&(|(cn=A)(cn=B))(!(status=disabled))(mail=*@example.com))",
			wantType: FilterAnd,
			check: func(t *testing.T, f *Filter) {
				t.Helper()
				if len(f.Children) != 3 {
					t.Fatalf("Children count = %d, want 3", len(f.Children))
				}
				// First child: OR
				if f.Children[0].Type != FilterOr {
					t.Errorf("child[0].Type = %v, want OR", f.Children[0].Type)
				}
				if len(f.Children[0].Children) != 2 {
					t.Errorf("OR children count = %d, want 2", len(f.Children[0].Children))
				}
				// Second child: NOT
				if f.Children[1].Type != FilterNot {
					t.Errorf("child[1].Type = %v, want NOT", f.Children[1].Type)
				}
				notChild := f.Children[1].Children[0]
				if notChild.Attr != "status" || notChild.Value != "disabled" {
					t.Errorf("NOT child = %v, want Equal status=disabled", notChild)
				}
				// Third child: Substring
				if f.Children[2].Type != FilterSubstring {
					t.Errorf("child[2].Type = %v, want Substring", f.Children[2].Type)
				}
			},
		},
		{
			name:    "Empty filter",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Unbalanced parens",
			input:   "((cn=John)",
			wantErr: true,
		},
		{
			name:    "Missing opening paren",
			input:   "cn=John)",
			wantErr: true,
		},
		{
			name:    "Invalid filter syntax",
			input:   "(cn)",
			wantErr: true,
		},
		{
			name:      "LDAP escape star in value",
			input:     "(cn=John\\2aDoe)",
			wantType:  FilterEqual,
			wantAttr:  "cn",
			wantValue: "John*Doe",
		},
		{
			name:      "LDAP escape parens in value",
			input:     "(cn=John\\28Sr\\29)",
			wantType:  FilterEqual,
			wantAttr:  "cn",
			wantValue: "John(Sr)",
		},
		{
			name:      "LDAP escape backslash in value",
			input:     "(cn=John\\5cDoe)",
			wantType:  FilterEqual,
			wantAttr:  "cn",
			wantValue: "John\\Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Parse(%q) = %v, want error", tt.input, f)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.input, err)
			}
			if f.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", f.Type, tt.wantType)
			}
			if tt.wantAttr != "" && f.Attr != tt.wantAttr {
				t.Errorf("Attr = %q, want %q", f.Attr, tt.wantAttr)
			}
			if tt.wantValue != "" && f.Value != tt.wantValue {
				t.Errorf("Value = %q, want %q", f.Value, tt.wantValue)
			}
			if tt.check != nil {
				tt.check(t, f)
			}
		})
	}
}

func TestFilterString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Equal",
			input: "(cn=John)",
			want:  "(cn=John)",
		},
		{
			name:  "Presence",
			input: "(cn=*)",
			want:  "(cn=*)",
		},
		{
			name:  "AND",
			input: "(&(cn=John)(mail=test))",
			want:  "(&(cn=John)(mail=test))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.input, err)
			}
			got := f.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
