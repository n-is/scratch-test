package utils

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"reflect"
	"testing"
)

func TestFilterURLQuery(t *testing.T) {

	tests := []struct {
		url    string
		filter map[string][]interface{}
		limit  int
	}{
		{
			"http://example.com?state=fl&state=ak&&limit=2",
			map[string][]interface{}{
				"state": {"fl", "ak"},
				"limit": {"2"},
			},
			2,
		},
		{
			"http://example.com?state=fl&state=ak&&",
			map[string][]interface{}{
				"state": {"fl", "ak"},
			},
			-1,
		},
		{
			"http://example.com?state=fl&state=alaska",
			map[string][]interface{}{
				"state": {"fl", "alaska"},
			},
			-1,
		},
		{
			"http://example.com?state=fl&state=ak&&from=09:00",
			map[string][]interface{}{
				"state": {"fl", "ak"},
				"from":  {"09:00"},
			},
			-1,
		},
		{
			"http://example.com?state=fl&state=ak&&limit=2&to=20:00&to=19:00",
			map[string][]interface{}{
				"state": {"fl", "ak"},
				"to":    {"20:00", "19:00"},
				"limit": {"2"},
			},
			2,
		},
		{
			"http://example.com?name=Dental Clinic&limit=20",
			map[string][]interface{}{
				"name":  {"Dental Clinic"},
				"limit": {"20"},
			},
			20,
		},
		{
			"http://example.com?name=Dental Clinic&limit=-20",
			map[string][]interface{}{
				"name":  {"Dental Clinic"},
				"limit": {"-20"},
			},
			-1,
		},
	}

	for _, tst := range tests {
		uri, err := url.ParseRequestURI(tst.url)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		filter, limit := FilterURLQuery(uri.Query())
		assert.Equal(t, tst.limit, limit)

		if !reflect.DeepEqual(filter, tst.filter) {
			t.Errorf("Error: Want: %v, Got: %v", tst.filter, filter)
		}
	}
}
