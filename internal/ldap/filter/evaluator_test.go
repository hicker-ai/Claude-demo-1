package filter

import (
	"strings"
	"testing"

	"entgo.io/ent/dialect/sql"
)

// mockMapper implements AttrMapper for testing.
type mockMapper struct {
	mapping map[string]string
}

func (m *mockMapper) MapAttribute(ldapAttr string) (string, bool) {
	col, ok := m.mapping[strings.ToLower(ldapAttr)]
	return col, ok
}

func newMockMapper() *mockMapper {
	return &mockMapper{
		mapping: map[string]string{
			"cn":     "name",
			"mail":   "email",
			"age":    "age",
			"status": "status",
			"sn":     "surname",
		},
	}
}

// predicateToSQL renders a *sql.Predicate as a SQL WHERE clause for testing.
func predicateToSQL(p *sql.Predicate) string {
	selector := sql.Select("*").From(sql.Table("users"))
	selector.Where(p)
	query, args := selector.Query()
	// Replace placeholders with args for readable output.
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			query = strings.Replace(query, "?", "'"+v+"'", 1)
		default:
			query = strings.Replace(query, "?", "?", 1)
		}
	}
	return query
}

func TestEvaluateEqual(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{Type: FilterEqual, Attr: "cn", Value: "John"}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	if !strings.Contains(query, "`name`") || !strings.Contains(query, "'John'") {
		t.Errorf("Equal predicate SQL = %q, want to contain `name` and 'John'", query)
	}
}

func TestEvaluatePresent(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{Type: FilterPresent, Attr: "cn"}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	if !strings.Contains(query, "NOT") || !strings.Contains(query, "NULL") {
		t.Errorf("Present predicate SQL = %q, want NOT ... NULL", query)
	}
}

func TestEvaluateSubstringPrefix(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{
		Type: FilterSubstring,
		Attr: "cn",
		Substr: &SubstringFilter{
			Initial: "Jo",
		},
	}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	if !strings.Contains(query, "LIKE") || !strings.Contains(query, "Jo%") {
		t.Errorf("Substring prefix SQL = %q, want LIKE 'Jo%%'", query)
	}
}

func TestEvaluateSubstringSuffix(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{
		Type: FilterSubstring,
		Attr: "cn",
		Substr: &SubstringFilter{
			Final: "ohn",
		},
	}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	if !strings.Contains(query, "LIKE") || !strings.Contains(query, "%ohn") {
		t.Errorf("Substring suffix SQL = %q, want LIKE '%%ohn'", query)
	}
}

func TestEvaluateSubstringComplex(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{
		Type: FilterSubstring,
		Attr: "mail",
		Substr: &SubstringFilter{
			Initial: "user",
			Any:     []string{"mid"},
			Final:   "example.com",
		},
	}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	if !strings.Contains(query, "LIKE") {
		t.Errorf("Substring complex SQL = %q, want LIKE pattern", query)
	}
	if !strings.Contains(query, "user%mid%example.com") {
		t.Errorf("Substring complex SQL = %q, want pattern user%%mid%%example.com", query)
	}
}

func TestEvaluateGreaterOrEqual(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{Type: FilterGreaterOrEqual, Attr: "age", Value: "18"}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	if !strings.Contains(query, ">=") || !strings.Contains(query, "'18'") {
		t.Errorf("GTE predicate SQL = %q, want >= '18'", query)
	}
}

func TestEvaluateLessOrEqual(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{Type: FilterLessOrEqual, Attr: "age", Value: "65"}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	if !strings.Contains(query, "<=") || !strings.Contains(query, "'65'") {
		t.Errorf("LTE predicate SQL = %q, want <= '65'", query)
	}
}

func TestEvaluateApproxMatch(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{Type: FilterApproxMatch, Attr: "cn", Value: "Jon"}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	// EqualFold uses COLLATE or LOWER comparison.
	if !strings.Contains(query, "`name`") {
		t.Errorf("ApproxMatch predicate SQL = %q, want to reference `name`", query)
	}
}

