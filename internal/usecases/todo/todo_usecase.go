package todo

import (
	"context"
	"github.com/savioruz/mikti-task/internal/domain/model"
)

type TodoUsecase interface {
	Create(ctx context.Context, request *model.TodoCreateRequest) (*model.TodoResponse, error)
	Update(ctx context.Context, request *model.TodoUpdateIDRequest, update *model.TodoUpdateRequest) (*model.TodoResponse, error)
	Get(ctx context.Context, request *model.TodoGetRequest) (*model.TodoResponse, error)
	GetAll(ctx context.Context, request *model.TodoGetAllRequest) (*model.Response[[]*model.TodoResponse], error)
	Delete(ctx context.Context, request *model.TodoDeleteRequest) (bool, error)
}
