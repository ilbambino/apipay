package persistent

import (
	"apipay/model"
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mongoHost = "localhost"
	mongoPort = 27017
)

func createTestDB(ctx context.Context, dbName string) (Client, error) {

	client, err := Connect(ctx, mongoHost, mongoPort, "", "")
	if err != nil {
		return Client{}, err
	}

	client.UseDatabase(dbName)

	err = client.DropDatabase(ctx) // for the test we want an empty DB every time
	if err != nil {
		return Client{}, err
	}
	return client, nil
}

func testPayment(id model.PaymentID) model.Payment {

	return model.Payment{
		Type: "Payment",
		ID:   id,
	}
}
func TestEmpty(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), defaultDBTimeout)
	defer cancel()

	client, err := createTestDB(ctx, "emptyDB")
	assert.NoError(t, err, "We can connect to DB")

	paymentsDB, err := GetPayments(ctx, client)
	assert.NoError(t, err, "We can init DB")

	payments, err := paymentsDB.Last100(ctx)
	assert.NoError(t, err, "We can fetch list from DB")
	assert.Equal(t, 0, len(payments), "We got nothing from DB")
}

func TestSave(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), defaultDBTimeout)
	defer cancel()

	client, err := createTestDB(ctx, "oneItemDB")
	assert.NoError(t, err, "We can connect to DB")

	paymentsDB, err := GetPayments(ctx, client)
	assert.NoError(t, err, "We can init DB")

	payments, err := paymentsDB.Last100(ctx)
	assert.NoError(t, err, "We can fetch list from DB")
	assert.Equal(t, 0, len(payments), "We got nothing from DB")

	payment1 := testPayment(model.PaymentID("12345"))

	err = paymentsDB.Save(ctx, payment1)
	assert.NoError(t, err, "We can save one item to DB")

	payments, err = paymentsDB.Last100(ctx)
	assert.NoError(t, err, "We can fetch list from DB")
	assert.Equal(t, 1, len(payments), "We got what we inserted")

}

func TestInsertDelete(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), defaultDBTimeout)
	defer cancel()

	client, err := createTestDB(ctx, "insertDeleteDB")
	assert.NoError(t, err, "We can connect to DB")

	paymentsDB, err := GetPayments(ctx, client)
	assert.NoError(t, err, "We can init DB")

	payment1 := testPayment(model.PaymentID("12345"))

	err = paymentsDB.Save(ctx, payment1)
	assert.NoError(t, err, "We can save one item to DB")

	payments, err := paymentsDB.Last100(ctx)
	assert.NoError(t, err, "We can fetch list from DB")
	assert.Equal(t, 1, len(payments), "We got what we inserted")

	numberDeleted, err := paymentsDB.Delete(ctx, payment1.ID)
	assert.NoError(t, err, "We can delete from DB")
	assert.Equal(t, int64(1), numberDeleted, "We deleted one item")

	payments, err = paymentsDB.Last100(ctx)
	assert.NoError(t, err, "We can fetch list from DB")
	assert.Equal(t, 0, len(payments), "List is empty now")

	numberDeleted, err = paymentsDB.Delete(ctx, payment1.ID)
	assert.NoError(t, err, "We can try to delete from DB")
	assert.Equal(t, int64(0), numberDeleted, "We deleted zero items")
}

func TestGet(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), defaultDBTimeout)
	defer cancel()

	client, err := createTestDB(ctx, "getDB")
	assert.NoError(t, err, "We can connect to DB")

	paymentsDB, err := GetPayments(ctx, client)
	assert.NoError(t, err, "We can init DB")

	payment1 := testPayment(model.PaymentID("12345"))

	_, err = paymentsDB.Get(ctx, payment1.ID)
	assert.True(t, IsErrorNoDBResults(err), "We didn't get anything from DB")

	err = paymentsDB.Save(ctx, payment1)
	assert.NoError(t, err, "We can save one item to DB")

	dbItem, err := paymentsDB.Get(ctx, payment1.ID)
	assert.NoError(t, err, "We can get items from DB")
	assert.True(t, reflect.DeepEqual(payment1, dbItem), "We loaded from DB what we saved")
}

func TestInsertUpdate(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), defaultDBTimeout)
	defer cancel()

	client, err := createTestDB(ctx, "insertUpdateDB")
	assert.NoError(t, err, "We can connect to DB")

	paymentsDB, err := GetPayments(ctx, client)
	assert.NoError(t, err, "We can init DB")

	payment1 := testPayment(model.PaymentID("12345"))

	err = paymentsDB.Save(ctx, payment1)
	assert.NoError(t, err, "We can save one item to DB")

	payments, err := paymentsDB.Last100(ctx)
	assert.NoError(t, err, "We can fetch list from DB")
	assert.Equal(t, 1, len(payments), "We got what we inserted")

	payment1.OrganisationID = "newOrg"
	err = paymentsDB.Save(ctx, payment1)
	assert.NoError(t, err, "We can update from DB")

	dbItem, err := paymentsDB.Get(ctx, payment1.ID)
	assert.NoError(t, err, "We can get items from DB")
	assert.True(t, reflect.DeepEqual(payment1, dbItem), "We loaded from DB what we saved")

}
func TestList100(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), defaultDBTimeout)
	defer cancel()

	client, err := createTestDB(ctx, "list100DB")
	assert.NoError(t, err, "We can connect to DB")

	paymentsDB, err := GetPayments(ctx, client)
	assert.NoError(t, err, "We can init DB")

	payment1 := testPayment(model.PaymentID("12345"))

	payments, err := paymentsDB.Last100(ctx)
	assert.NoError(t, err, "We can fetch list from DB")
	assert.Equal(t, 0, len(payments), "We got nothing from DB")

	err = paymentsDB.Save(ctx, payment1)
	assert.NoError(t, err, "We can save one item to DB")

	payments, err = paymentsDB.Last100(ctx)
	assert.NoError(t, err, "We can fetch list from DB")
	assert.Equal(t, 1, len(payments), "We got one item from DB")

	payment2 := testPayment(model.PaymentID("12346"))
	err = paymentsDB.Save(ctx, payment2)
	assert.NoError(t, err, "We can save one item to DB")

	payments, err = paymentsDB.Last100(ctx)
	assert.NoError(t, err, "We can fetch list from DB")
	assert.Equal(t, 2, len(payments), "We got two items from DB")

	payment3 := testPayment(model.PaymentID("12347"))
	err = paymentsDB.Save(ctx, payment3)
	assert.NoError(t, err, "We can save one item to DB")

	payments, err = paymentsDB.Last100(ctx)
	assert.NoError(t, err, "We can fetch list from DB")
	assert.Equal(t, 3, len(payments), "We got three items from DB")

}
