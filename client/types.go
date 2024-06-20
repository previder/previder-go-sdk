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
