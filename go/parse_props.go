package vbml

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	sectionStartPattern = regexp.MustCompile(`\{\{([#^])(\w+)\}\}`)
	variablePattern     = regexp.MustCompile(`\{\{(\w+|\.)\}\}`)
)

func parseProps(props map[string]any, template string) string {
	if props == nil {
		props = map[string]any{}
	}
	return renderTemplate(template, props)
}

func renderTemplate(template string, context map[string]any) string {
	for {
		match := sectionStartPattern.FindStringSubmatchIndex(template)
		if match == nil {
			break
		}

		sigil := template[match[2]:match[3]]
		key := template[match[4]:match[5]]
		closeTag := "{{/" + key + "}}"
		closeIndex := strings.Index(template[match[1]:], closeTag)
		if closeIndex < 0 {
			break
		}

		bodyStart := match[1]
		bodyEnd := match[1] + closeIndex
		body := template[bodyStart:bodyEnd]
		value := context[key]

		replacement := renderSection(sigil, body, value, context)
		template = template[:match[0]] + replacement + template[bodyEnd+len(closeTag):]
	}

	return variablePattern.ReplaceAllStringFunc(template, func(token string) string {
		parts := variablePattern.FindStringSubmatch(token)
		if len(parts) != 2 {
			return token
		}

		key := parts[1]
		value, ok := resolveTemplateValue(context, key)
		if !ok {
			return ""
		}
		return stringifyTemplateValue(value)
	})
}

func renderSection(sigil, body string, value any, context map[string]any) string {
	switch sigil {
	case "#":
		iterable := reflect.ValueOf(value)
		if iterable.IsValid() && (iterable.Kind() == reflect.Slice || iterable.Kind() == reflect.Array) {
			var builder strings.Builder
			for index := 0; index < iterable.Len(); index++ {
				builder.WriteString(renderTemplate(body, mergeTemplateContext(context, iterable.Index(index).Interface())))
			}
			return builder.String()
		}

		if isTruthy(value) {
			return renderTemplate(body, context)
		}
		return ""
	case "^":
		if isTruthy(value) {
			return ""
		}
		return renderTemplate(body, context)
	default:
		return ""
	}
}

func mergeTemplateContext(parent map[string]any, current any) map[string]any {
	merged := make(map[string]any, len(parent)+1)
	for key, value := range parent {
		merged[key] = value
	}
	merged["."] = current
	return merged
}

func resolveTemplateValue(context map[string]any, key string) (any, bool) {
	value, ok := context[key]
	return value, ok
}

func stringifyTemplateValue(value any) string {
	switch value := value.(type) {
	case string:
		return value
	case nil:
		return ""
	default:
		return fmt.Sprint(value)
	}
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}

	switch value := value.(type) {
	case bool:
		return value
	case string:
		return value != ""
	}

	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return rv.Len() > 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() != 0
	case reflect.Interface, reflect.Pointer:
		return !rv.IsNil()
	default:
		return true
	}
}
