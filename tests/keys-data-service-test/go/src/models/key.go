package models

import (
	"encoding/json"
	"time"
)

type Key struct {
	Reference string          `json:"reference"`
	Name      string          `json:"name"`
	Group     string          `json:"group,omitempty"`
	Value     json.RawMessage `json:"value,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type KeyResponse struct {
	Key Key `json:"key"`
}

type KeysResponse struct {
	Keys  []Key `json:"keys"`
	Skip  int   `json:"skip"`
	Limit int   `json:"limit"`
	Total int   `json:"total"`
	Pages int   `json:"pages"`
	Page  int   `json:"page"`
}
