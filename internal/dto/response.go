package dto

type ResponseNewsList struct {
	Success bool   `json:"Success"`
	News    []News `json:"News"`
}
