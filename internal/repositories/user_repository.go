package repositories

import (
	"github.com/savioruz/mikti-task/tree/week-3/internal/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entities.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByToken(db *gorm.DB, user *entities.User, token string) error {
	return db.Where("token = ?", token).First(user).Error
}
