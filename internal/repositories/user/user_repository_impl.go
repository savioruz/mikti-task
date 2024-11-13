package user

import (
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/repositories"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	repositories.RepositoryImpl[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(db *gorm.DB, log *logrus.Logger) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		RepositoryImpl: repositories.RepositoryImpl[entity.User]{DB: db},
		Log:            log,
	}
}

func (r *UserRepositoryImpl) GetByID(db *gorm.DB, user *entity.User, id string) error {
	return db.Where("id = ?", id).Take(&user).Error
}

func (r *UserRepositoryImpl) GetByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).Take(&user).Error
}

func (r *UserRepositoryImpl) CountByRole(db *gorm.DB, role string) (int64, error) {
	var count int64
	err := db.Model(&entity.User{}).Where("role = ?", role).Count(&count).Error
	return count, err
}
