package converter

import "github.com/savioruz/mikti-task/tree/post-1/pkg/models"

func BookToResponse(book *models.Book) *models.BookResponse {
	return &models.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Year:        book.Year,
		BookStatus:  book.BookStatus,
		Picture:     book.Picture,
		Description: book.Description,
		Rating:      book.Rating,
		CreatedAt:   book.CreatedAt.String(),
		UpdatedAt:   book.UpdatedAt.String(),
	}
}
