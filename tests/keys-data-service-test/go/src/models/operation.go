package models

type Operation struct {
	Status    string `json:"status"`
	Procedure string `json:"procedure"`
}

type OperationResponse struct {
	Operation Operation `json:"operation"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}
