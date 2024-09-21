package services

import "github.com/savioruz/mikti-task-1/internal/cores/ports"

type MathService struct {
	MathRepository ports.MathRepository
}

func NewMathService(mathRepository ports.MathRepository) *MathService {
	return &MathService{MathRepository: mathRepository}
}

func (s *MathService) Add(num ...float64) (*float64, error) {
	return s.MathRepository.Add(num...)
}

func (s *MathService) Subtract(num ...float64) (*float64, error) {
	return s.MathRepository.Subtract(num...)
}

func (s *MathService) Multiply(num ...float64) (*float64, error) {
	return s.MathRepository.Multiply(num...)
}

func (s *MathService) Divide(num ...float64) (*float64, error) {
	return s.MathRepository.Divide(num...)
}

func (s *MathService) Factorial(num int) (*int, error) {
	return s.MathRepository.Factorial(num)
}

func (s *MathService) Fibonacci(num int) (*int, error) {
	return s.MathRepository.Fibonacci(num)
}
