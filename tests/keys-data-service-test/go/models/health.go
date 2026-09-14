package models

type Health struct {
	Status string `json:"status"`
}

type HealthResponse struct {
	Health  Health `json:"health"`
	Service string `json:"service"`
	Version string `json:"version"`
}
