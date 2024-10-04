package cli

import (
	"fmt"
	"github.com/savioruz/mikti-task/week-2/internal/cores/entities"
	"github.com/savioruz/mikti-task/week-2/internal/cores/ports"
	"time"
)

type CLI struct {
	ItemRepository ports.ItemRepository
}

func NewCLI(itemRepository ports.ItemRepository) *CLI {
	return &CLI{
		ItemRepository: itemRepository,
	}
}

func (c *CLI) Run() {
	c.ItemRepository.AddItem("Nasi Goreng", 11000)
	c.ItemRepository.AddItem("Mie Goreng", 12000)
	c.ItemRepository.AddItem("Ayam Bakar", 15000)

	c.printMenu(c.ItemRepository.GetItems())

	orderChannel := make(chan entities.Order)

	restaurant := &entities.Restaurant{
		Menu: *c.ItemRepository.GetItems(),
	}

	go func() {
		c.ItemRepository.TakeOrder(restaurant, orderChannel)
	}()

	c.ItemRepository.WaitOrder()

	order := <-orderChannel

	fmt.Println("\nYour Order:")
	for _, item := range order.Items {
		fmt.Printf("- %s (%d)\n", item.Name, item.Quantity)
	}

	total := c.calculateTotal(order)
	fmt.Printf("Total: %.2fIDR\n\n", total)

	encodedOrder := c.ItemRepository.EncodeOrder(order)
	fmt.Printf("Order (encoded base64):\n%s\n\n", encodedOrder)

	c.handlePayment(total)

	c.processOrder()

	fmt.Println("Done...")
}

func (c *CLI) printMenu(menu *[]entities.Item) {
	fmt.Println("Menu:")
	for _, item := range *menu {
		fmt.Printf(" - %s: %.2fIDR\n", item.Name, item.Price)
	}
}

func (c *CLI) calculateTotal(order entities.Order) float64 {
	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func (c *CLI) handlePayment(total float64) {
	for {
		fmt.Print("Input payment: ")
		var paymentInput string
		fmt.Scanln(&paymentInput)

		payment, err := c.ItemRepository.ValidatePrice(paymentInput)
		if err != nil {
			fmt.Println("Invalid input. Please input a valid number.")
			continue
		}

		if payment >= total {
			change := payment - total
			fmt.Printf("Change: %.2fIDR\n", change)
			break
		} else {
			fmt.Println("Insufficient payment. Please input a valid number.")
		}
	}
}

func (c *CLI) processOrder() {
	fmt.Print("Processing...")
	done := make(chan bool)
	go func() {
		time.Sleep(3 * time.Second)
		done <- true
	}()
	<-done
	fmt.Println(" Done!")
}
