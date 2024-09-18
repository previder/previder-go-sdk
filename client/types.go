package client

import "encoding/json"

type PageRequest struct {
	Page  int
	Size  int
	Sort  string
	Query string
}

type Page struct {
	TotalPages       int
	TotalElements    int
	NumberOfElements int
	Size             int
	Number           int
	Content          json.RawMessage
}

type OwnerReference struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Audit struct {
	CreatedBy      AuditUserRef `json:"createdBy"`
	CreatedAt      int          `json:"createdAt"`
	LastModifiedBy AuditUserRef `json:"lastModifiedBy"`
	LastModifiedAt int          `json:"lastModifiedAt"`
}

type AuditUserRef struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
