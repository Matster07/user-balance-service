package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/matster07/user-balance-service/internal/app/data/dto"
	"github.com/matster07/user-balance-service/internal/app/data/entity"
	"github.com/matster07/user-balance-service/internal/pkg/client/postgresql"
)

// Process Обработка результата оказания услуги
func (h *Handler) Process(deliverStatus dto.DeliverStatusDto) error {
	order, err := h.Order.FindById(deliverStatus.OrderId)
	if err != nil {
		return errors.New(fmt.Sprintf("order %d wasn't fonud", order.ID))
	} else if order.Status != "IN_PROGRESS" {
		return errors.New(fmt.Sprintf("order %d was already processed", order.ID))
	}

	tx, err := h.Begin(context.TODO())
	if err != nil {
		return errors.New("error while starting transaction")
	}

	defer postgresql.RollbackTx(tx)

	err = h.Order.UpdateStatus(tx, deliverStatus.OrderId, deliverStatus.Status)
	if err != nil {
		return err
	}

	var transactionType string
	var receiverAccount entity.Account

	service, err := h.Service.FindById(order.ServiceId)
	if err != nil {
		return errors.New(fmt.Sprintf("service %d wasn't found", order.ServiceId))
	}

	switch deliverStatus.Status {
	case "COMPLETED":
		{
			transactionType = "PROFIT"

			receiverAccount, err = h.Account.FindByType("PROFIT_ACCOUNT")
			if err != nil {
				return errors.New(fmt.Sprintf("failed to find company receiverAccount"))
			}
		}
	default:
		{
			transactionType = "REFUND"

			receiverAccId := order.UserAccountId
			receiverAccount, err = h.Account.FindById(receiverAccId)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to get company receiverAccount with id %d", receiverAccId))
			}
		}
	}

	reservingAccountId := service.AccountId
	reservingAcc, err := h.Account.FindById(reservingAccountId)
	if err != nil {
		return errors.New(fmt.Sprintf("reserving receiverAccount %d wasn't found", reservingAccountId))
	}

	err = h.Account.Update(tx, entity.Account{
		ID:      reservingAcc.ID,
		Balance: reservingAcc.Balance - deliverStatus.Amount,
	})
	if err != nil {
		return errors.New("failed companyAcc update receiverAccount")
	}

	err = h.Account.Update(tx, entity.Account{
		ID:      receiverAccount.ID,
		Balance: receiverAccount.Balance + deliverStatus.Amount,
	})
	if err != nil {
		return errors.New("failed companyAcc update receiverAccount")
	}

	err = h.Transaction.Create(tx, entity.Transaction{
		Amount:          order.Price,
		From:            reservingAcc.ID,
		To:              receiverAccount.ID,
		TransactionType: transactionType,
		Comment:         fmt.Sprintf("order_id: %d", order.ID),
	})
	if err != nil {
		return errors.New("failed companyAcc save transaction")
	}

	postgresql.CommitTx(tx)

	return nil
}

func (h *Handler) determineReceiver(order entity.Order) (account entity.Account, err error) {
	var accountId uint

	switch order.Status {
	case "PROFIT":
		{
			service, err := h.Service.FindById(order.ServiceId)
			if err != nil {
				return account, errors.New(fmt.Sprintf("service %d wasn't found", order.ServiceId))
			}
			accountId = service.AccountId
		}
	default:
		{
			accountId = order.UserAccountId
		}
	}

	account, err = h.Account.FindById(accountId)
	if err != nil {
		return account, err
	}

	return account, err
}
