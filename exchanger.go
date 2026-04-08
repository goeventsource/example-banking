package banking

import (
	"context"

	"github.com/Rhymond/go-money"
)

type ExchangerName string

type ExchangerRate uint

type Exchanger interface {
	Exchange(ctx context.Context, amount money.Amount, from money.Currency, to money.Currency) (ExchangerName, ExchangerRate, money.Amount, error)
}
