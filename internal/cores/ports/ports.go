package ports

import "github.com/savioruz/mikti-task-1/internal/cores/entities"

type MathRepository interface {
	Add(num ...float64) (*float64, error)
	Subtract(num ...float64) (*float64, error)
	Multiply(num ...float64) (*float64, error)
	Divide(num ...float64) (*float64, error)
	Factorial(num int) (*int, error)
	Fibonacci(num int) (*int, error)
}

type DataRepository interface {
	Insert(ID string, Name string, Email string, Hp int) (*[]entities.Data, error)
	GetAll() (*[]entities.Data, error)
	GetByID(ID string) (*entities.Data, error)
	Update(ID string, Name string, Email string, Hp int) (*[]entities.Data, error)
	Delete(ID string) (*[]entities.Data, error)
}
