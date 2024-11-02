package usecases

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/savioruz/mikti-task/tree/week-4/internal/delivery/http/auth"
	"github.com/savioruz/mikti-task/tree/week-4/internal/entities"
	"github.com/savioruz/mikti-task/tree/week-4/internal/models"
	"github.com/savioruz/mikti-task/tree/week-4/internal/models/converter"
	"github.com/savioruz/mikti-task/tree/week-4/internal/repositories"
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
	JWTService     *auth.JWTService
}

func NewUserUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repositories.UserRepository, jwtService *auth.JWTService) *UserUsecase {
	return &UserUsecase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
		JWTService:     jwtService,
	}
}

func (u *UserUsecase) Create(ctx context.Context, request *models.RegisterRequest) (*models.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	existingUser := &entities.User{}
	if err := u.UserRepository.GetByEmail(tx, existingUser, request.Email); err == nil {
		u.Log.Errorf("email already exists: %v", request.Email)
		return nil, errors.New(http.StatusText(http.StatusConflict))
	}

	if request.Role == "admin" {
		count, err := u.UserRepository.CountByRole(tx, request.Role)
		if err != nil {
			u.Log.Errorf("failed to count user by role: %v", err)
			return nil, errors.New(http.StatusText(http.StatusInternalServerError))
		}
		if count > 0 {
			u.Log.Errorf("admin already exists")
			return nil, errors.New(http.StatusText(http.StatusConflict))
		}
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.Errorf("failed to hash password: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	user := &entities.User{
		ID:       uuid.New().String(),
		Email:    request.Email,
		Password: string(password),
		Role:     request.Role,
		Status:   true,
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

func (u *UserUsecase) Login(ctx context.Context, request *models.LoginRequest) (*models.TokenResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	user := &entities.User{}
	if err := u.UserRepository.GetByEmail(tx, user, request.Email); err != nil {
		u.Log.Errorf("failed to get user by email: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		u.Log.Errorf("failed to compare password: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	accessToken, err := u.JWTService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		u.Log.Errorf("failed to generate access token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	refreshToken, err := u.JWTService.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		u.Log.Errorf("failed to generate refresh token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.LoginToTokenResponse(accessToken, refreshToken), nil
}

func (u *UserUsecase) RefreshToken(request *models.RefreshTokenRequest) (*models.TokenResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Errorf("failed to validate refresh token request: %v", err)
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	claims, err := u.JWTService.ValidateToken(request.RefreshToken)
	if err != nil {
		u.Log.Errorf("failed to validate refresh token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	accessToken, err := u.JWTService.GenerateAccessToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		u.Log.Errorf("failed to generate access token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	refreshToken, err := u.JWTService.GenerateRefreshToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		u.Log.Errorf("failed to generate refresh token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.LoginToTokenResponse(accessToken, refreshToken), nil
}
