package main

import (
	"atlas-cashshop/cashshop"
	character2 "atlas-cashshop/cashshop/character"
	wishlist2 "atlas-cashshop/cashshop/character/wishlist"
	"atlas-cashshop/cashshop/item"
	"atlas-cashshop/database"
	"atlas-cashshop/kafka"
	"atlas-cashshop/logger"
	"atlas-cashshop/rest"
	"atlas-cashshop/tracing"
	"atlas-cashshop/wz"
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const serviceName = "atlas-cashshop"
const consumerGroupId = "Cash Shop Orchestration Service"

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}
	defer func(tc io.Closer) {
		err = tc.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close tracer.")
		}
	}(tc)

	wzDir := os.Getenv("WZ_DIR")
	wz.GetFileCache().Init(wzDir)

	err = item.GetCache().Init()
	if err != nil {
		l.WithError(err).Errorf("Unable to load quest cache.")
	}

	db := database.Connect(l, database.SetMigrations(character2.Migration, wishlist2.Migration))

	kafka.CreateConsumers(l, ctx, wg,
		cashshop.EnterCashShopCommandConsumer()(consumerGroupId),
		cashshop.PurchaseCashShopItemCommandConsumer(db)(consumerGroupId),
		character2.CreatedConsumer(db)(consumerGroupId),
		character2.AwardCreditConsumer(db)(consumerGroupId),
		character2.AwardPointsConsumer(db)(consumerGroupId),
		character2.AwardPrepaidConsumer(db)(consumerGroupId),
		wishlist2.ModifyWishlistConsumer(db)(consumerGroupId))

	rest.CreateService(l, db, ctx, wg, "/ms/cashshop", character2.InitResource, wishlist2.InitResource)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	l.Infoln("Service shutdown.")
}
