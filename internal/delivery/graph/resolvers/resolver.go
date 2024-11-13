package resolvers

import (
	"github.com/savioruz/mikti-task/internal/usecases/todo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TodoUsecase todo.TodoUsecase
}

func NewResolver(t todo.TodoUsecase) *Resolver {
	return &Resolver{
		TodoUsecase: t,
	}
}
