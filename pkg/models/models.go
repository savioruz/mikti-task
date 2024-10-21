package models

type WebResponse[T any] struct {
	Data    T       `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	Error   string  `json:"error,omitempty"`
}

type GetsRequest struct {
	Page int `query:"page" validate:"numeric"`
	Size int `query:"size" validate:"numeric"`
}
