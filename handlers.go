package main

import (
	"apipay/model"
	"apipay/persistent"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

const (
	defaultTimeout = time.Second * 10
)

func getPayments(logger *zap.Logger, paymentDb persistent.Payments) func(ginCtx *gin.Context) {

	return func(ginCtx *gin.Context) {

		ctx, cancel := context.WithTimeout(ginCtx.Request.Context(), defaultTimeout)
		defer cancel()

		items, err := paymentDb.Last100(ctx)
		if err != nil {
			logger.Sugar().Warnw("get-payments-db", "error", err)
			ginCtx.Status(http.StatusInternalServerError)
		} else {
			ginCtx.JSON(http.StatusOK, items)
		}
	}
}

func getOnePayment(logger *zap.Logger, paymentDb persistent.Payments) func(ginCtx *gin.Context) {

	return func(ginCtx *gin.Context) {

		ctx, cancel := context.WithTimeout(ginCtx.Request.Context(), defaultTimeout)
		defer cancel()

		id := model.PaymentID(ginCtx.Param("paymentID"))
		// TODO Can we do some validation of the ID?

		item, err := paymentDb.Get(ctx, id)
		if err != nil {
			if persistent.IsErrorNoDBResults(err) {
				logger.Info("get-one-payments-db-not-found")
				ginCtx.Status(http.StatusNotFound)
			} else {
				logger.Sugar().Warnw("get-one-payments-db", "error", err)
				ginCtx.Status(http.StatusInternalServerError)

			}
		} else {
			ginCtx.JSON(http.StatusOK, item)
		}
	}
}

func updatePayment(logger *zap.Logger, paymentDb persistent.Payments) func(ginCtx *gin.Context) {

	return func(ginCtx *gin.Context) {

		ctx, cancel := context.WithTimeout(ginCtx.Request.Context(), defaultTimeout)
		defer cancel()

		id := model.PaymentID(ginCtx.Param("paymentID"))
		// TODO Can we do some validation of the ID?

		item, err := paymentDb.Get(ctx, id)
		if err != nil {
			if persistent.IsErrorNoDBResults(err) {
				logger.Info("get-one-payments-db-not-found")
				ginCtx.Status(http.StatusNotFound)
			} else {
				logger.Sugar().Warnw("get-one-payments-db", "error", err)
				ginCtx.Status(http.StatusInternalServerError)

			}
		} else {
			ginCtx.JSON(http.StatusOK, item)
		}
	}
}

func deletePayment(logger *zap.Logger, paymentDb persistent.Payments) func(ginCtx *gin.Context) {

	return func(ginCtx *gin.Context) {

		ctx, cancel := context.WithTimeout(ginCtx.Request.Context(), defaultTimeout)
		defer cancel()

		id := model.PaymentID(ginCtx.Param("paymentID"))
		// TODO Can we do some validation of the ID?

		_, err := paymentDb.Delete(ctx, id)
		if err != nil {
			logger.Sugar().Warnw("delete-payments-db", "error", err)
			ginCtx.Status(http.StatusInternalServerError)
		} else {
			ginCtx.Status(http.StatusNoContent)
		}
	}
}

func createPayment(logger *zap.Logger, paymentDb persistent.Payments) func(ginCtx *gin.Context) {

	return func(ginCtx *gin.Context) {

		ctx, cancel := context.WithTimeout(ginCtx.Request.Context(), defaultTimeout)
		defer cancel()

		received := &model.Payment{}
		if err := binding.JSON.Bind(ginCtx.Request, received); err != nil {
			logger.Sugar().Warnw("create-payments-db-json", "error", err)
			ginCtx.Status(http.StatusBadRequest)
			return
		}
		if !received.Valid() {
			logger.Warn("create-payments-db-invalid")
			ginCtx.Status(http.StatusBadRequest)
			return
		}
		err := paymentDb.Save(ctx, *received)
		if err != nil {
			logger.Sugar().Warnw("create-payments-db", "error", err)
			ginCtx.Status(http.StatusInternalServerError)
		} else {
			ginCtx.Status(http.StatusCreated) //TODO. Should we return the ID?
		}
	}
}
