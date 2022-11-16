package handlersImpl

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/matster07/user-balance-service/internal/app/entity/accounts"
	"github.com/matster07/user-balance-service/internal/app/entity/categories"
	"github.com/matster07/user-balance-service/internal/app/entity/orders"
	"github.com/matster07/user-balance-service/internal/app/entity/transactions"
)

func (h *handler) Process(deliverStatus accounts.DeliverStatusDto) error {
	tx, _ := h.client.Begin(context.TODO())

	defer func(tx pgx.Tx) {
		err := tx.Rollback(context.TODO())
		if err != nil {
			h.logger.Trace("Rollback transaction")
		}
	}(tx)

	order, err := h.orderRepository.UpdateStatus(tx, deliverStatus.OrderId, deliverStatus.Status)
	if err != nil {
		return errors.New(fmt.Sprintf("order %d wasn't fonud", order.ID))
	} else if order.Status != "IN_PROGRESS" {
		return errors.New(fmt.Sprintf("order %d was already processed", order.ID))
	}

	var category categories.Category
	var transactionType string
	var receiverAccount accounts.Account

	switch deliverStatus.Status {
	case "COMPLETED":
		{
			transactionType = "PROFIT"

			category, err = h.categoryRepository.FindById(order.CategoryId)
			if err != nil {
				return errors.New(fmt.Sprintf("category %d wasn't found", order.CategoryId))
			}

			receiverAccount, err = h.accountRepository.FindByType("PROFIT_ACCOUNT")
			if err != nil {
				return errors.New(fmt.Sprintf("failed to find company receiverAccount"))
			}
		}
	default:
		{
			transactionType = "REFUND"

			receiverAccId := order.UserAccountId
			receiverAccount, err = h.accountRepository.FindById(receiverAccId)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to get company receiverAccount with id %d", receiverAccId))
			}
		}
	}

	reservingAccountId := category.AccountId
	reservingAcc, err := h.accountRepository.FindById(reservingAccountId)
	if err != nil {
		return errors.New(fmt.Sprintf("reserving receiverAccount %d wasn't found", reservingAccountId))
	}

	err = h.accountRepository.Update(tx, accounts.Account{
		ID:      reservingAcc.ID,
		Balance: reservingAcc.Balance - deliverStatus.Amount,
	})
	if err != nil {
		return errors.New("failed companyAcc update receiverAccount")
	}

	err = h.accountRepository.Update(tx, accounts.Account{
		ID:      receiverAccount.ID,
		Balance: receiverAccount.Balance + deliverStatus.Amount,
	})
	if err != nil {
		return errors.New("failed companyAcc update receiverAccount")
	}

	err = h.transactionRepository.Create(tx, transactions.Transaction{
		Amount:          order.Price,
		From:            reservingAcc.ID,
		To:              receiverAccount.ID,
		TransactionType: transactionType,
		Comment:         fmt.Sprintf("order_id: %d", order.ID),
	})
	if err != nil {
		return errors.New("failed companyAcc save transaction")
	}

	if err = tx.Commit(context.TODO()); err != nil {
		return errors.New("failed companyAcc commit transaction")
	}

	return nil
}

func (h *handler) determineReceiver(order orders.Order) (account accounts.Account, err error) {
	var accountId uint

	switch order.Status {
	case "PROFIT":
		{
			category, err := h.categoryRepository.FindById(order.CategoryId)
			if err != nil {
				return account, errors.New(fmt.Sprintf("category %d wasn't found", order.CategoryId))
			}
			accountId = category.AccountId
		}
	default:
		{
			accountId = order.UserAccountId
		}
	}

	account, err = h.accountRepository.FindById(accountId)
	if err != nil {
		return account, err
	}

	return account, err
}
