package main

import (
	"apipay/model"
	"apipay/persistent"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

const (
	defaultTimeout = time.Second * 10
)

// getPayments handler for getting all the payments
// @Summary Get the last 100 payments
// @Description Gets the last 100 payments. It should have a pagination
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Payment
// @Failure 404 {object} APIError "Can not find ID"
// @Failure 500 {object} APIError "Cannot process the request"
// @Router /payments [get]
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

// getOnePayment handler for getting one Payment by ID
// @Summary Get a Payment by ID
// @Accept  json
// @Produce  json
// @Param paymentID path string true "Payment ID"
// @Success 200 {object} model.Payment
// @Failure 404 {object} APIError "Can not find ID"
// @Failure 500 {object} APIError "Cannot process the request"
// @Router /payments/{paymentID} [get]
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

// updatePayment handler for getting one Payment by ID
// @Summary Update a Payment by ID
// @Accept  json
// @Produce  json
// @Param paymentID path string true "Payment ID"
// @Param payment body model.Payment true "The payment to be updated"
// @Success 201 {object} model.Payment "Empty result, it is not the created object" TODO fix this
// @Failure 400 {object} APIError "Invalid payment received"
// @Failure 404 {object} APIError "Can not find ID"
// @Failure 500 {object} APIError "Cannot process the request"
// @Router /payments/{paymentID} [put]
func updatePayment(logger *zap.Logger, paymentDb persistent.Payments) func(ginCtx *gin.Context) {

	return func(ginCtx *gin.Context) {

		ctx, cancel := context.WithTimeout(ginCtx.Request.Context(), defaultTimeout)
		defer cancel()

		id := model.PaymentID(ginCtx.Param("paymentID"))
		// TODO Can we do some validation of the ID?

		received := &model.Payment{}
		if err := binding.JSON.Bind(ginCtx.Request, received); err != nil {
			logger.Sugar().Warnw("update-payments-db-json", "error", err)
			ginCtx.Status(http.StatusBadRequest)
			return
		}
		if !received.Valid() {
			logger.Warn("update-payments-db-invalid")
			ginCtx.Status(http.StatusBadRequest)
			return
		}
		if received.ID != id {
			logger.Warn("update-payments-db-id-mismatch")
			ginCtx.Status(http.StatusBadRequest)
			return
		}

		err := paymentDb.Update(ctx, *received)
		fmt.Println("updated", err)
		if err != nil {
			if persistent.IsErrorNoDBResults(err) {
				logger.Info("get-one-payments-db-not-found")
				ginCtx.Status(http.StatusNotFound)
			} else {
				logger.Sugar().Warnw("get-one-payments-db", "error", err)
				ginCtx.Status(http.StatusInternalServerError)
			}
		} else {
			ginCtx.Status(http.StatusOK)
		}
	}
}

// deletePayment handler for deleting one Payment by ID
// @Summary Delete a Payment by ID
// @Accept  json
// @Produce  json
// @Param paymentID path string true "Payment ID"
// @Success 204 {object} model.Payment "Empty result, it is not the created object" TODO fix this
// @Failure 404 {object} APIError "Can not find ID"
// @Failure 500 {object} APIError "Cannot process the request"
// @Router /payments/{paymentID} [delete]
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

// createPayment handler for creating a new Payment
// @Summary Create  a new Payment
// @Accept  json
// @Produce  json
// @Param payment body model.Payment true "The payment to be created"
// @Success 204 {object} model.Payment "Empty result, it is not the created object" TODO fix this
// @Failure 400 {object} APIError "Payment with invalid format"
// @Failure 500 {object} APIError "Cannot process the request"
// @Router /payments [post]
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
