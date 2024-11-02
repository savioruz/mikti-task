package converter

import (
	"github.com/savioruz/mikti-task/tree/week-4/internal/entities"
	"github.com/savioruz/mikti-task/tree/week-4/internal/models"
)

func TodoToResponse(todo *entities.Todo) *models.TodoResponse {
	return &models.TodoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt.String(),
		UpdatedAt: todo.UpdatedAt.String(),
	}
}

func TodosToResponses(todos []entities.Todo) []*models.TodoResponse {
	todoResponses := make([]*models.TodoResponse, len(todos))
	for i := range todos {
		todoResponses[i] = TodoToResponse(&todos[i])
	}
	return todoResponses
}
