package models

import "encoding/json"

type Message struct {
	Category string          `json:"category"`
	Data     json.RawMessage `json:"data"`
}
