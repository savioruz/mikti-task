package converter

import (
	"github.com/savioruz/mikti-task/tree/week-3/internal/entities"
	"github.com/savioruz/mikti-task/tree/week-3/internal/models"
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
