package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/ui/moneyinput"
)

type parseHandlers struct{}

func newParseHandlers() *parseHandlers {
	return &parseHandlers{}
}

func (p *parseHandlers) register(r chi.Router) {
	r.Get("/api/parse/decimal", moneyinput.DecimalHandler())
	r.Get("/api/parse/money", moneyinput.MoneyHandler())
	r.Get("/api/parse/money-restricted", moneyinput.MoneyHandler("USD", "EUR"))
}
