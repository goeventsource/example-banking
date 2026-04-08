package banking

import (
	"encoding/json"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
)

type accountJSON struct {
	ID          uuid.UUID `json:"id"`
	Version     uint      `json:"version"`
	Balance     uint      `json:"balance"`
	Currency    string    `json:"currency"`
	OpenedAt    time.Time `json:"opened_at"`
	ActivatedAt time.Time `json:"activated_at"`
	ActivatedBy uuid.UUID `json:"activated_by"`
	ClosedAt    time.Time `json:"closed_at"`
	ClosedBy    uuid.UUID `json:"closed_by"`
}

type RootEncodeDecoder struct {
	factoryFn goeventsource.FactoryFunc[uuid.UUID, *Account]
}

func NewRootEncodeDecoder(factoryFn goeventsource.FactoryFunc[uuid.UUID, *Account]) *RootEncodeDecoder {
	return &RootEncodeDecoder{
		factoryFn: factoryFn,
	}
}

func (r *RootEncodeDecoder) Encode(a *Account) ([]byte, error) {
	return json.Marshal(accountJSON{
		ID:          goeventsource.RootID(a),
		Version:     uint(goeventsource.RootVersion(a)),
		Balance:     uint(a.balance),
		Currency:    a.currency.Code,
		OpenedAt:    a.openedAt,
		ActivatedAt: a.activatedAt,
		ActivatedBy: a.activatedBy,
		ClosedAt:    a.closedAt,
		ClosedBy:    a.closedBy,
	})
}

func (r *RootEncodeDecoder) Decode(data []byte, a **Account) error {
	var accountJSON accountJSON
	if err := json.Unmarshal(data, &accountJSON); err != nil {
		return err
	}

	acc := r.factoryFn(accountJSON.ID, goeventsource.Version(accountJSON.Version))
	acc.balance = money.Amount(accountJSON.Balance)
	acc.currency = *money.GetCurrency(accountJSON.Currency)
	acc.openedAt = accountJSON.OpenedAt
	acc.activatedAt = accountJSON.ActivatedAt
	acc.activatedBy = accountJSON.ActivatedBy
	acc.closedAt = accountJSON.ClosedAt
	acc.closedBy = accountJSON.ClosedBy

	*a = acc
	return nil
}
