package usecase

import (
	"github.com/eneridangelis/golangExpert/proj-3/internal/entity"
	"github.com/eneridangelis/golangExpert/proj-3/pkg/events"
)

type Order struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrderOutputDTO struct {
	Orders []Order `json:"orders"`
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *ListOrdersUseCase) Execute() (ListOrderOutputDTO, error) {
	orders, err := c.OrderRepository.ListAll()
	if err != nil {
		return ListOrderOutputDTO{}, err
	}

	var ordersDTO []Order
	for _, order := range orders {
		ordersDTO = append(ordersDTO, Order{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		})
	}

	dto := ListOrderOutputDTO{
		Orders: ordersDTO,
	}

	//[ene] checar dispach de eventos
	// c.OrderCreated.SetPayload(dto)
	// c.EventDispatcher.Dispatch(c.OrderCreated)

	return dto, nil
}
