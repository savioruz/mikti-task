package converter

import (
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/domain/model"
)

func TodoToResponse(todo *entity.Todo, isAdmin bool) *model.TodoResponse {
	response := &model.TodoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Done:      todo.Done,
		CreatedAt: todo.CreatedAt.String(),
		UpdatedAt: todo.UpdatedAt.String(),
	}

	if isAdmin {
		response.UserID = &todo.UserID
	}

	return response
}

func TodosToResponses(todos []entity.Todo, isAdmin bool) []*model.TodoResponse {
	todoResponses := make([]*model.TodoResponse, len(todos))
	for i := range todos {
		todoResponses[i] = TodoToResponse(&todos[i], isAdmin)
	}
	return todoResponses
}

func TodosToPaginatedResponse(todos []entity.Todo, totalItems int64, page, size int, isAdmin bool) *model.Response[[]*model.TodoResponse] {
	todoResponses := TodosToResponses(todos, isAdmin)
	totalPages := (int(totalItems) + size - 1) / size

	return model.NewResponse(todoResponses, &model.PageMetadata{
		Page:       page,
		Size:       size,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	})
}
