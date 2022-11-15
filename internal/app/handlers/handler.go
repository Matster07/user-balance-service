package handlers

import (
	"github.com/gorilla/mux"
	"github.com/matster07/user-balance-service/internal/app/entity/accounts"
)

type Handler interface {
	Register(router *mux.Router)
	Process(deliverStatus accounts.DeliverStatusDto) error
}
