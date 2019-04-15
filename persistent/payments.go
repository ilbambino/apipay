package persistent

import (
	"apipay/model"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	defaultPaymentsCollection = "payments"

	defaultDBTimeout = time.Second * 3
)

// GetPayments is to get the Payments object (to interact with DB) with a
// given DB connection
func GetPayments(ctx context.Context, cl Client) (Payments, error) {

	obj := Payments{
		collection: cl.db.Collection(defaultPaymentsCollection),
	}

	err := obj.init(ctx)
	return obj, err
}

// Payments is the way to persist payments to DB
// it abstracts all the interactions to the DB. New methods should be added here
// depending on the needs
// also here is where the needed checks should be added
type Payments struct {
	collection *mongo.Collection
}

// init the collection, setting up indices…
func (p *Payments) init(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, defaultDBTimeout)
	defer cancel()

	indexOps := options.Index()
	indexOps.SetBackground(true)
	indexOps.SetUnique(true)

	index := mongo.IndexModel{
		Options: indexOps,
		Keys:    bson.D{{Key: "id", Value: 1}},
	}

	//TODO depending on the usage more indices should be created
	_, err := p.collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}
	return nil
}

// Save saves a payment to DB. If it is already there it will replace it, otherwise
// it will save it
func (p *Payments) Save(ctx context.Context, obj model.Payment) error {

	ctx, cancel := context.WithTimeout(ctx, defaultDBTimeout)
	defer cancel()

	filter := bson.D{{"id", obj.ID}}
	ops := options.Replace()
	ops.SetUpsert(true)
	_, err := p.collection.ReplaceOne(ctx, filter, obj, ops)
	if err != nil {
		return err
	}
	return nil
}

// Get tries to find a payment in the DB and returns it
func (p *Payments) Get(ctx context.Context, id model.PaymentID) (model.Payment, error) {

	ctx, cancel := context.WithTimeout(ctx, defaultDBTimeout)
	defer cancel()

	filter := bson.D{{"id", id}}

	var result model.Payment

	err := p.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// Delete tries to delete a payment in the DB, returns the number of deleted items
func (p *Payments) Delete(ctx context.Context, id model.PaymentID) (int64, error) {

	ctx, cancel := context.WithTimeout(ctx, defaultDBTimeout)
	defer cancel()

	filter := bson.D{{"id", id}}

	res, err := p.collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// Last100 gets the last 100 items of the list
func (p *Payments) Last100(ctx context.Context) ([]*model.Payment, error) {

	ctx, cancel := context.WithTimeout(ctx, defaultDBTimeout)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(100)
	findOptions.Sort = bson.D{{"_id", -1}} //descending

	var results []*model.Payment

	cur, err := p.collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {

		var elem model.Payment
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(ctx)

	return results, nil
}
