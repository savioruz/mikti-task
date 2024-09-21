package repositories

import "errors"

var (
	ErrNoNumberProvided = errors.New("no number provided")
	ErrInvalidNumber    = errors.New("invalid number")
	ErrInfinity         = errors.New("infinity")
	ErrInvalidRequest   = errors.New("invalid request")
	ErrDataExist        = errors.New("data already exist")
	ErrDataNotFound     = errors.New("data not found")
)
