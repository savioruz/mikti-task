package user

import (
	"context"
	"github.com/savioruz/mikti-task/tree/week-4/internal/domain/model"
)

type UserUsecase interface {
	Create(ctx context.Context, request *model.RegisterRequest) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginRequest) (*model.TokenResponse, error)
	RefreshToken(request *model.RefreshTokenRequest) (*model.TokenResponse, error)
}
