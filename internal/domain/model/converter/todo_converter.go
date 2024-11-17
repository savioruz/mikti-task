package converter

import (
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/domain/model"
)

func TodoToResponse(todo *entity.Todo) *model.TodoResponse {
	return &model.TodoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Done:      todo.Done,
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

func TodosToPaginatedResponse(todos []entity.Todo, totalItems int64, page, size int) *model.Response[[]*model.TodoResponse] {
	todoResponses := TodosToResponses(todos)
	totalPages := (int(totalItems) + size - 1) / size

	return model.NewResponse(todoResponses, &model.PageMetadata{
		Page:       page,
		Size:       size,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	})
}
