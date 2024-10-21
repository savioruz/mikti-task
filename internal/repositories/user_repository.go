package repositories

import (
	"github.com/savioruz/mikti-task/tree/week-4/internal/entities"
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

func (r *UserRepository) GetByID(db *gorm.DB, user *entities.User, id string) error {
	return db.Where("id = ?", id).Take(&user).Error
}

func (r *UserRepository) GetByEmail(db *gorm.DB, user *entities.User, email string) error {
	return db.Where("email = ?", email).Take(&user).Error
}

func (r *UserRepository) CountByRole(db *gorm.DB, role string) (int64, error) {
	var count int64
	err := db.Model(&entities.User{}).Where("role = ?", role).Count(&count).Error
	return count, err
}
