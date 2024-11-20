package model

type TodoResponse struct {
	ID        string  `json:"id"`
	UserID    *string `json:"user_id,omitempty"`
	Title     string  `json:"title"`
	Done      bool    `json:"done"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type TodoCreateRequest struct {
	Title string `json:"title" validate:"required,gte=5,lte=255"`
}

type TodoUpdateIDRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

type TodoUpdateRequest struct {
	Title *string `json:"title,omitempty" validate:"omitempty,gte=5,lte=255"`
	Done  *bool   `json:"done,omitempty" validate:"omitempty,boolean"`
}

type TodoDeleteRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

type TodoGetRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

type TodoSearchRequest struct {
	Title string `query:"title" validate:"omitempty,gte=2,lte=255"`
	Page  int    `query:"page" validate:"numeric"`
	Size  int    `query:"size" validate:"numeric"`
}

type TodoGetAllRequest struct {
	Page int `query:"page" validate:"numeric"`
	Size int `query:"size" validate:"numeric"`
}

type TodoQueryOptions struct {
	UserID  *string
	Title   *string
	Page    int
	Size    int
	IsAdmin bool
}
