package user

import (
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/repositories"
	"gorm.io/gorm"
)

type UserRepository interface {
	repositories.Repository[entity.User]
	GetFirst(db *gorm.DB, user *entity.User) error
	GetByID(db *gorm.DB, user *entity.User, id string) error
	GetByEmail(db *gorm.DB, user *entity.User, email string) error
	CountByRole(db *gorm.DB, role string) (int64, error)
}
