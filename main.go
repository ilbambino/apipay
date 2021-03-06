package main

import (
	"apipay/config"
	"apipay/persistent"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	gitHash string // This is set when building
)

func getHandler(logger *zap.Logger, paymentDb persistent.Payments) http.Handler {

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	paymentsRoute := router.Group("/payments/") //TODO add any specific AUTH
	// TODO add tracing and request id to the logger
	{
		paymentsRoute.GET("/", getPayments(logger, paymentDb))

		paymentsRoute.GET("/:paymentID", getOnePayment(logger, paymentDb))

		paymentsRoute.PUT("/:paymentID", updatePayment(logger, paymentDb))

		paymentsRoute.DELETE("/:paymentID", deletePayment(logger, paymentDb))

		paymentsRoute.POST("/", createPayment(logger, paymentDb))
	}

	return router
}

// @title APIPAY Payments API
// @version 1.0
// @description This is an example implementation of an API to serve Payments
// @description The documentation is automatically built from the code.
// @description There are many TODOs and things not fully finalized.
// @description Please read the README file in the same repository.

func main() {
	ctx := context.Background()
	err := config.Load()
	if err != nil {
		panic("Cannot load config " + err.Error())
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic("Cannot create logger")
	}

	//nolint
	defer logger.Sync() // flushes buffer, if any

	logger.Sugar().Infow("server-init", "gitHash", gitHash)

	db, err := persistent.Connect(ctx,
		viper.GetString(config.MongoHost),
		viper.GetInt(config.MongoPort),
		viper.GetString(config.MongoUser),
		viper.GetString(config.MongoPassword))

	if err != nil {
		logger.Sugar().Fatalw("init-db-error", "error", err)
		panic("init-error")
	}
	defer db.Close(ctx)

	paymentsDB, err := persistent.GetPayments(ctx, db)
	if err != nil {
		logger.Sugar().Fatalw("init-db-payments-error", "error", err)
		panic("init-error")
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: getHandler(logger, paymentsDB),
	}

	// Open the server for incoming connections
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Sugar().Fatalw("init-error", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 3)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Sugar().Fatalw("shutdown-error", "error", err)
	}

	<-ctx.Done()
}
