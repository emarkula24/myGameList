package utils

import (
	"strings"
)

func ParseSearchQuery(query string) string {
	querysplice := strings.ReplaceAll(query, " ", "-")
	return querysplice
}
