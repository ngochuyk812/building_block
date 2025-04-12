package dtos

type PagingModel[T any] struct {
	PageSize int `json:"page_size"`
	Page     int `json:"page"`
	Total    int `json:"total"`
	Items    []T `json:"items"`
}

type PagingRequest struct {
	PageSize int `json:"page_size" validate:"required"`
	Page     int `json:"page" validate:"required"`
}
