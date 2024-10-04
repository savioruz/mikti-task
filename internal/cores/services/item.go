package services

import (
	"github.com/savioruz/mikti-task/week-2/internal/cores/entities"
	"github.com/savioruz/mikti-task/week-2/internal/cores/ports"
)

type ItemService struct {
	ItemRepository ports.ItemRepository
}

func NewItemService(itemRepository ports.ItemRepository) *ItemService {
	return &ItemService{
		ItemRepository: itemRepository,
	}
}

func (i *ItemService) AddItem(name string, price float64) *[]entities.Item {
	return i.ItemRepository.AddItem(name, price)
}

func (i *ItemService) GetItems() *[]entities.Item {
	return i.ItemRepository.GetItems()
}

func (i *ItemService) TakeOrder(restaurant *entities.Restaurant, ch chan<- entities.Order) {
	i.ItemRepository.TakeOrder(restaurant, ch)
}

func (i *ItemService) ValidateOrder(restaurant *entities.Restaurant, item string) (*entities.Item, bool) {
	return i.ItemRepository.ValidateOrder(restaurant, item)
}

func (i *ItemService) ValidatePrice(price string) (float64, error) {
	return i.ItemRepository.ValidatePrice(price)
}

func (i *ItemService) EncodeOrder(order entities.Order) string {
	return i.ItemRepository.EncodeOrder(order)
}

func (i *ItemService) WaitOrder() {
	i.ItemRepository.WaitOrder()
}
