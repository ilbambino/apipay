package main

import (
	"apipay/config"
	"apipay/model"
	"apipay/persistent"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	testTimeOut = time.Second * 10
)

func createSupportItems(ctx context.Context, dbName string) (persistent.Payments, *zap.Logger, error) {

	err := config.Load()
	if err != nil {
		return persistent.Payments{}, nil, err
	}

	client, err := persistent.Connect(ctx,
		viper.GetString(config.MongoHost),
		viper.GetInt(config.MongoPort),
		viper.GetString(config.MongoUser),
		viper.GetString(config.MongoPassword))

	if err != nil {
		return persistent.Payments{}, nil, err
	}

	client.UseDatabase(dbName)
	paymentsDb, err := persistent.GetPayments(ctx, client)
	if err != nil {
		return persistent.Payments{}, nil, err
	}

	err = client.DropDatabase(ctx) // for the test we want an empty DB every time
	if err != nil {
		return persistent.Payments{}, nil, err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return persistent.Payments{}, nil, err
	}
	return paymentsDb, logger, nil
}

func testPayment(id model.PaymentID) model.Payment {

	return model.Payment{
		Type:           "Payment",
		ID:             id,
		OrganisationID: "testOrg",
	}
}
func TestGetList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeOut)
	defer cancel()

	db, logger, err := createSupportItems(ctx, "api_list")
	assert.NoError(t, err, "We can init the needed deps")

	router := getHandler(logger, db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/payments/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	obj := []model.Payment{}
	err = json.Unmarshal(w.Body.Bytes(), &obj)
	assert.NoError(t, err, "We can unmarshal the json")

	assert.Equal(t, 0, len(obj))
}

func TestGetOneInsertOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeOut)
	defer cancel()

	db, logger, err := createSupportItems(ctx, "api_one")
	assert.NoError(t, err, "We can init the needed deps")

	router := getHandler(logger, db)

	payment1 := testPayment(model.PaymentID("12345"))

	// get a payment (no payment available)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/payments/"+string(payment1.ID), nil)
	assert.NoError(t, err, "We can can the http request")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// create a new payment
	payment1Json, err := json.Marshal(payment1)
	assert.NoError(t, err, "We can marshal to json")
	req, err = http.NewRequest("POST", "/payments/", bytes.NewBuffer(payment1Json))
	w = httptest.NewRecorder()
	assert.NoError(t, err, "We can can the http request")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// get payment again, this time it is there
	req, err = http.NewRequest("GET", "/payments/"+string(payment1.ID), nil)
	w = httptest.NewRecorder()
	assert.NoError(t, err, "We can can the http request")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	obj := model.Payment{}
	err = json.Unmarshal(w.Body.Bytes(), &obj)
	assert.NoError(t, err, "We can unmarshal the json")

}

func TestInsertUpdateGetAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeOut)
	defer cancel()

	db, logger, err := createSupportItems(ctx, "api_insertupdategetall")
	assert.NoError(t, err, "We can init the needed deps")

	router := getHandler(logger, db)

	payment1 := testPayment(model.PaymentID("12345"))

	// get a payment (no payment available)
	w := httptest.NewRecorder()

	// create a new payment
	payment1Json, err := json.Marshal(payment1)
	assert.NoError(t, err, "We can marshal to json")
	req, err := http.NewRequest("POST", "/payments/", bytes.NewBuffer(payment1Json))
	assert.NoError(t, err, "We can can the http request")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// get payment again, this time it is there
	req, err = http.NewRequest("GET", "/payments/"+string(payment1.ID), nil)
	w = httptest.NewRecorder()
	assert.NoError(t, err, "We can can the http request")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	obj := model.Payment{}
	err = json.Unmarshal(w.Body.Bytes(), &obj)
	assert.NoError(t, err, "We can unmarshal the json")

}
