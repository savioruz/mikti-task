package cli

import (
	"fmt"
	"github.com/savioruz/mikti-task-1/internal/cores/ports"
)

type CLI struct {
	mathService ports.MathRepository
	dataService ports.DataRepository
}

func NewCLI(mathService ports.MathRepository, dataService ports.DataRepository) *CLI {
	return &CLI{
		mathService: mathService,
		dataService: dataService,
	}
}

func (c *CLI) Run() {
	menu := []string{
		"[1]. Hello World",
		"[2]. Basic Math Operations",
		"[3]. Slice & Map Data",
		"[4]. Recursive Operations",
		"[0]. Exit",
	}

	for {
		c.ClearScreen()
		c.DrawBox(menu, "Menu")
		choice := c.GetUserInput("Enter your choice")

		switch choice {
		case "1":
			c.GetOutput("Hello World !")
		case "2":
			c.mathMenu()
		case "3":
			c.dataMenu()
		case "4":
			c.recursiveMenu()
		case "0":
			c.GetOutputWithoutLines("K bye :(\n")
			return
		default:
			c.GetOutputWithoutLines("Invalid choice")
		}

		fmt.Println()
		c.GetUserInput("Press enter to continue...")
	}
}
