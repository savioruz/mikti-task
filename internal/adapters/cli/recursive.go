package cli

import (
	"fmt"
	"strconv"
)

func (c *CLI) recursiveMenu() {
	menu := []string{
		"[1]. Factorial",
		"[2]. Fibonacci",
		"[3]. All",
		"[0]. Back",
	}

	for {
		c.ClearScreen()
		c.DrawBox(menu, "Recursive Operations")
		choice := c.GetUserInput("Enter your choice")

		var sub string
		switch choice {
		case "1":
			sub = "factorial"
		case "2":
			sub = "fibonacci"
		case "3":
			sub = "all"
		case "0":
			return
		default:
			c.GetOutput("Invalid choice")
		}

		c.doRecursive(sub)
		c.GetUserInput("Press enter to continue")
	}
}

func (c *CLI) doRecursive(sub string) {
	var num int
	var input string
	if sub == "all" {
		c.GetOutputWithoutLines("Enter a number: ")
		input = c.GetUserInput("")
		if input == "" {
			c.GetOutputWithoutLines("Invalid input")
			return
		}
	}

	num, err := strconv.Atoi(input)
	if err != nil {
		c.GetOutputWithoutLines("Invalid number")
		return
	}

	calculate := func(operation string, services func(int) (*int, error)) {
		if sub != "all" {
			c.GetOutputWithoutLines("Enter a number: ")
			input = c.GetUserInput("")
			if input == "" {
				c.GetOutputWithoutLines("Invalid input")
				return
			}
		}

		result, err := services(num)
		if err != nil {
			c.GetOutputWithoutLines(fmt.Sprintf("Error: %v\n", err))
			return
		}
		c.GetOutput(fmt.Sprintf("%s Result: %v", operation, *result))
	}

	switch sub {
	case "factorial":
		calculate("Factorial", c.mathService.Factorial)
	case "fibonacci":
		calculate("Fibonacci", c.mathService.Fibonacci)
	case "all":
		calculate("Factorial", c.mathService.Factorial)
		c.OneLineUp()
		calculate("Fibonacci", c.mathService.Fibonacci)
	default:
		c.GetOutput("Invalid choice")
	}

}
