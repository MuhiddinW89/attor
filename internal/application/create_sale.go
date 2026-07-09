package application

import (
	"github.com/MuhiddinW89/attor/internal/clients"
	"github.com/MuhiddinW89/attor/internal/reminders"
	"github.com/MuhiddinW89/attor/internal/sales"
)

type CreateSaleUseCase struct {
	clientsService clients.Service
	salesService sales.Service
	remindersService reminders.Service
}

func NewCreateSaleUseCase (
	clientsService clients.Service,
	salesService sales.Service,
	remindersService reminders.Service,
) *CreateSaleUseCase {
	return &CreateSaleUseCase{
		clientsService: clientsService,
		salesService:  salesService,
		remindersService: remindersService,
	}
}