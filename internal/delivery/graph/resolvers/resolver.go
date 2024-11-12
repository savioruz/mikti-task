package resolvers

import "github.com/savioruz/mikti-task/tree/week-4/internal/usecases"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TodoUsecase *usecases.TodoUsecase
}

func NewResolver(todo *usecases.TodoUsecase) *Resolver {
	return &Resolver{
		TodoUsecase: todo,
	}
}
