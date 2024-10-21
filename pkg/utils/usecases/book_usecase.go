package usecases

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/models"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/utils"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/utils/converter"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

type BookUseCase struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate
	Query    *utils.Query[models.Book]
}

func NewBookUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, query *utils.Query[models.Book]) *BookUseCase {
	return &BookUseCase{
		DB:       db,
		Log:      log,
		Validate: validate,
		Query:    query,
	}
}

func (b *BookUseCase) CreateBook(ctx context.Context, request *models.BookCreateRequest) (*models.BookResponse, error) {
	tx := b.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := b.Validate.Struct(request); err != nil {
		b.Log.Errorf("failed to validate request: %v", err)
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	book := &models.Book{
		ID:          uuid.New().String(),
		Title:       request.Title,
		Author:      request.Author,
		BookStatus:  request.BookStatus,
		Year:        request.Year,
		Picture:     request.Picture,
		Description: request.Description,
		Rating:      request.Rating,
	}

	if err := b.Query.FindByField(tx, book, "title", request.Title); err == nil {
		b.Log.Errorf("book already exists")
		return nil, errors.New(http.StatusText(http.StatusConflict))
	}

	if err := b.Query.Create(tx, book); err != nil {
		b.Log.Errorf("failed to create book: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		b.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.BookToResponse(book), nil
}

func (b *BookUseCase) UpdateBook(ctx context.Context, id *models.BookIdRequest, request *models.BookUpdateRequest) (*models.BookResponse, error) {
	tx := b.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.Title == "" && request.Author == "" && request.Year == 0 && request.BookStatus == 0 {
		return nil, errors.New("no field to update")
	}

	if err := b.Validate.Struct(id); err != nil {
		b.Log.Errorf("failed to validate id: %v", err)
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	if err := b.Validate.Struct(request); err != nil {
		b.Log.Errorf("failed to validate request: %v", err)
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	book := &models.Book{}
	if err := b.Query.FindByID(tx, book, id.ID); err != nil {
		b.Log.Errorf("failed to find book: %v", err)
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	if request != nil {
		if request.Title != "" {
			book.Title = request.Title
		}
		if request.Author != "" {
			book.Author = request.Author
		}
		if request.Year != 0 {
			book.Year = request.Year
		}
		if request.BookStatus != 0 {
			book.BookStatus = request.BookStatus
		}
		if request.Picture != "" {
			book.Picture = request.Picture
		}
		if request.Description != "" {
			book.Description = request.Description
		}
		if request.Rating != 0 {
			book.Rating = request.Rating
		}
	}

	if err := b.Query.Update(tx, book); err != nil {
		b.Log.Errorf("failed to update book: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		b.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.BookToResponse(book), nil
}

func (b *BookUseCase) DeleteBook(ctx context.Context, request *models.BookIdRequest) error {
	tx := b.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := b.Validate.Struct(request); err != nil {
		b.Log.Errorf("failed to validate request: %v", err)
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	book := &models.Book{}
	if err := b.Query.FindByID(tx, book, request.ID); err != nil {
		b.Log.Errorf("failed to find book: %v", err)
		return errors.New(http.StatusText(http.StatusNotFound))
	}

	if err := b.Query.Delete(tx, book); err != nil {
		b.Log.Errorf("failed to delete book: %v", err)
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		b.Log.Errorf("failed to commit transaction: %v", err)
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func (b *BookUseCase) GetBook(ctx context.Context, request *models.BookIdRequest) (*models.BookResponse, error) {
	tx := b.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := b.Validate.Struct(request); err != nil {
		b.Log.Errorf("failed to validate request: %v", err)
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	book := &models.Book{}
	if err := b.Query.FindByID(tx, book, request.ID); err != nil {
		b.Log.Errorf("failed to find book: %v", err)
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	return converter.BookToResponse(book), nil
}

func (b *BookUseCase) GetBooks(ctx context.Context) ([]*models.BookResponse, error) {
	tx := b.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var books []models.Book
	if err := b.Query.FindAll(tx, &books); err != nil {
		b.Log.Errorf("failed to find books: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	var response []*models.BookResponse
	for _, book := range books {
		response = append(response, converter.BookToResponse(&book))
	}

	return response, nil
}
