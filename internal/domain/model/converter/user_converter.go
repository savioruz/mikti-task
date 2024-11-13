package converter

import (
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/domain/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}

func LoginToTokenResponse(accessToken, refreshToken string) *model.TokenResponse {
	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
