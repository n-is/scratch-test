package schema

import (
	"encoding/json"
)

type Token struct {
	Id    string      `json:"id"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// For printing purposes
func (t Token) String() string {
	bts, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(bts)
}
