package application

import (
	"context"
	"errors"
	"time"

	"github.com/MuhiddinW89/attor/internal/clients"
	"github.com/MuhiddinW89/attor/internal/reminders"
	"github.com/MuhiddinW89/attor/internal/sales"
)

type CreateSaleUseCase struct {
	clientService   clients.Service
	saleService     sales.Service
	reminderService reminders.Service
}

func NewCreateSaleUseCase(
	clientService clients.Service,
	saleService sales.Service,
	reminderService reminders.Service,
) *CreateSaleUseCase {
	return &CreateSaleUseCase{
		clientService:   clientService,
		saleService:     saleService,
		reminderService: reminderService,
	}
}

type CreateSaleRequest struct {
	Phone       string
	FullName    string
	PerfumeName string
	VolumeML    int
	Price       float64
	Comment     *string
	SaleDate    time.Time
}

func (uc *CreateSaleUseCase) Execute(
	ctx context.Context,
	req CreateSaleRequest,
) (*sales.Sale, error) {

	client, err := uc.clientService.GetByPhone(ctx, req.Phone)
	if err != nil {
		if errors.Is(err, clients.ErrClientNotFound) {
			client, err = uc.clientService.Create(
				ctx,
				clients.CreateClientInput{
					FullName: req.FullName,
					Phone:    req.Phone,
				},
			)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	sale, err := uc.saleService.Create(
		ctx,
		sales.CreateSaleInput{
			ClientID:    client.ID,
			PerfumeName: req.PerfumeName,
			VolumeML:    req.VolumeML,
			Price:       req.Price,
			Comment:     req.Comment,
			SaleDate:    req.SaleDate,
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = uc.reminderService.Create(
		ctx,
		reminders.CreateReminderInput{
			ClientID:   client.ID,
			SaleID:     sale.ID,
			ReminderAt: sale.SaleDate.AddDate(0, 0, reminders.DefaultReminderAfterDays),
		},
	)
	if err != nil {
		return nil, err
	}

	return sale, nil
}
