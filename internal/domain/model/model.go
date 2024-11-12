package model

type Response[T any] struct {
	Data    T             `json:"data,omitempty"`
	Message *Message      `json:"message,omitempty"`
	Paging  *PageMetadata `json:"paging,omitempty"`
}

type Message struct {
	Message string `json:"message"`
}

type PageMetadata struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
