package repositories

import (
	"github.com/savioruz/mikti-task-1/internal/cores/entities"
)

type DataRepository struct {
	data map[string]entities.Data
}

func NewDataRepository() *DataRepository {
	return &DataRepository{}
}

func (r *DataRepository) Insert(ID string, Name string, Email string, Hp int) (*[]entities.Data, error) {
	if ID == "" || Name == "" || Email == "" || Hp == 0 {
		return nil, ErrInvalidRequest
	}

	if r.data == nil {
		r.data = make(map[string]entities.Data)
	}

	if _, ok := r.data[ID]; ok {
		return nil, ErrDataExist
	}

	r.data[ID] = entities.Data{
		ID:    ID,
		Name:  Name,
		Email: Email,
		Hp:    Hp,
	}

	data := make([]entities.Data, 0, len(r.data))
	for _, v := range r.data {
		data = append(data, v)
	}

	return &data, nil
}

func (r *DataRepository) GetAll() (*[]entities.Data, error) {
	if r.data == nil {
		return nil, ErrDataNotFound
	}

	data := make([]entities.Data, 0, len(r.data))
	for _, v := range r.data {
		data = append(data, v)
	}

	return &data, nil
}

func (r *DataRepository) GetByID(ID string) (*entities.Data, error) {
	if r.data == nil {
		return nil, ErrDataNotFound
	}

	if data, ok := r.data[ID]; ok {
		return &data, nil
	}

	return nil, ErrDataNotFound
}

func (r *DataRepository) Update(ID string, Name string, Email string, Hp int) (*[]entities.Data, error) {
	if ID == "" || Name == "" || Email == "" || Hp == 0 {
		return nil, ErrInvalidRequest
	}

	if r.data == nil {
		return nil, ErrDataNotFound
	}

	if _, ok := r.data[ID]; !ok {
		return nil, ErrDataNotFound
	}

	r.data[ID] = entities.Data{
		ID:    ID,
		Name:  Name,
		Email: Email,
		Hp:    Hp,
	}

	data := make([]entities.Data, 0, len(r.data))
	for _, v := range r.data {
		data = append(data, v)
	}

	return &data, nil
}

func (r *DataRepository) Delete(ID string) (*[]entities.Data, error) {
	if r.data == nil {
		return nil, ErrDataNotFound
	}

	if _, ok := r.data[ID]; !ok {
		return nil, ErrDataNotFound
	}

	delete(r.data, ID)

	data := make([]entities.Data, 0, len(r.data))
	for _, v := range r.data {
		data = append(data, v)
	}

	return &data, nil
}
