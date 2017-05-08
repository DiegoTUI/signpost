package models

import (
	"encoding/json"
)

// Capital defines a capital model
type Capital struct {
	Country string `json:"country"`
	Capital string `json:"capital"`
	Type    string `json:"type"`
}

// NewCapital creates a new capital from a JSON
func NewCapital(jsonByte []byte) (*Capital, error) {
	var capital Capital

	if err := json.Unmarshal(jsonByte, &capital); err != nil {
		return nil, err
	}

	return &capital, nil
}
