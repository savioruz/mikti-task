package resolvers

import (
	"github.com/savioruz/mikti-task/tree/week-4/internal/usecases/todo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TodoUsecase todo.TodoUsecase
}

func NewResolver(todo todo.TodoUsecase) *Resolver {
	return &Resolver{
		TodoUsecase: todo,
	}
}
