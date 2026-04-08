package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"

	banking "github.com/goeventsource/example-banking"
)

var (
	ErrServiceOpenAccount     = errors.New("could not open account")
	ErrServiceActivate        = errors.New("could not activate account")
	ErrServiceWithdraw        = errors.New("could not withdraw account")
	ErrServiceDeposit         = errors.New("could not deposit account")
	ErrServiceClose           = errors.New("could not close account")
	ErrExchangerNotConfigured = errors.New("currency exchange is not configured")
)

type Service struct {
	Repo      banking.Repo
	Exchanger banking.Exchanger
}

func NewService(repo banking.Repo, exchanger banking.Exchanger) *Service {
	return &Service{
		Repo:      repo,
		Exchanger: exchanger,
	}
}

func (s *Service) OpenAccount(ctx context.Context, c money.Currency) (*banking.Account, error) {
	acc := banking.OpenAccount(c)
	if err := s.Repo.Write(ctx, acc); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceOpenAccount, err)
	}

	return acc, nil
}

func (s *Service) Activate(ctx context.Context, accountID, agentID uuid.UUID) (*banking.Account, error) {
	acc, err := s.Repo.Read(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceActivate, err)
	}

	if err := acc.Activate(agentID); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceActivate, err)
	}

	if err := s.Repo.Write(ctx, acc); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceActivate, err)
	}

	return acc, nil
}

func (s *Service) Withdraw(
	ctx context.Context,
	accountID uuid.UUID,
	amount money.Amount,
	currency money.Currency,
) (*banking.Account, error) {
	acc, err := s.Repo.Read(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceWithdraw, err)
	}

	switch {
	case acc.Currency() == currency:
		if err := acc.Withdraw(amount, currency); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrServiceWithdraw, err)
		}
	default:
		if s.Exchanger == nil {
			return nil, fmt.Errorf("%w: %w", ErrServiceWithdraw, ErrExchangerNotConfigured)
		}
		name, rate, toAmount, err := s.Exchanger.Exchange(ctx, amount, currency, acc.Currency())
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrServiceWithdraw, err)
		}

		if err := acc.WithdrawWithExchange(
			amount,
			currency,
			name,
			rate,
			toAmount,
			acc.Currency(),
		); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrServiceWithdraw, err)
		}
	}

	if err := s.Repo.Write(ctx, acc); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceWithdraw, err)
	}

	return acc, nil
}

func (s *Service) Deposit(
	ctx context.Context,
	accountID uuid.UUID,
	amount money.Amount,
	currency money.Currency,
) (*banking.Account, error) {
	acc, err := s.Repo.Read(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceDeposit, err)
	}

	switch {
	case acc.Currency() == currency:
		if err := acc.Deposit(amount, currency); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrServiceDeposit, err)
		}
	default:
		if s.Exchanger == nil {
			return nil, fmt.Errorf("%w: %w", ErrServiceDeposit, ErrExchangerNotConfigured)
		}
		name, rate, toAmount, err := s.Exchanger.Exchange(ctx, amount, currency, acc.Currency())
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrServiceDeposit, err)
		}

		if err := acc.DepositWithExchange(
			amount,
			currency,
			name,
			rate,
			toAmount,
			acc.Currency(),
		); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrServiceDeposit, err)
		}
	}

	if err := s.Repo.Write(ctx, acc); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceDeposit, err)
	}

	return acc, nil
}

func (s *Service) Close(ctx context.Context, accountID, agentID uuid.UUID) (*banking.Account, error) {
	acc, err := s.Repo.Read(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceClose, err)
	}

	if err := acc.Close(agentID); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceClose, err)
	}

	if err := s.Repo.Write(ctx, acc); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceClose, err)
	}

	return acc, nil
}
