package cli

import (
	"fmt"
	"github.com/google/uuid"
	toxic "github.com/savioruz/mikti-task/week-1/pkg/word"
	"strconv"
	"strings"
)

func (c *CLI) dataMenu() {
	menu := []string{
		"[1]. Insert",
		"[2]. Get All",
		"[3]. Get By ID",
		"[4]. Update",
		"[5]. Delete",
		"[0]. Back",
	}

	for {
		c.ClearScreen()
		c.DrawBox(menu, "Data Operations")
		choice := c.GetUserInput("Enter your choice")

		var sub string
		switch choice {
		case "1":
			sub = "insert"
		case "2":
			sub = "getAll"
		case "3":
			sub = "getByID"
		case "4":
			sub = "update"
		case "5":
			sub = "delete"
		case "0":
			return
		default:
			c.GetOutput("Invalid choice")
		}

		c.doData(sub)
		c.GetUserInput("Press enter to continue")
	}
}

func (c *CLI) doData(sub string) {
	switch sub {
	case "insert":
		name := c.GetUserInput("Enter Name: ")
		email := c.GetUserInput("Enter Email: ")
		hp := c.GetUserInput("Enter HP: ")

		if !c.validate(name, email, hp) {
			return
		}

		hpParse, err := strconv.Atoi(hp)
		if err != nil {
			c.GetOutputWithoutLines("Invalid HP")
			return
		}

		id := uuid.New().String()[:5]
		_, err = c.dataService.Insert(id, name, email, hpParse)

		if err != nil {
			c.GetOutput(fmt.Sprintf("Error: %v", err))
			return
		}
		c.GetOutputWithoutLines("Data inserted")
	case "getAll":
		all, err := c.dataService.GetAll()
		if err != nil {
			return
		}
		c.GetOutputWithoutLines("Data:\n\n")
		c.GetOutputWithoutLines(fmt.Sprintf("ID%sName%50sEmail%25sHP\n", "\t", "\t", "\t"))
		for _, data := range *all {
			c.GetOutputWithoutLines(fmt.Sprintf("%s\t%s\t%50s\t%25d\n", data.ID, data.Name, data.Email, data.Hp))
		}
		fmt.Println()
	case "getByID":
		id := c.GetUserInput("Enter ID: ")
		if id == "" {
			c.GetOutput("Invalid input\n")
			return
		}

		data, err := c.dataService.GetByID(id)
		if err != nil {
			c.GetOutputWithoutLines(fmt.Sprintf("Error: %v", err))
		}

		c.GetOutputWithoutLines("Data:\n\n")
		c.GetOutput("ID\tName\tEmail\tHP\n")
		c.GetOutput(fmt.Sprintf("%s\t%s\t%s\t%d\n", data.ID, data.Name, data.Email, data.Hp))
		fmt.Println()
	case "update":
		id := c.GetUserInput("Enter ID: ")
		if id == "" {
			c.GetOutputWithoutLines("Invalid input\n")
			return
		}

		name := c.GetUserInput("Enter Name: ")
		email := c.GetUserInput("Enter Email: ")
		hp := c.GetUserInput("Enter HP: ")

		if !c.validate(name, email, hp) {
			return
		}

		hpParse, err := strconv.Atoi(hp)
		if err != nil {
			c.GetOutputWithoutLines("Invalid HP\n")
			return
		}

		_, err = c.dataService.Update(id, name, email, hpParse)
		if err != nil {
			c.GetOutputWithoutLines(fmt.Sprintf("Error: %v", err))
			return
		}
	case "delete":
		id := c.GetUserInput("Enter ID: ")
		if id == "" {
			c.GetOutputWithoutLines("Invalid input\n")
			return
		}

		_, err := c.dataService.Delete(id)
		if err != nil {
			c.GetOutputWithoutLines(fmt.Sprintf("Error: %v", err))
			return
		}

		c.GetOutputWithoutLines(fmt.Sprintf("Data with ID %s deleted", id))
	}
}

func (c *CLI) validate(name, email, hp string) bool {
	filterWords := func(word string) bool {
		lower := strings.ToLower(word)
		for _, w := range toxic.ToxicWords {
			if strings.Contains(lower, w) {
				return true
			}
		}
		return false
	}

	filterEmails := func(email string) bool {
		return strings.Contains(email, "@")
	}

	if name == "" || email == "" || hp == "" {
		c.GetOutputWithoutLines("Invalid input\n")
		return false
	}

	if len(name) > 100 || len(email) > 50 || len(hp) > 14 {
		c.GetOutputWithoutLines("Input too long\n")
		return false
	}

	if filterWords(name) || filterWords(email) {
		c.GetOutputWithoutLines("Contains toxic words\n")
		return false
	}

	if !filterEmails(email) {
		c.GetOutputWithoutLines("Invalid email\n")
		return false
	}

	return true
}
