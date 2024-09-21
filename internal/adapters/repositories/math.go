package repositories

type MathRepository struct{}

func NewMathRepository() *MathRepository {
	return &MathRepository{}
}

func (r *MathRepository) Add(num ...float64) (*float64, error) {
	if len(num) == 0 {
		return nil, ErrNoNumberProvided
	}

	var result float64
	for _, n := range num {
		result += n
	}
	return &result, nil
}

func (r *MathRepository) Subtract(num ...float64) (*float64, error) {
	if len(num) == 0 {
		return nil, ErrNoNumberProvided
	}

	var result float64
	for i, n := range num {
		if i == 0 {
			result = n
			continue
		}
		result -= n
	}
	return &result, nil
}

func (r *MathRepository) Multiply(num ...float64) (*float64, error) {
	if len(num) == 0 {
		return nil, ErrNoNumberProvided
	}

	var result float64
	for i, n := range num {
		if i == 0 {
			result = n
			continue
		}
		result *= n
	}
	return &result, nil
}

func (r *MathRepository) Divide(num ...float64) (*float64, error) {
	if len(num) == 0 {
		return nil, ErrNoNumberProvided
	}

	var result float64
	for i, n := range num {
		if i == 0 {
			result = n
			continue
		}
		if n == 0 {
			return nil, ErrInfinity
		}
		result /= n
	}
	return &result, nil
}

func (r *MathRepository) Factorial(num int) (*int, error) {
	if num < 0 || num > 20 {
		return nil, ErrInvalidNumber
	}

	temp := func(value int) *int {
		return &value
	}
	if num == 0 || num == 1 {
		return temp(1), nil
	}

	result, err := r.Factorial(num - 1)
	if err != nil {
		return nil, err
	}

	factorial := num * *result
	return &factorial, nil
}

func (r *MathRepository) Fibonacci(num int) (*int, error) {
	if num < 0 || num > 20 {
		return nil, ErrInvalidNumber
	}

	temp := func(value int) *int {
		return &value
	}
	if num == 0 || num == 1 {
		return temp(num), nil
	}

	resultA, err := r.Fibonacci(num - 1)
	if err != nil {
		return nil, err
	}

	resultB, err := r.Fibonacci(num - 2)
	if err != nil {
		return nil, err
	}

	fibonacci := *resultA + *resultB
	return &fibonacci, nil
}
