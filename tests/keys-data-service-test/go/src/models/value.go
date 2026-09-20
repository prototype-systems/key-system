package models

import "encoding/json"

type ValueResponse struct {
	Value json.RawMessage `json:"value"`
}
