package converter

import (
	"github.com/savioruz/mikti-task/tree/week-3/internal/entities"
	"github.com/savioruz/mikti-task/tree/week-3/internal/models"
)

func UserToResponse(user *entities.User) *models.UserResponse {
	return &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}

func UserToTokenResponse(user *entities.User) *models.UserResponse {
	return &models.UserResponse{
		Token: user.Token,
	}
}
