package models

import (
	"encoding/json"
	"time"
)

type Key struct {
	Reference string `json:"reference"`

	Name  string          `json:"name"`
	Group string          `json:"group,omitempty"`
	Value json.RawMessage `json:"value,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
