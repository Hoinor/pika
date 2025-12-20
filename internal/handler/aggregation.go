package handler

import "strings"

func normalizeAggregation(raw string) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	switch value {
	case "avg", "max":
		return value
	default:
		return ""
	}
}
