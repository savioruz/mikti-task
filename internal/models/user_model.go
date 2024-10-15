package models

type UserResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Token     string `json:"token,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type VerifyUserRequest struct {
	Token string `json:"token" validate:"required,max=100"`
}

type RegisterUserRequest struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required,min=5"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UpdateUserRequest struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name,omitempty" validate:"omitempty,min=5"`
	Pass string `json:"password,omitempty" validate:"omitempty,min=8,max=100"`
}

type LoginUserRequest struct {
	ID       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required,max=100"`
}

type LogoutUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type GetUserRequest struct {
	ID string `param:"id" validate:"required,max=100"`
}