func TestEvaluateAnd(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{
		Type: FilterAnd,
		Children: []*Filter{
			{Type: FilterEqual, Attr: "cn", Value: "John"},
			{Type: FilterEqual, Attr: "mail", Value: "j@e.com"},
		},
	}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	if !strings.Contains(query, "`name`") || !strings.Contains(query, "`email`") {
		t.Errorf("AND predicate SQL = %q, want both `name` and `email`", query)
	}
}

func TestEvaluateOr(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{
		Type: FilterOr,
		Children: []*Filter{
			{Type: FilterEqual, Attr: "cn", Value: "John"},
			{Type: FilterEqual, Attr: "cn", Value: "Jane"},
		},
	}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	if !strings.Contains(query, "OR") {
		t.Errorf("OR predicate SQL = %q, want OR clause", query)
	}
}

func TestEvaluateNot(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{
		Type: FilterNot,
		Children: []*Filter{
			{Type: FilterEqual, Attr: "cn", Value: "John"},
		},
	}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	if !strings.Contains(query, "NOT") {
		t.Errorf("NOT predicate SQL = %q, want NOT clause", query)
	}
}

func TestEvaluateUnmappedAttribute(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	f := &Filter{Type: FilterEqual, Attr: "unknownAttr", Value: "test"}

	_, err := e.Evaluate(f)
	if err == nil {
		t.Fatal("Evaluate() expected error for unmapped attribute, got nil")
	}
	if !strings.Contains(err.Error(), "unmapped attribute") {
		t.Errorf("error = %q, want to contain 'unmapped attribute'", err.Error())
	}
}

func TestEvaluateObjectClassSkipped(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	// AND with objectClass and a real filter: objectClass should be skipped.
	f := &Filter{
		Type: FilterAnd,
		Children: []*Filter{
			{Type: FilterEqual, Attr: "objectClass", Value: "inetOrgPerson"},
			{Type: FilterEqual, Attr: "cn", Value: "John"},
		},
	}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	query := predicateToSQL(p)
	if strings.Contains(query, "objectClass") || strings.Contains(query, "objectclass") {
		t.Errorf("objectClass should be skipped, SQL = %q", query)
	}
	if !strings.Contains(query, "`name`") {
		t.Errorf("SQL = %q, want to contain `name` for cn filter", query)
	}
}

func TestEvaluateNilFilter(t *testing.T) {
	e := NewEvaluator(newMockMapper())
	_, err := e.Evaluate(nil)
	if err == nil {
		t.Fatal("Evaluate(nil) expected error, got nil")
	}
}

func TestEvaluateIntegration(t *testing.T) {
	// Parse a filter string and evaluate it end-to-end.
	e := NewEvaluator(newMockMapper())
	f, err := Parse("(&(cn=John)(mail=*@example.com))")
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	p, err := e.Evaluate(f)
	if err != nil {
		t.Fatalf("Evaluate() error: %v", err)
	}

	// Build the full query to verify it is valid SQL.
	query := predicateToSQL(p)
	if query == "" {
		t.Error("generated SQL query is empty")
	}
	t.Logf("Generated SQL: %s", query)
}

func TestBuildLikePattern(t *testing.T) {
	tests := []struct {
		name   string
		substr *SubstringFilter
		want   string
	}{
		{
			name:   "prefix only",
			substr: &SubstringFilter{Initial: "Jo"},
			want:   "Jo%",
		},
		{
			name:   "suffix only",
			substr: &SubstringFilter{Final: "ohn"},
			want:   "%ohn",
		},
		{
			name:   "any only",
			substr: &SubstringFilter{Any: []string{"oh"}},
			want:   "%oh%",
		},
		{
			name:   "complex",
			substr: &SubstringFilter{Initial: "J", Any: []string{"o"}, Final: "hn"},
			want:   "J%o%hn",
		},
		{
			name:   "escape percent in value",
			substr: &SubstringFilter{Initial: "100%"},
			want:   "100\\%%",
		},
		{
			name:   "escape underscore in value",
			substr: &SubstringFilter{Initial: "a_b"},
			want:   "a\\_b%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildLikePattern(tt.substr)
			if got != tt.want {
				t.Errorf("buildLikePattern() = %q, want %q", got, tt.want)
			}
		})
	}
}
