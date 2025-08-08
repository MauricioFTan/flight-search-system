package model

type SearchRequest struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Date       string `json:"date"`
	Passengers int    `json:"passengers"`
}

type SearchResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    SearchData `json:"data"`
}

type SearchData struct {
	SearchID string `json:"search_id"`
	Status   string `json:"status"`
}

type SearchRequestData struct {
	SearchID   string
	From       string
	To         string
	Date       string
	Passengers int
}
