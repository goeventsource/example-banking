package banking

import (
	"errors"
	"fmt"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
)

const AccountAggregateName = "banking_account"

var (
	ErrAccountNotActive   = errors.New("not active")
	ErrAccountNotInReview = errors.New("not in review")

	ErrAccountClose    = errors.New("could not close")
	ErrAccountActivate = errors.New("could not activate")
	ErrAccountDeposit  = errors.New("could not deposit")
	ErrAccountWithdraw = errors.New("could not withdraw")
)

type Account struct {
	*goeventsource.BaseRoot[uuid.UUID]
	balance     money.Amount
	currency    money.Currency
	openedAt    time.Time
	activatedAt time.Time
	activatedBy uuid.UUID
	closedAt    time.Time
	closedBy    uuid.UUID
}

func OpenAccount(currency money.Currency) *Account {
	id := uuid.New()
	a := &Account{
		BaseRoot: goeventsource.NewBase(id, AccountAggregateName, 0),
	}

	goeventsource.Record(a, AccountWasOpenedV1{ID: id, Currency: currency, OpenedAt: time.Now()})
	return a
}

func (a *Account) Activate(agentID uuid.UUID) error {
	if a.Status() != InReview {
		return fmt.Errorf("%w: %w", ErrAccountActivate, ErrAccountNotInReview)
	}

	goeventsource.Record(a, AccountWasActivatedV1{ActivatedAt: time.Now(), ActivatedBy: agentID})
	return nil
}

func (a *Account) Deposit(amount money.Amount, currency money.Currency) error {
	if a.Status() != Active {
		return fmt.Errorf("%w: %w", ErrAccountDeposit, ErrAccountNotActive)
	}

	if a.currency != currency {
		return fmt.Errorf("%w: %w", ErrAccountDeposit, money.ErrCurrencyMismatch)
	}

	goeventsource.Record(a, AccountBalanceWasDepositedV1{Amount: amount})
	return nil
}

func (a *Account) DepositWithExchange(
	fromAmount money.Amount,
	fromCurrency money.Currency,
	name ExchangerName,
	rate ExchangerRate,
	toAmount money.Amount,
	toCurrency money.Currency,
) error {
	if a.Status() != Active {
		return fmt.Errorf("%w: %w", ErrAccountDeposit, ErrAccountNotActive)
	}

	if a.currency != toCurrency {
		return fmt.Errorf("%w: %w", ErrAccountDeposit, money.ErrCurrencyMismatch)
	}

	ev := AccountBalanceCurrencyExchangeWasDepositedV1{
		FromAmount:   fromAmount,
		FromCurrency: fromCurrency,
		Name:         name,
		Rate:         rate,
		ToAmount:     toAmount,
		ToCurrency:   toCurrency,
	}

	goeventsource.Record(a, ev)

	return nil
}

func (a *Account) Withdraw(amount money.Amount, currency money.Currency) error {
	if a.Status() != Active {
		return fmt.Errorf("%w: %w", ErrAccountWithdraw, ErrAccountNotActive)
	}

	if a.currency != currency {
		return fmt.Errorf("%w: %w", ErrAccountWithdraw, money.ErrCurrencyMismatch)
	}

	if a.balance < amount {
		return fmt.Errorf("%w: not enough balance", ErrAccountWithdraw)
	}

	goeventsource.Record(a, AccountBalanceWasWithdrewV1{Amount: amount})

	return nil
}

func (a *Account) WithdrawWithExchange(
	fromAmount money.Amount,
	fromCurrency money.Currency,
	name ExchangerName,
	rate ExchangerRate,
	toAmount money.Amount,
	toCurrency money.Currency,
) error {
	if a.Status() != Active {
		return fmt.Errorf("%w: %w", ErrAccountWithdraw, ErrAccountNotActive)
	}

	if a.currency != toCurrency {
		return fmt.Errorf("%w: %w", ErrAccountWithdraw, money.ErrCurrencyMismatch)
	}

	if a.balance < toAmount {
		return fmt.Errorf("%w: not enough balance", ErrAccountWithdraw)
	}

	goeventsource.Record(a, AccountBalanceCurrencyExchangeWasWithdrewV1{
		FromAmount:   fromAmount,
		FromCurrency: fromCurrency,
		Name:         name,
		Rate:         rate,
		ToAmount:     toAmount,
		ToCurrency:   toCurrency,
	})

	return nil
}

func (a *Account) Close(agentID uuid.UUID) error {
	if a.Status() != Active {
		return fmt.Errorf("%w: %w", ErrAccountClose, ErrAccountNotActive)
	}

	goeventsource.Record(a, AccountWasClosedV1{ClosedAt: time.Now(), ClosedBy: agentID})
	return nil
}

func (a *Account) ID() uuid.UUID {
	return goeventsource.RootID(a)
}

func (a *Account) Balance() money.Amount {
	return a.balance
}

func (a *Account) Currency() money.Currency {
	return a.currency
}

func (a *Account) OpenedAt() time.Time {
	return a.openedAt
}

func (a *Account) ActivatedAt() time.Time {
	return a.activatedAt
}

func (a *Account) ClosedAt() time.Time {
	return a.closedAt
}

func (a *Account) Status() Status {
	switch {
	case !a.closedAt.IsZero():
		return Closed
	case !a.activatedAt.IsZero():
		return Active
	}
	return InReview
}

func (a *Account) Apply(ev goeventsource.DomainEvent) {
	switch e := ev.(type) {
	case AccountWasOpenedV1:
		a.currency = e.Currency
		a.openedAt = e.OpenedAt
	case AccountWasActivatedV1:
		a.activatedAt = e.ActivatedAt
		a.activatedBy = e.ActivatedBy
	case AccountBalanceWasDepositedV1:
		a.balance += e.Amount
	case AccountBalanceCurrencyExchangeWasDepositedV1:
		a.balance += e.ToAmount
	case AccountBalanceWasWithdrewV1:
		a.balance -= e.Amount
	case AccountBalanceCurrencyExchangeWasWithdrewV1:
		a.balance -= e.ToAmount
	case AccountWasClosedV1:
		a.closedAt = e.ClosedAt
		a.closedBy = e.ClosedBy
	}
}
