package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"context"
	"github.com/savioruz/mikti-task/internal/delivery/graph"
	graphmodel "github.com/savioruz/mikti-task/internal/delivery/graph/model"
	"github.com/savioruz/mikti-task/internal/domain/model"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, title string) (*model.TodoResponse, error) {
	return r.TodoUsecase.Create(ctx, &model.TodoCreateRequest{Title: title})
}

// UpdateTodo is the resolver for the updateTodo field.
func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, input model.TodoUpdateRequest) (*model.TodoResponse, error) {
	u, err := r.TodoUsecase.Update(ctx, &model.TodoUpdateIDRequest{ID: id}, &input)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// DeleteTodo is the resolver for the deleteTodo field.
func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (bool, error) {
	d, err := r.TodoUsecase.Delete(ctx, &model.TodoDeleteRequest{ID: id})
	if err != nil {
		return false, err
	}

	return d, nil
}

// Todo is the resolver for the todo field.
func (r *queryResolver) Todo(ctx context.Context, id string) (*model.TodoResponse, error) {
	return r.TodoUsecase.Get(ctx, &model.TodoGetRequest{ID: id})
}

// SearchTodos is the resolver for the searchTodos field.
func (r *queryResolver) SearchTodos(ctx context.Context, title *string, page *int, size *int) (*graphmodel.TodoResponse, error) {
	var paginated *model.Response[[]*model.TodoResponse]
	var err error
	if page != nil && size != nil {
		paginated, err = r.TodoUsecase.Search(ctx, &model.TodoSearchRequest{
			Title: *title,
			Page:  *page,
			Size:  *size,
		})
	} else {
		paginated, err = r.TodoUsecase.Search(ctx, &model.TodoSearchRequest{
			Title: *title,
			Page:  1,
			Size:  10,
		})
	}

	if err != nil {
		return nil, err
	}

	return &graphmodel.TodoResponse{
		Data: *paginated.Data,
		Paging: &graphmodel.PageMetadata{
			Page:       paginated.Paging.Page,
			Size:       paginated.Paging.Size,
			TotalItems: paginated.Paging.TotalItems,
			TotalPages: paginated.Paging.TotalPages,
		},
	}, nil
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context, page *int, size *int) (*graphmodel.TodoResponse, error) {
	var paginated *model.Response[[]*model.TodoResponse]
	var err error
	if page != nil && size != nil {
		paginated, err = r.TodoUsecase.GetAll(ctx, &model.TodoGetAllRequest{
			Page: *page,
			Size: *size,
		})
	} else {
		paginated, err = r.TodoUsecase.GetAll(ctx, &model.TodoGetAllRequest{
			Page: 1,
			Size: 10,
		})
	}

	if err != nil {
		return nil, err
	}

	return &graphmodel.TodoResponse{
		Data: *paginated.Data,
		Paging: &graphmodel.PageMetadata{
			Page:       paginated.Paging.Page,
			Size:       paginated.Paging.Size,
			TotalItems: paginated.Paging.TotalItems,
			TotalPages: paginated.Paging.TotalPages,
		},
	}, nil
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
