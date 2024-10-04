package main

import (
	"github.com/savioruz/mikti-task/week-1/internal/adapters/cli"
	"github.com/savioruz/mikti-task/week-1/internal/adapters/repositories"
	"github.com/savioruz/mikti-task/week-1/internal/cores/services"
)

func main() {
	mathRepository := repositories.NewMathRepository()
	mathService := services.NewMathService(mathRepository)

	dataRepository := repositories.NewDataRepository()
	dataService := services.NewDataService(dataRepository)

	c := cli.NewCLI(mathService, dataService)
	c.Run()
}
