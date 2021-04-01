package utils

import (
	"net/url"
	"strconv"
	"strings"
)

// Remove cases where keys are not available
func FilterURLQuery(values url.Values) (map[string][]interface{}, int) {

	limit := -1
	if vs, ok := values["limit"]; ok {
		for _, v := range vs {
			if vi, err := strconv.Atoi(v); err == nil {
				if vi > limit {
					limit = vi
				}
			}
		}
	}

	if vs, ok := values["all"]; ok {
		for _, v := range vs {
			// In this scenario all the data are read without any filter
			// The limit still works
			if strings.EqualFold(v, "true") {
				return nil, limit
			}
		}
	}

	filtered := make(map[string][]interface{})
	for k, vs := range values {
		if k != "" {
			for _, v := range vs {
				filtered[k] = append(filtered[k], v)
			}
		}
	}

	return filtered, limit
}
