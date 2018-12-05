package httputil

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	qualityValueEntityRE = regexp.MustCompile(`^((?:\*|[\w+]+)/(?:\*|[\w+]+))(?:;q=([01](?:\.\d{1,3})?))?$`)
)

// See https://developer.mozilla.org/en-US/docs/Glossary/quality_values.
type QualityValue struct {
	MIMEType string
	Priority float64
}

// ParseQualityValues parses quality value text from HTTP header.
// This function does not sort the entities.
// Use SortQualityValues to sort quality values by priority.
// See https://developer.mozilla.org/en-US/docs/Glossary/quality_values.
func ParseQualityValues(s string) ([]QualityValue, error) {
	// example: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"

	var values []QualityValue

	// split into individual elements (separated by comma)
	for _, t := range tokenizeString(s, ",") {
		m := qualityValueEntityRE.FindAllStringSubmatch(strings.TrimSpace(t), -1)
		if len(m) != 1 || len(m[0]) != 3 {
			return nil, fmt.Errorf("invalid HTTP quality value: %q", s)
		}

		var priority float64 = 1 // default priority: 1.0
		if m[0][2] != "" {
			f, err := strconv.ParseFloat(m[0][2], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid HTTP quality value (containing invalid priority value: %s): %q", m[0][2], s)
			}

			priority = f
		}

		values = append(values, QualityValue{
			MIMEType: m[0][1],
			Priority: priority,
		})

	}

	return values, nil
}

// SortQualityValues sorts quality values based on their priority and MIME types.
// See https://developer.mozilla.org/en-US/docs/Glossary/quality_values.
func SortQualityValues(values []QualityValue) {
	sort.Slice(values, func(i, j int) bool {
		switch {
		case values[i].Priority > values[j].Priority:
			if values[i].Priority - values[j].Priority >= 0.001 {
				return true
			}
		case values[j].Priority > values[i].Priority:
			if values[j].Priority - values[i].Priority >= 0.001 {
				return false
			}
		}

		// equal priority: more specific MIME values have priority over less specific ones
		si := 2
		if values[i].MIMEType == "*/*" {
			si = 0
		} else if strings.HasSuffix(values[i].MIMEType, "/*") {
			si = 1
		}
		sj := 2
		if values[j].MIMEType == "*/*" {
			sj = 0
		} else if strings.HasSuffix(values[j].MIMEType, "/*") {
			sj = 1
		}
		return si > sj
	})
}

// tokenizeString slices s into all substrings separated by sep
// and returns a slice of the substrings between those separators.
// Empty elements are ignored.
func tokenizeString(s, sep string) []string {
	var tokens []string = nil

	for _, t := range strings.Split(s, sep) {
		if t != "" {
			tokens = append(tokens, t)
		}
	}

	return tokens
}
