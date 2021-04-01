package models

import (
	"encoding/json"
	"log"
	"os"
	"scratch-test/schema"
	"strings"
)

// Setup the schema, resolvers and the checkers for the clinics
func SetupClinic() (schema.IResolver, error) {
	s := `{
		"name": 	 "string",
		"from":  	 "string",
		"to": 	 	 "string",
		"state": 	 "string"
	}`

	sch, err := schema.ParseJSONToBaseResolver(strings.NewReader(s))
	if err != nil {
		return nil, err
	}

	sch.AddResolver("string", func(s string) (interface{}, error) {
		return s, nil
	})

	sch.AddEqualChecker("", func(v1, v2 interface{}) (bool, error) {
		// We are sure that all variables are string
		v1s, _ := v1.(string)
		v2s, _ := v2.(string)
		return strings.EqualFold(v1s, v2s), nil
	})

	// For states, both state name or state code could be used
	states := make(map[string]string)
	f, err := os.Open("states-ansi.json")
	if err != nil {
		return nil, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	stateCodes := make(map[string]string)
	if err := json.NewDecoder(f).Decode(&stateCodes); err != nil {
		return nil, err
	}
	// Lower all the entries in the tables
	for k, v := range stateCodes {
		states[strings.ToLower(k)] = strings.ToLower(v)
	}

	// States are equal if names and codes match
	// Alaska == alaska == ak == AK
	sch.AddEqualChecker("state", func(v1, v2 interface{}) (bool, error) {
		// We are sure that all variables are string
		v1s, _ := v1.(string)
		v2s, _ := v2.(string)

		k1 := strings.ToLower(v1s)
		k2 := strings.ToLower(v2s)

		if k1 == k2 {
			return true, nil
		}
		if code, ok := states[k1]; ok {
			if code == k2 {
				return true, nil
			}
		}
		if code, ok := states[k2]; ok {
			if code == k1 {
				return true, nil
			}
		}

		return false, nil
	})

	return sch, nil
}
