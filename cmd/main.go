package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	banking "github.com/goeventsource/example-banking"
	binInternal "github.com/goeventsource/example-banking/cmd/internal"
	"github.com/goeventsource/example-banking/internal"
	"github.com/goeventsource/goeventsource"
	"github.com/goeventsource/pgx"
	"github.com/goeventsource/pgx/pgxtest"
)

var factoryFn = func(id uuid.UUID, ver goeventsource.Version) *banking.Account {
	return &banking.Account{
		BaseRoot: goeventsource.NewBase(id, banking.AccountAggregateName, ver),
	}
}

func main() {
	ctx := context.Background()
	db, cleanup, err := pgxtest.NewDemoPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	cfg := pgxtest.NewRepositoryConfig(db, factoryFn)

	cfg.StoreConfig.Codec = goeventsource.NewDomainEventEncodeDecoderWrapper(map[goeventsource.DomainEventName]goeventsource.DomainEventEncodeDecoder{
		banking.AccountWasOpenedV1{}.DomainEventName():                           goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountWasOpenedV1](),
		banking.AccountWasActivatedV1{}.DomainEventName():                        goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountWasActivatedV1](),
		banking.AccountBalanceWasDepositedV1{}.DomainEventName():                 goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountBalanceWasDepositedV1](),
		banking.AccountBalanceCurrencyExchangeWasDepositedV1{}.DomainEventName(): goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountBalanceCurrencyExchangeWasDepositedV1](),
		banking.AccountBalanceWasWithdrewV1{}.DomainEventName():                  goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountBalanceWasWithdrewV1](),
		banking.AccountBalanceCurrencyExchangeWasWithdrewV1{}.DomainEventName():  goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountBalanceCurrencyExchangeWasWithdrewV1](),
		banking.AccountWasClosedV1{}.DomainEventName():                           goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountWasClosedV1](),
	})

	codec := banking.NewRootEncodeDecoder(factoryFn)
	snapCfg := pgxtest.NewSnapshotterConfig(
		db,
		codec,
		goeventsource.SnapshotterWriteStrategyByVersionStep[uuid.UUID, *banking.Account](goeventsource.Version(3)),
	)
	snap, err := pgxtest.NewSnapshotter(snapCfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.Opts = append(cfg.Opts, pgx.WithSnapshotterOpt(snap))

	repo, _, err := pgxtest.NewRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}
	svc := internal.NewService(repo, nil)

	srv := binInternal.New(svc, codec)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	if err := http.ListenAndServe(addr, srv); err != http.ErrServerClosed {
		panic(err)
	}
}
