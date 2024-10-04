package ports

import "github.com/savioruz/mikti-task/week-2/internal/cores/entities"

type ItemRepository interface {
	AddItem(name string, price float64) *[]entities.Item
	GetItems() *[]entities.Item
	TakeOrder(restaurant *entities.Restaurant, ch chan<- entities.Order)
	ValidateOrder(restaurant *entities.Restaurant, item string) (*entities.Item, bool)
	ValidatePrice(price string) (float64, error)
	EncodeOrder(order entities.Order) string
	WaitOrder()
}
