package main

import (
	"github.com/savioruz/mikti-task/week-2/internal/adapters/cli"
	"github.com/savioruz/mikti-task/week-2/internal/adapters/repositories"
	"github.com/savioruz/mikti-task/week-2/internal/cores/services"
)

func main() {
	itemRepository := repositories.NewItemRepository()
	itemService := services.NewItemService(itemRepository)

	c := cli.NewCLI(itemService)
	c.Run()
}
