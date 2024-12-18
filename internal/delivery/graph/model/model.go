// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphmodel

import (
	"github.com/savioruz/mikti-task/internal/domain/model"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Mutation struct {
}

type PageMetadata struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type Query struct {
}

type TodoResponse struct {
	Data   []*model.TodoResponse `json:"data,omitempty"`
	Paging *PageMetadata         `json:"paging"`
	Error  *Error                `json:"error,omitempty"`
}
