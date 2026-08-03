package clients

import (
	"context"
	"errors"
	"time"

	"github.com/MuhiddinW89/attor/internal/clients"
	"github.com/MuhiddinW89/attor/internal/reminders"
	"github.com/MuhiddinW89/attor/internal/sales"
	"github.com/google/uuid"
)

type ClientHistoryUseCase struct {
	clientService   clients.Service
	saleService     sales.Service
	reminderService reminders.Service
}

func NewClientHistoryUseCase(
	clientService clients.Service,
	saleService sales.Service,
	reminderService reminders.Service,
) *ClientHistoryUseCase {
	return &ClientHistoryUseCase{
		clientService:   clientService,
		saleService:     saleService,
		reminderService: reminderService,
	}
}

type ClientHistory struct {
	Client      ClientInfo    `json:"client"`
	TotalSales  int           `json:"total_sales"`
	TotalAmount float64       `json:"total_amount"`
	Sales       []SaleInfo    `json:"sales"`
	Reminder    *ReminderInfo `json:"reminder,omitempty"`
}

type ClientInfo struct {
	ID        string  `json:"id"`
	FullName  string  `json:"full_name"`
	Phone     string  `json:"phone"`
	Instagram *string `json:"instagram,omitempty"`
}

type SaleInfo struct {
	ID          string  `json:"id"`
	PerfumeName string  `json:"perfume_name"`
	VolumeML    int     `json:"volume_ml"`
	Price       float64 `json:"price"`
	Comment     *string `json:"comment"`
	SaleDate    string  `json:"sale_date"`
}

type ReminderInfo struct {
	ID         string `json:"id"`
	ReminderAt string `json:"reminder_at"`
	IsSent     bool   `json:"is_sent"`
}

func (uc *ClientHistoryUseCase) Execute(
	ctx context.Context,
	clientID uuid.UUID,
) (*ClientHistory, error) {

	client, err := uc.clientService.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	salesList, err := uc.saleService.ListByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	reminder, err := uc.reminderService.GetNearestByClientID(ctx, clientID)
	if err != nil && !errors.Is(err, reminders.ErrReminderNotFound) {
		return nil, err
	}

	result := &ClientHistory{
		Client: ClientInfo{
			ID:        client.ID.String(),
			FullName:  client.FullName,
			Phone:     client.Phone,
			Instagram: client.Instagram,
		},
		Sales: make([]SaleInfo, 0, len(salesList)),
	}

	result.TotalSales = len(salesList)

	var totalAmount float64

	for _, sale := range salesList {

		totalAmount += sale.Price

		result.Sales = append(
			result.Sales,
			SaleInfo{
				ID:          sale.ID.String(),
				PerfumeName: sale.PerfumeName,
				VolumeML:    sale.VolumeML,
				Price:       sale.Price,
				Comment:     sale.Comment,
				SaleDate:    sale.SaleDate.Format(time.RFC3339),
			},
		)
	}

	result.TotalAmount = totalAmount

	if reminder != nil {
		result.Reminder = &ReminderInfo{
			ID:         reminder.ID.String(),
			ReminderAt: reminder.ReminderAt.Format(time.RFC3339),
			IsSent:     reminder.IsSent,
		}
	}

	return result, nil
}
