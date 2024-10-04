package cli

import (
	"fmt"
	"github.com/savioruz/mikti-task/week-1/internal/adapters/repositories"
	"strconv"
	"strings"
)

func (c *CLI) mathMenu() {
	menu := []string{
		"[1]. Add",
		"[2]. Subtract",
		"[3]. Multiply",
		"[4]. Divide",
		"[5]. All",
		"[0]. Back",
	}

	for {
		c.ClearScreen()
		c.DrawBox(menu, "Math Operations")
		choice := c.GetUserInput("Enter your choice")

		var sub string
		switch choice {
		case "1":
			sub = "add"
		case "2":
			sub = "subtract"
		case "3":
			sub = "multiply"
		case "4":
			sub = "divide"
		case "5":
			sub = "all"
		case "0":
			return
		default:
			c.GetOutput("Invalid choice")
		}

		c.doMath(sub)
		c.GetUserInput("Press enter to continue")
	}
}

func (c *CLI) doMath(sub string) {
	var nums *[]float64
	var err error
	if sub == "all" {
		nums, err = c.userInput()
		if err != nil {
			c.GetOutputWithoutLines(fmt.Sprintf("Error: %v\n", err))
			return
		}
	}

	calculate := func(operation string, services func(...float64) (*float64, error)) {
		var result *float64
		var err error

		if sub != "all" {
			nums, err = c.userInput()
			if err != nil {
				c.GetOutputWithoutLines(fmt.Sprintf("Error: %v\n", err))
				return
			}
		}

		result, err = services(*nums...)
		if err != nil {
			c.GetOutputWithoutLines(fmt.Sprintf("Error: %v\n", err))
			return
		}
		c.GetOutput(fmt.Sprintf("%s Result: %v", operation, *result))
	}

	switch sub {
	case "add":
		calculate("Addition", c.mathService.Add)
	case "subtract":
		calculate("Subtraction", c.mathService.Subtract)
	case "multiply":
		calculate("Multiplication", c.mathService.Multiply)
	case "divide":
		calculate("Division", c.mathService.Divide)
	case "all":
		calculate("Addition", c.mathService.Add)
		c.OneLineUp()
		calculate("Subtraction", c.mathService.Subtract)
		c.OneLineUp()
		calculate("Multiplication", c.mathService.Multiply)
		c.OneLineUp()
		calculate("Division", c.mathService.Divide)
	default:
		c.GetOutputWithoutLines("Invalid operation")
	}
}

func (c *CLI) userInput() (*[]float64, error) {
	c.GetOutputWithoutLines("Enter numbers separated by space")
	numbers := c.GetUserInput("")

	if len(strings.Fields(numbers)) < 2 {
		return nil, repositories.ErrInvalidNumber
	}

	nums := make([]float64, 0)
	for _, n := range strings.Fields(numbers) {
		num, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return nil, repositories.ErrInvalidNumber
		}
		nums = append(nums, num)
	}

	return &nums, nil
}
