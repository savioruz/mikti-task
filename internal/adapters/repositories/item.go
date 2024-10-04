package repositories

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/savioruz/mikti-task/week-2/internal/cores/entities"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type ItemRepository struct {
	data []entities.Item
	WG   sync.WaitGroup
}

func NewItemRepository() *ItemRepository {
	return &ItemRepository{}
}

func (i *ItemRepository) AddItem(name string, price float64) *[]entities.Item {
	i.data = append(i.data, entities.Item{Name: name, Price: price})
	return &i.data
}

func (i *ItemRepository) GetItems() *[]entities.Item {
	return &i.data
}

func (i *ItemRepository) TakeOrder(restaurant *entities.Restaurant, ch chan<- entities.Order) {
	i.WG.Add(1)
	defer i.WG.Done()

	order := entities.Order{}
	var itemName string
	var itemQuantity int
	s := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Enter item name: ")
		s.Scan()
		itemName = strings.ToLower(s.Text())

		if itemName == "done" {
			break
		}

		if item, ok := i.ValidateOrder(restaurant, itemName); ok {
			fmt.Println("Enter item quantity: ")
			s.Scan()
			fmt.Sscanf(s.Text(), "%d", &itemQuantity)
			order.Items = append(order.Items, entities.Item{
				Name:     item.Name,
				Price:    item.Price,
				Quantity: itemQuantity,
			})
		} else {
			fmt.Println("Item not found in menu")
		}
	}

	ch <- order
}

func (i *ItemRepository) ValidateOrder(restaurant *entities.Restaurant, item string) (*entities.Item, bool) {
	for _, menuItem := range restaurant.Menu {
		if strings.ToLower(menuItem.Name) == item {
			return &menuItem, true
		}
	}
	return nil, false
}

func (i *ItemRepository) ValidatePrice(price string) (float64, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Invalid price")
		}
	}()

	if matched, _ := regexp.MatchString(`^\d+(\.\d{1,2})?$`, price); !matched {
		panic("Invalid price")
	}

	return strconv.ParseFloat(price, 64)
}

func (i *ItemRepository) EncodeOrder(order entities.Order) string {
	data := fmt.Sprintf("Order: %v, Total: %.2fIDR", order.Items, order.Total)
	encodedData := base64.StdEncoding.EncodeToString([]byte(data))
	return encodedData
}

func (i *ItemRepository) WaitOrder() {
	i.WG.Wait()
}
