package converter

import (
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/domain/model"
)

func TodoToResponse(todo *entity.Todo) *model.TodoResponse {
	return &model.TodoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt.String(),
		UpdatedAt: todo.UpdatedAt.String(),
	}
}

func TodosToResponses(todos []entity.Todo) []*model.TodoResponse {
	todoResponses := make([]*model.TodoResponse, len(todos))
	for i := range todos {
		todoResponses[i] = TodoToResponse(&todos[i])
	}
	return todoResponses
}
