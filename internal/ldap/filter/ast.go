// Package filter implements an LDAP filter parser (RFC 4515) and an evaluator
// that converts filter ASTs into Ent ORM SQL predicates.
package filter

import "fmt"

// FilterType represents the type of an LDAP filter node.
type FilterType int

const (
	// FilterAnd represents an AND filter (&).
	FilterAnd FilterType = iota
	// FilterOr represents an OR filter (|).
	FilterOr
	// FilterNot represents a NOT filter (!).
	FilterNot
	// FilterEqual represents an equality match (=).
	FilterEqual
	// FilterSubstring represents a substring match (=*...*).
	FilterSubstring
	// FilterGreaterOrEqual represents a greater-or-equal match (>=).
	FilterGreaterOrEqual
	// FilterLessOrEqual represents a less-or-equal match (<=).
	FilterLessOrEqual
	// FilterPresent represents a presence test (=*).
	FilterPresent
	// FilterApproxMatch represents an approximate match (~=).
	FilterApproxMatch
)

// String returns a human-readable name for the filter type.
func (ft FilterType) String() string {
	switch ft {
	case FilterAnd:
		return "AND"
	case FilterOr:
		return "OR"
	case FilterNot:
		return "NOT"
	case FilterEqual:
		return "Equal"
	case FilterSubstring:
		return "Substring"
	case FilterGreaterOrEqual:
		return "GreaterOrEqual"
	case FilterLessOrEqual:
		return "LessOrEqual"
	case FilterPresent:
		return "Present"
	case FilterApproxMatch:
		return "ApproxMatch"
	default:
		return fmt.Sprintf("Unknown(%d)", int(ft))
	}
}

// Filter represents a parsed LDAP filter as an abstract syntax tree node.
type Filter struct {
	// Type is the kind of filter (AND, OR, NOT, Equal, etc.).
	Type FilterType
	// Attr is the LDAP attribute name for simple filters.
	Attr string
	// Value is the assertion value for equality, GTE, LTE, and approx filters.
	Value string
	// Children holds sub-filters for AND, OR, and NOT compound filters.
	Children []*Filter
	// Substr holds substring match components when Type is FilterSubstring.
	Substr *SubstringFilter
}

// SubstringFilter holds the components of a substring assertion.
type SubstringFilter struct {
	// Initial is the prefix before the first wildcard.
	Initial string
	// Any contains the middle parts between wildcards.
	Any []string
	// Final is the suffix after the last wildcard.
	Final string
}

// String returns a human-readable representation of the filter tree.
func (f *Filter) String() string {
	if f == nil {
		return "<nil>"
	}
	switch f.Type {
	case FilterAnd:
		s := "(&"
		for _, c := range f.Children {
			s += c.String()
		}
		return s + ")"
	case FilterOr:
		s := "(|"
		for _, c := range f.Children {
			s += c.String()
		}
		return s + ")"
	case FilterNot:
		if len(f.Children) > 0 {
			return "(!" + f.Children[0].String() + ")"
		}
		return "(!)"
	case FilterEqual:
		return fmt.Sprintf("(%s=%s)", f.Attr, f.Value)
	case FilterPresent:
		return fmt.Sprintf("(%s=*)", f.Attr)
	case FilterGreaterOrEqual:
		return fmt.Sprintf("(%s>=%s)", f.Attr, f.Value)
	case FilterLessOrEqual:
		return fmt.Sprintf("(%s<=%s)", f.Attr, f.Value)
	case FilterApproxMatch:
		return fmt.Sprintf("(%s~=%s)", f.Attr, f.Value)
	case FilterSubstring:
		s := fmt.Sprintf("(%s=", f.Attr)
		if f.Substr != nil {
			s += f.Substr.Initial + "*"
			for _, a := range f.Substr.Any {
				s += a + "*"
			}
			s += f.Substr.Final
		}
		return s + ")"
	default:
		return fmt.Sprintf("(?type=%d)", int(f.Type))
	}
}
