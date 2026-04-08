package banking

import (
	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
)

type Repo = goeventsource.Repository[uuid.UUID, *Account]
