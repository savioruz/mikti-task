package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/models"
	"github.com/savioruz/mikti-task/tree/post-1/pkg/utils/usecases"
	"github.com/sirupsen/logrus"
	"net/http"
)

type BookController struct {
	Log     *logrus.Logger
	Usecase *usecases.BookUseCase
}

func NewBookController(log *logrus.Logger, usecase *usecases.BookUseCase) *BookController {
	return &BookController{
		Log:     log,
		Usecase: usecase,
	}
}

// CreateBook creates a new book
// @Summary Create a new book
// @Description Create a new book
// @Tags book
// @Accept json
// @Produce json
// @Param request body models.BookCreateRequest true "Book Create Request"
// @Success 201 {object} object{data=models.BookResponse}
// @Failure 400 {object} object{error=string}
// @Failure 409 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /book [post]
func (b *BookController) CreateBook(c echo.Context) error {
	request := new(models.BookCreateRequest)
	if err := c.Bind(request); err != nil {
		b.Log.Errorf("failed to bind request: %v", err)
		return HandleError(c, http.StatusBadRequest, ErrorBindingRequest)
	}

	response, err := b.Usecase.CreateBook(c.Request().Context(), request)
	if err != nil {
		b.Log.Errorf("failed to create book: %v", err)
		if err.Error() == http.StatusText(http.StatusBadRequest) {
			return HandleError(c, http.StatusBadRequest, ErrorValidation)
		} else if err.Error() == http.StatusText(http.StatusConflict) {
			return HandleError(c, http.StatusConflict, ErrorConflict)
		} else {
			return HandleError(c, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return c.JSON(http.StatusCreated, models.WebResponse[*models.BookResponse]{
		Data: response,
	})
}

// GetBookByID gets a book by ID
// @Summary Get a book by ID
// @Description Get a book by ID
// @Tags book
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} object{data=models.BookResponse}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /book/{id} [get]
func (b *BookController) GetBookByID(c echo.Context) error {
	id := &models.BookIdRequest{ID: c.Param("id")}

	response, err := b.Usecase.GetBook(c.Request().Context(), id)
	if err != nil {
		b.Log.Errorf("failed to get book by ID: %v", err)
		if err.Error() == http.StatusText(http.StatusNotFound) {
			return HandleError(c, http.StatusNotFound, ErrorNotFound)
		} else {
			return HandleError(c, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return c.JSON(http.StatusOK, models.WebResponse[*models.BookResponse]{
		Data: response,
	})
}

// GetAllBooks gets all books
// @Summary Get all books
// @Description Get all books
// @Tags book
// @Accept json
// @Produce json
// @Success 200 {object} object{data=[]models.BookResponse}
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /books [get]
func (b *BookController) GetAllBooks(c echo.Context) error {
	request := new(models.GetsRequest)
	if err := c.Bind(request); err != nil {
		b.Log.Errorf("failed to bind request: %v", err)
		return HandleError(c, http.StatusBadRequest, ErrorBindingRequest)
	}

	response, err := b.Usecase.GetBooks(c.Request().Context())
	if err != nil {
		b.Log.Errorf("failed to get books: %v", err)
		return HandleError(c, http.StatusInternalServerError, ErrorInternalServer)
	}

	return c.JSON(http.StatusOK, models.WebResponse[[]*models.BookResponse]{
		Data: response,
	})
}

// UpdateBook updates a book by ID
// @Summary Update a book by ID
// @Description Update a book by ID
// @Tags book
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param request body models.BookUpdateRequest true "Book Update Request"
// @Success 200 {object} models.BookResponse
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /book/{id} [put]
func (b *BookController) UpdateBook(c echo.Context) error {
	id := &models.BookIdRequest{ID: c.Param("id")}

	request := new(models.BookUpdateRequest)
	if err := c.Bind(request); err != nil {
		b.Log.Errorf("failed to bind request: %v", err)
		return HandleError(c, http.StatusBadRequest, ErrorBindingRequest)
	}

	response, err := b.Usecase.UpdateBook(c.Request().Context(), id, request)
	if err != nil {
		b.Log.Errorf("failed to update book: %v", err)
		if err.Error() == http.StatusText(http.StatusBadRequest) {
			return HandleError(c, http.StatusBadRequest, ErrorValidation)
		} else if err.Error() == http.StatusText(http.StatusNotFound) {
			return HandleError(c, http.StatusNotFound, err)
		} else {
			return HandleError(c, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return c.JSON(http.StatusOK, models.WebResponse[*models.BookResponse]{
		Data: response,
	})
}

// DeleteBook deletes a book by ID
// @Summary Delete a book by ID
// @Description Delete a book by ID
// @Tags book
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /book/{id} [delete]
func (b *BookController) DeleteBook(c echo.Context) error {
	id := &models.BookIdRequest{ID: c.Param("id")}

	err := b.Usecase.DeleteBook(c.Request().Context(), id)
	if err != nil {
		b.Log.Errorf("failed to delete book: %v", err)
		if err.Error() == http.StatusText(http.StatusBadRequest) {
			return HandleError(c, http.StatusBadRequest, ErrorValidation)
		} else if err.Error() == http.StatusText(http.StatusNotFound) {
			return HandleError(c, http.StatusNotFound, err)
		} else {
			return HandleError(c, http.StatusInternalServerError, ErrorInternalServer)
		}
	}

	return c.JSON(http.StatusNoContent, nil)
}
