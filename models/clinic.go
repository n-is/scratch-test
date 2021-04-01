package models

import (
	"bytes"
	"encoding/json"
)

type available struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Clinic struct {
	Name         string     `json:"name"`
	State        string     `json:"stateName"`
	Availability *available `json:"availability"`

	// Alternate names
	ClinicName string     `json:"clinicName,omitempty"`
	StateCode  string     `json:"stateCode,omitempty"`
	Opening    *available `json:"opening,omitempty"`
}

// Corrects fields that have multiple possible keys
func (c *Clinic) Init() {
	if c.Name == "" && c.ClinicName != "" {
		c.Name = c.ClinicName
	}
	if c.State == "" && c.StateCode != "" {
		c.State = c.StateCode
	}
	if c.Availability == nil && c.Opening != nil {
		c.Availability = c.Opening
	}
}

func (c Clinic) Get(s string) interface{} {
	switch s {
	case "name":
		return c.Name
	case "state":
		return c.State
	case "from":
		return c.Availability.From
	case "to":
		return c.Availability.To
	}
	return nil
}

func (c Clinic) String() string {
	bts, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(bts)
}

// GetClinics reads the database output of the clinics and formats
// it in the structure of the clinics
func GetClinics(entries []map[string]interface{}) ([]Clinic, error) {

	bts, err := json.Marshal(entries)
	if err != nil {
		return nil, err
	}

	var values []struct {
		Name  string `json:"name"`
		State string `json:"state"`
		From  string `json:"from"`
		To    string `json:"to"`
	}

	err = json.NewDecoder(bytes.NewReader(bts)).Decode(&values)
	if err != nil {
		return nil, err
	}

	var clinics []Clinic
	for _, v := range values {
		clinics = append(clinics, Clinic{
			Name:  v.Name,
			State: v.State,
			Availability: &available{
				From: v.From,
				To:   v.To,
			},
		})
	}

	return clinics, nil
}
