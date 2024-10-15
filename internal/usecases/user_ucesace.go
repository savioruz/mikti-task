package usecases

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/savioruz/mikti-task/tree/week-3/internal/entities"
	"github.com/savioruz/mikti-task/tree/week-3/internal/models"
	"github.com/savioruz/mikti-task/tree/week-3/internal/models/converter"
	"github.com/savioruz/mikti-task/tree/week-3/internal/repositories"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type UserUsecase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repositories.UserRepository
}

func NewUserUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repositories.UserRepository) *UserUsecase {
	return &UserUsecase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (u *UserUsecase) Verify(ctx context.Context, request *models.VerifyUserRequest) (*models.Auth, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Errorf("failed to validate request: %v", err)
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	user := new(entities.User)
	if err := u.UserRepository.FindByToken(tx, user, request.Token); err != nil {
		u.Log.Errorf("failed to find user: %v", err)
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return &models.Auth{
		ID: user.ID,
	}, nil
}

func (u *UserUsecase) Login(ctx context.Context, request *models.LoginUserRequest) (*models.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Errorf("failed to validate request: %v", err)
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	user := new(entities.User)
	if err := u.UserRepository.FindByID(tx, user, request.ID); err != nil {
		u.Log.Errorf("failed to find user: %v", err)
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		u.Log.Errorf("failed to compare password: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	user.Token = uuid.New().String()
	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Errorf("failed to save user: %+v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.UserToTokenResponse(user), nil
}

func (u *UserUsecase) Create(ctx context.Context, request *models.RegisterUserRequest) (*models.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Errorf("failed to validate request: %v", err)
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	total, err := u.UserRepository.CountByID(tx, request.ID)
	if err != nil {
		u.Log.Errorf("failed to count user: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if total > 0 {
		u.Log.Errorf("user already exists")
		return nil, errors.New(http.StatusText(http.StatusConflict))
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.Errorf("failed to hash password: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	user := &entities.User{
		ID:       request.ID,
		Name:     request.Name,
		Password: string(password),
	}

	if err := u.UserRepository.Create(tx, user); err != nil {
		u.Log.Errorf("failed to create user: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUsecase) Current(ctx context.Context, request *models.GetUserRequest) (*models.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Errorf("failed to validate request: %v", err)
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	user := new(entities.User)
	if err := u.UserRepository.FindByID(tx, user, request.ID); err != nil {
		u.Log.Errorf("failed to find user: %v", err)
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUsecase) Logout(ctx context.Context, request *models.LogoutUserRequest) (bool, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Errorf("failed to validate request: %v", err)
		return false, errors.New(http.StatusText(http.StatusBadRequest))
	}

	user := new(entities.User)
	if err := u.UserRepository.FindByID(tx, user, request.ID); err != nil {
		u.Log.Errorf("failed to find user: %v", err)
		return false, errors.New(http.StatusText(http.StatusNotFound))
	}

	user.Token = ""
	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Errorf("failed to save user: %+v", err)
		return false, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return false, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return true, nil
}

func (u *UserUsecase) Update(ctx context.Context, request *models.UpdateUserRequest) (*models.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Errorf("failed to validate request: %v", err)
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	user := new(entities.User)
	if err := u.UserRepository.FindByID(tx, user, request.ID); err != nil {
		u.Log.Errorf("failed to find user: %v", err)
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Pass != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Pass), bcrypt.DefaultCost)
		if err != nil {
			u.Log.Errorf("failed to hash password: %v", err)
			return nil, errors.New(http.StatusText(http.StatusInternalServerError))
		}
		user.Password = string(password)
	}

	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Errorf("failed to save user: %+v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.UserToResponse(user), nil
}
