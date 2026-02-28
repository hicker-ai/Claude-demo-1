package filter

import (
	"fmt"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
)

// BER tag constants mirroring github.com/go-ldap/ldap/v3 filter types.
const (
	tagFilterAnd             = 0
	tagFilterOr              = 1
	tagFilterNot             = 2
	tagFilterEqualityMatch   = 3
	tagFilterSubstrings      = 4
	tagFilterGreaterOrEqual  = 5
	tagFilterLessOrEqual     = 6
	tagFilterPresent         = 7
	tagFilterApproxMatch     = 8
	tagFilterExtensibleMatch = 9

	tagSubstringsInitial = 0
	tagSubstringsAny     = 1
	tagSubstringsFinal   = 2
)

// Parse parses an RFC 4515 LDAP filter string into a Filter AST.
// It returns a descriptive error for malformed filters.
func Parse(filterStr string) (*Filter, error) {
	if filterStr == "" {
		return nil, fmt.Errorf("ldap filter: empty filter string")
	}

	packet, err := ldap.CompileFilter(filterStr)
	if err != nil {
		return nil, fmt.Errorf("ldap filter: failed to compile filter %q: %w", filterStr, err)
	}

	f, err := packetToFilter(packet)
	if err != nil {
		return nil, fmt.Errorf("ldap filter: failed to convert BER packet: %w", err)
	}

	return f, nil
}

// packetToFilter recursively converts a BER packet tree into a Filter AST.
func packetToFilter(p *ber.Packet) (*Filter, error) {
	if p == nil {
		return nil, fmt.Errorf("nil BER packet")
	}

	tag := int(p.Tag)
	switch tag {
	case tagFilterAnd:
		return parseCompoundFilter(FilterAnd, p)
	case tagFilterOr:
		return parseCompoundFilter(FilterOr, p)
	case tagFilterNot:
		return parseNotFilter(p)
	case tagFilterEqualityMatch:
		return parseComparisonFilter(FilterEqual, p)
	case tagFilterSubstrings:
		return parseSubstringFilter(p)
	case tagFilterGreaterOrEqual:
		return parseComparisonFilter(FilterGreaterOrEqual, p)
	case tagFilterLessOrEqual:
		return parseComparisonFilter(FilterLessOrEqual, p)
	case tagFilterPresent:
		return parsePresentFilter(p)
	case tagFilterApproxMatch:
		return parseComparisonFilter(FilterApproxMatch, p)
	case tagFilterExtensibleMatch:
		return nil, fmt.Errorf("extensible match filters are not supported")
	default:
		return nil, fmt.Errorf("unknown BER filter tag: %d", tag)
	}
}

// parseCompoundFilter handles AND and OR filters with multiple children.
func parseCompoundFilter(filterType FilterType, p *ber.Packet) (*Filter, error) {
	f := &Filter{Type: filterType}
	for _, child := range p.Children {
		cf, err := packetToFilter(child)
		if err != nil {
			return nil, fmt.Errorf("in %s filter: %w", filterType, err)
		}
		f.Children = append(f.Children, cf)
	}
	if len(f.Children) == 0 {
		return nil, fmt.Errorf("%s filter must have at least one child", filterType)
	}
	return f, nil
}

// parseNotFilter handles NOT filters with exactly one child.
func parseNotFilter(p *ber.Packet) (*Filter, error) {
	if len(p.Children) != 1 {
		return nil, fmt.Errorf("NOT filter must have exactly one child, got %d", len(p.Children))
	}
	child, err := packetToFilter(p.Children[0])
	if err != nil {
		return nil, fmt.Errorf("in NOT filter: %w", err)
	}
	return &Filter{
		Type:     FilterNot,
		Children: []*Filter{child},
	}, nil
}

// parseComparisonFilter handles equality, GTE, LTE, and approx match filters.
func parseComparisonFilter(filterType FilterType, p *ber.Packet) (*Filter, error) {
	if len(p.Children) != 2 {
		return nil, fmt.Errorf("%s filter must have 2 children (attr, value), got %d", filterType, len(p.Children))
	}
	attr := packetStringValue(p.Children[0])
	value := packetStringValue(p.Children[1])
	return &Filter{
		Type:  filterType,
		Attr:  attr,
		Value: value,
	}, nil
}

// parsePresentFilter handles presence filters (attr=*).
func parsePresentFilter(p *ber.Packet) (*Filter, error) {
	attr := packetStringValue(p)
	if attr == "" {
		return nil, fmt.Errorf("present filter has empty attribute name")
	}
	return &Filter{
		Type: FilterPresent,
		Attr: attr,
	}, nil
}

// parseSubstringFilter handles substring match filters (attr=init*any*final).
func parseSubstringFilter(p *ber.Packet) (*Filter, error) {
	if len(p.Children) < 2 {
		return nil, fmt.Errorf("substring filter must have at least 2 children (attr, sequence), got %d", len(p.Children))
	}

	attr := packetStringValue(p.Children[0])
	seq := p.Children[1]

	substr := &SubstringFilter{}
	for _, child := range seq.Children {
		val := packetStringValue(child)
		switch int(child.Tag) {
		case tagSubstringsInitial:
			substr.Initial = val
		case tagSubstringsAny:
			substr.Any = append(substr.Any, val)
		case tagSubstringsFinal:
			substr.Final = val
		default:
			return nil, fmt.Errorf("unknown substring filter tag: %d", child.Tag)
		}
	}

	return &Filter{
		Type:   FilterSubstring,
		Attr:   attr,
		Substr: substr,
	}, nil
}

// packetStringValue extracts a string value from a BER packet.
func packetStringValue(p *ber.Packet) string {
	if p == nil {
		return ""
	}
	if p.Value != nil {
		if s, ok := p.Value.(string); ok {
			return s
		}
	}
	// Fall back to reading from the Data buffer.
	if p.Data != nil && p.Data.Len() > 0 {
		return p.Data.String()
	}
	// Fall back to ByteValue.
	if len(p.ByteValue) > 0 {
		return string(p.ByteValue)
	}
	return ""
}
