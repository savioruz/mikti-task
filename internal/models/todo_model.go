package models

type TodoResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TodoCreateRequest struct {
	Title string `json:"title" validate:"required,alphanum,min=5"`
}

type TodoUpdateIDRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

type TodoUpdateRequest struct {
	Title     *string `json:"title,omitempty" validate:"alphanum,min=5"`
	Completed *bool   `json:"completed,omitempty" validate:"boolean"`
}

type TodoDeleteRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

type TodoGetRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

type TodoListRequest struct {
	Page int `query:"page" validate:"numeric"`
	Size int `query:"size" validate:"numeric"`
}
