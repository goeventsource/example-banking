package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
	"github.com/goeventsource/pgx"
	"github.com/goeventsource/pgx/pgxtest"

	banking "github.com/goeventsource/example-banking"
	binInternal "github.com/goeventsource/example-banking/cmd/internal"
	"github.com/goeventsource/example-banking/internal"
)

var factoryFn = func(id uuid.UUID, ver goeventsource.Version) *banking.Account {
	return &banking.Account{
		BaseRoot: goeventsource.NewBase(id, banking.AccountAggregateName, ver),
	}
}

func domainEventCodecs() map[goeventsource.DomainEventName]goeventsource.DomainEventEncodeDecoder {
	m := make(map[goeventsource.DomainEventName]goeventsource.DomainEventEncodeDecoder)
	m[banking.AccountWasOpenedV1{}.DomainEventName()] =
		goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountWasOpenedV1]()
	m[banking.AccountWasActivatedV1{}.DomainEventName()] =
		goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountWasActivatedV1]()
	m[banking.AccountBalanceWasDepositedV1{}.DomainEventName()] =
		goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountBalanceWasDepositedV1]()
	m[banking.AccountBalanceCurrencyExchangeWasDepositedV1{}.DomainEventName()] =
		goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountBalanceCurrencyExchangeWasDepositedV1]()
	m[banking.AccountBalanceWasWithdrawnV1{}.DomainEventName()] =
		goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountBalanceWasWithdrawnV1]()
	m[banking.AccountBalanceCurrencyExchangeWasWithdrawnV1{}.DomainEventName()] =
		goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountBalanceCurrencyExchangeWasWithdrawnV1]()
	m[banking.AccountWasClosedV1{}.DomainEventName()] =
		goeventsource.NewJSONDomainEventEncodeDecoder[banking.AccountWasClosedV1]()
	return m
}

func run(ctx context.Context) error {
	db, cleanup, err := pgxtest.NewDemoPool(ctx)
	if err != nil {
		return fmt.Errorf("demo pool: %w", err)
	}
	defer cleanup()

	cfg := pgxtest.NewRepositoryConfig(db, factoryFn)
	cfg.Codec = goeventsource.NewDomainEventEncodeDecoderWrapper(domainEventCodecs())

	codec := banking.NewRootEncodeDecoder(factoryFn)
	const snapshotEveryNVersions uint64 = 3
	snapCfg := pgxtest.NewSnapshotterConfig(
		db,
		codec,
		goeventsource.SnapshotterWriteStrategyByVersionStep[uuid.UUID, *banking.Account](goeventsource.Version(snapshotEveryNVersions)),
	)
	snap, err := pgxtest.NewSnapshotter(snapCfg)
	if err != nil {
		return fmt.Errorf("snapshotter: %w", err)
	}
	cfg.Opts = append(cfg.Opts, pgx.WithSnapshotterOpt(snap))

	repo, _, err := pgxtest.NewRepository(cfg)
	if err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	svc := internal.NewService(repo, nil)

	srv := binInternal.New(svc, codec)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	if err := http.ListenAndServe(addr, srv); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
