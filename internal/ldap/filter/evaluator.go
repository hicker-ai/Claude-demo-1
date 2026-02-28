package filter

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
)

// AttrMapper maps LDAP attribute names to database column names.
type AttrMapper interface {
	// MapAttribute returns the database column name for the given LDAP attribute.
	// If the attribute is not recognized, ok is false.
	MapAttribute(ldapAttr string) (dbColumn string, ok bool)
}

// Evaluator converts a Filter AST into Ent ORM SQL predicates.
type Evaluator struct {
	mapper AttrMapper
}

// NewEvaluator creates a new Evaluator with the given attribute mapper.
func NewEvaluator(mapper AttrMapper) *Evaluator {
	return &Evaluator{mapper: mapper}
}

// Evaluate converts a Filter AST into an Ent SQL predicate.
// It returns an error if the filter references an unmapped attribute.
func (e *Evaluator) Evaluate(f *Filter) (*sql.Predicate, error) {
	if f == nil {
		return nil, fmt.Errorf("ldap evaluator: nil filter")
	}
	return e.eval(f)
}

// eval recursively converts a filter node into a SQL predicate.
func (e *Evaluator) eval(f *Filter) (*sql.Predicate, error) {
	switch f.Type {
	case FilterAnd:
		return e.evalAnd(f)
	case FilterOr:
		return e.evalOr(f)
	case FilterNot:
		return e.evalNot(f)
	case FilterEqual:
		return e.evalEqual(f)
	case FilterPresent:
		return e.evalPresent(f)
	case FilterSubstring:
		return e.evalSubstring(f)
	case FilterGreaterOrEqual:
		return e.evalGreaterOrEqual(f)
	case FilterLessOrEqual:
		return e.evalLessOrEqual(f)
	case FilterApproxMatch:
		return e.evalApproxMatch(f)
	default:
		return nil, fmt.Errorf("ldap evaluator: unsupported filter type: %s", f.Type)
	}
}

// evalAnd builds an AND predicate from child filters.
func (e *Evaluator) evalAnd(f *Filter) (*sql.Predicate, error) {
	preds, err := e.evalChildren(f.Children)
	if err != nil {
		return nil, err
	}
	if len(preds) == 0 {
		return nil, fmt.Errorf("ldap evaluator: AND filter has no evaluable children")
	}
	if len(preds) == 1 {
		return preds[0], nil
	}
	return sql.And(preds...), nil
}

// evalOr builds an OR predicate from child filters.
func (e *Evaluator) evalOr(f *Filter) (*sql.Predicate, error) {
	preds, err := e.evalChildren(f.Children)
	if err != nil {
		return nil, err
	}
	if len(preds) == 0 {
		return nil, fmt.Errorf("ldap evaluator: OR filter has no evaluable children")
	}
	if len(preds) == 1 {
		return preds[0], nil
	}
	return sql.Or(preds...), nil
}

// evalNot builds a NOT predicate from the single child filter.
func (e *Evaluator) evalNot(f *Filter) (*sql.Predicate, error) {
	if len(f.Children) == 0 {
		return nil, fmt.Errorf("ldap evaluator: NOT filter has no child")
	}
	child, err := e.eval(f.Children[0])
	if err != nil {
		return nil, err
	}
	return sql.Not(child), nil
}

// evalEqual builds an equality predicate.
func (e *Evaluator) evalEqual(f *Filter) (*sql.Predicate, error) {
	col, err := e.resolveAttr(f.Attr)
	if err != nil {
		return nil, err
	}
	return sql.EQ(col, f.Value), nil
}

// evalPresent builds a NOT NULL predicate.
func (e *Evaluator) evalPresent(f *Filter) (*sql.Predicate, error) {
	col, err := e.resolveAttr(f.Attr)
	if err != nil {
		return nil, err
	}
	return sql.Not(sql.IsNull(col)), nil
}

// evalSubstring builds a LIKE predicate from substring components.
func (e *Evaluator) evalSubstring(f *Filter) (*sql.Predicate, error) {
	col, err := e.resolveAttr(f.Attr)
	if err != nil {
		return nil, err
	}
	if f.Substr == nil {
		return nil, fmt.Errorf("ldap evaluator: substring filter has nil SubstringFilter")
	}

	pattern := buildLikePattern(f.Substr)
	return sql.Like(col, pattern), nil
}

// buildLikePattern constructs a SQL LIKE pattern from substring components.
// Wildcards (*) are converted to SQL % wildcards.
func buildLikePattern(s *SubstringFilter) string {
	var b strings.Builder

	if s.Initial == "" {
		b.WriteByte('%')
	} else {
		b.WriteString(escapeLikeValue(s.Initial))
		b.WriteByte('%')
	}

	for _, a := range s.Any {
		b.WriteString(escapeLikeValue(a))
		b.WriteByte('%')
	}

	if s.Final != "" {
		b.WriteString(escapeLikeValue(s.Final))
	}

	return b.String()
}

// escapeLikeValue escapes SQL LIKE special characters in a value string.
func escapeLikeValue(s string) string {
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

// evalGreaterOrEqual builds a GTE predicate.
func (e *Evaluator) evalGreaterOrEqual(f *Filter) (*sql.Predicate, error) {
	col, err := e.resolveAttr(f.Attr)
	if err != nil {
		return nil, err
	}
	return sql.GTE(col, f.Value), nil
}

// evalLessOrEqual builds a LTE predicate.
func (e *Evaluator) evalLessOrEqual(f *Filter) (*sql.Predicate, error) {
	col, err := e.resolveAttr(f.Attr)
	if err != nil {
		return nil, err
	}
	return sql.LTE(col, f.Value), nil
}

// evalApproxMatch builds a case-insensitive equality comparison.
func (e *Evaluator) evalApproxMatch(f *Filter) (*sql.Predicate, error) {
	col, err := e.resolveAttr(f.Attr)
	if err != nil {
		return nil, err
	}
	return sql.EqualFold(col, f.Value), nil
}

// resolveAttr maps an LDAP attribute to a database column.
// The objectClass attribute is skipped (returns a tautology predicate handled at
// handler level for routing). All other unmapped attributes produce an error.
func (e *Evaluator) resolveAttr(attr string) (string, error) {
	if strings.EqualFold(attr, "objectClass") {
		return "", fmt.Errorf("ldap evaluator: objectClass attribute should be handled at handler level")
	}
	col, ok := e.mapper.MapAttribute(attr)
	if !ok {
		return "", fmt.Errorf("ldap evaluator: unmapped attribute %q", attr)
	}
	return col, nil
}

// evalChildren evaluates a slice of child filters, skipping objectClass filters.
func (e *Evaluator) evalChildren(children []*Filter) ([]*sql.Predicate, error) {
	var preds []*sql.Predicate
	for _, child := range children {
		if isObjectClassFilter(child) {
			continue
		}
		p, err := e.eval(child)
		if err != nil {
			return nil, err
		}
		preds = append(preds, p)
	}
	return preds, nil
}

// isObjectClassFilter returns true if the filter is a simple filter on the
// objectClass attribute (equality or presence).
func isObjectClassFilter(f *Filter) bool {
	if f == nil {
		return false
	}
	switch f.Type {
	case FilterEqual, FilterPresent, FilterSubstring, FilterApproxMatch:
		return strings.EqualFold(f.Attr, "objectClass")
	default:
		return false
	}
}
