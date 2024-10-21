package converter

import (
	"github.com/savioruz/mikti-task/tree/week-4/internal/entities"
	"github.com/savioruz/mikti-task/tree/week-4/internal/models"
)

func UserToResponse(user *entities.User) *models.UserResponse {
	return &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}

func LoginToTokenResponse(accessToken, refreshToken string) *models.TokenResponse {
	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
