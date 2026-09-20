package models

type Version struct {
	Number  string `json:"number"`
	Service string `json:"service"`
}

type VersionResponse struct {
	Version Version `json:"version"`
}
