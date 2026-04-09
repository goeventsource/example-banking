package banking

import (
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
)

type AccountWasOpenedV1 struct {
	ID       uuid.UUID
	Currency money.Currency
	OpenedAt time.Time
}

func (e AccountWasOpenedV1) DomainEventName() goeventsource.DomainEventName {
	return "banking_account_was_opened_v1"
}

type AccountWasActivatedV1 struct {
	ActivatedBy uuid.UUID
	ActivatedAt time.Time
}

func (e AccountWasActivatedV1) DomainEventName() goeventsource.DomainEventName {
	return "banking_account_was_activated_v1"
}

type AccountBalanceWasDepositedV1 struct {
	Amount money.Amount
}

func (e AccountBalanceWasDepositedV1) DomainEventName() goeventsource.DomainEventName {
	return "banking_account_balance_was_deposited_v1"
}

type AccountBalanceCurrencyExchangeWasDepositedV1 struct {
	FromAmount   money.Amount
	FromCurrency money.Currency
	Name         ExchangerName
	Rate         ExchangerRate
	ToAmount     money.Amount
	ToCurrency   money.Currency
}

func (e AccountBalanceCurrencyExchangeWasDepositedV1) DomainEventName() goeventsource.DomainEventName {
	return "banking_account_balance_currency_exchange_was_deposited_v1"
}

type AccountBalanceWasWithdrawnV1 struct {
	Amount money.Amount
}

func (e AccountBalanceWasWithdrawnV1) DomainEventName() goeventsource.DomainEventName {
	return "banking_account_balance_was_withdrawn_v1"
}

type AccountBalanceCurrencyExchangeWasWithdrawnV1 struct {
	FromAmount   money.Amount
	FromCurrency money.Currency
	Name         ExchangerName
	Rate         ExchangerRate
	ToAmount     money.Amount
	ToCurrency   money.Currency
}

func (e AccountBalanceCurrencyExchangeWasWithdrawnV1) DomainEventName() goeventsource.DomainEventName {
	return "banking_account_balance_currency_exchange_was_withdrawn_v1"
}

type AccountWasClosedV1 struct {
	ClosedAt time.Time
	ClosedBy uuid.UUID
}

func (e AccountWasClosedV1) DomainEventName() goeventsource.DomainEventName {
	return "banking_account_was_closed_v1"
}
