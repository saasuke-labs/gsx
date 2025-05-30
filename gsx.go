package gsx

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

var componentTag = regexp.MustCompile(`<([A-Z][A-Za-z0-9]*)\s*([^>]*)\s*/?>`)
var attributePattern = regexp.MustCompile(`(\w+)=("[^"]*"|\{[^}]*\}|[^\s"{}=<>` + "`]+)")

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

// LintWarning represents a non-fatal parsing issue
type LintWarning struct {
	Message string
	Line    int
	Column  int
}

// Options can be extended later to control parsing behavior

func ParseString(name, input string) (*template.Template, []LintWarning, error) {
	return ParseStringInto(name, input, nil)
}

func ParseStringInto(name, input string, parent *template.Template) (*template.Template, []LintWarning, error) {
	funcs := template.FuncMap{
		"props": func(attrs ...any) map[string]any {
			result := map[string]any{}

			for i := 0; i < len(attrs); i += 2 {
				if i+1 >= len(attrs) {
					continue // Skip if there's no value for the key
				}
				key := attrs[i]
				value := attrs[i+1]

				result[fmt.Sprintf("%v", key)] = value
			}

			return result
		}}

	if parent == nil {
		parent = template.New(name)
	}
	var warnings []LintWarning

	matches := componentTag.FindAllStringSubmatchIndex(input, -1)
	var buf strings.Builder
	lastIndex := 0

	for _, match := range matches {
		fullStart, fullEnd := match[0], match[1]
		nameStart, nameEnd := match[2], match[3]
		attrStart, attrEnd := match[4], match[5]

		tagName := input[nameStart:nameEnd]
		attrStr := input[attrStart:attrEnd]

		// Write everything before this match
		buf.WriteString(input[lastIndex:fullStart])

		attrs := parseAttributes(attrStr, &warnings)
		parts := []string{}
		for k, v := range attrs {
			parts = append(parts, fmt.Sprintf("\"%s\" %s", k, v))
		}

		if len(parts) == 0 {
			buf.WriteString(fmt.Sprintf(`{{ template "%s" }}`, tagName))
		} else {
			buf.WriteString(fmt.Sprintf(`{{ template "%s" (props %s) }}`, tagName, strings.Join(parts, " ")))
		}

		lastIndex = fullEnd
	}

	// Write the remainder of the input
	buf.WriteString(input[lastIndex:])

	output := buf.String()

	// Wrap the final output in a template definition so it can be invoked
	// in other templates
	output = fmt.Sprintf(`{{define "%s"}}%s{{end}}`, name, output)
	parent, err := parent.New(name).Funcs(funcs).Parse(output)

	if err != nil {
		return nil, warnings, fmt.Errorf("failed to parse GSX: %w", err)
	}

	return parent, warnings, nil
}

func parseAttributes(attrStr string, warnings *[]LintWarning) map[string]string {
	result := map[string]string{}
	matches := attributePattern.FindAllStringSubmatch(attrStr, -1)
	for _, m := range matches {
		key := m[1]
		val := strings.TrimSpace(m[2])

		// Detect format
		if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
			result[key] = val
		} else if strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}") {
			inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(val, "{"), "}"))

			if isNumeric(inner) {
				result[key] = inner
			} else {
				result[key] = fmt.Sprintf(".%s", inner)

			}
		} else {
			*warnings = append(*warnings, LintWarning{
				Message: fmt.Sprintf("unquoted value for attribute \"%s\"", key),
			})
			result[key] = fmt.Sprintf("\"%s\"", val)
		}
	}
	return result
}

// RenderTemplate executes the parsed GSX content as a Go template
func RenderTemplate(parsed string, data any, funcs template.FuncMap) (string, error) {
	tmpl, err := template.New("gsx").Funcs(funcs).Parse(parsed)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
