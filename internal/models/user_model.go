package models

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,gte=8,lte=255"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,gte=8,lte=255"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,jwt"`
}
