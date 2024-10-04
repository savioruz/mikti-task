package services

import (
	"github.com/savioruz/mikti-task/week-1/internal/cores/entities"
	"github.com/savioruz/mikti-task/week-1/internal/cores/ports"
)

type DataRepository struct {
	DataRepository ports.DataRepository
}

func NewDataService(dataRepository ports.DataRepository) *DataRepository {
	return &DataRepository{DataRepository: dataRepository}
}

func (s *DataRepository) Insert(ID string, Name string, Email string, Hp int) (*[]entities.Data, error) {
	return s.DataRepository.Insert(ID, Name, Email, Hp)
}

func (s *DataRepository) GetAll() (*[]entities.Data, error) {
	return s.DataRepository.GetAll()
}

func (s *DataRepository) GetByID(ID string) (*entities.Data, error) {
	return s.DataRepository.GetByID(ID)
}

func (s *DataRepository) Update(ID string, Name string, Email string, Hp int) (*[]entities.Data, error) {
	return s.DataRepository.Update(ID, Name, Email, Hp)
}

func (s *DataRepository) Delete(ID string) (*[]entities.Data, error) {
	return s.DataRepository.Delete(ID)
}
