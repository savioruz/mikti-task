package user

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/savioruz/mikti-task/internal/domain/model"
	"github.com/savioruz/mikti-task/internal/domain/model/converter"
	"github.com/savioruz/mikti-task/internal/platform/jwt"
	"github.com/savioruz/mikti-task/internal/repositories/user"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type UserUsecaseImpl struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *user.UserRepositoryImpl
	JWTService     jwt.JWTService
}

func NewUserUsecaseImpl(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *user.UserRepositoryImpl, jwtService jwt.JWTService) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
		JWTService:     jwtService,
	}
}

func (u *UserUsecaseImpl) Create(ctx context.Context, request *model.RegisterRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	existingUser := &entity.User{}
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

	data := &entity.User{
		ID:       uuid.New().String(),
		Email:    request.Email,
		Password: string(password),
		Role:     request.Role,
		Status:   true,
	}

	if err := u.UserRepository.Create(tx, data); err != nil {
		u.Log.Errorf("failed to create user: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.UserToResponse(data), nil
}

func (u *UserUsecaseImpl) Login(ctx context.Context, request *model.LoginRequest) (*model.TokenResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	data := &entity.User{}
	if err := u.UserRepository.GetByEmail(tx, data, request.Email); err != nil {
		u.Log.Errorf("failed to get user by email: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(request.Password)); err != nil {
		u.Log.Errorf("failed to compare password: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	accessToken, err := u.JWTService.GenerateAccessToken(data.ID, data.Email, data.Role)
	if err != nil {
		u.Log.Errorf("failed to generate access token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	refreshToken, err := u.JWTService.GenerateRefreshToken(data.ID, data.Email, data.Role)
	if err != nil {
		u.Log.Errorf("failed to generate refresh token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.LoginToTokenResponse(accessToken, refreshToken), nil
}

func (u *UserUsecaseImpl) RefreshToken(request *model.RefreshTokenRequest) (*model.TokenResponse, error) {
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
