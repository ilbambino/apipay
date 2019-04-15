package persistent

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultDB = "apipay"
)

// Client holds a connection to a database
type Client struct {
	mongoClient *mongo.Client
	db          *mongo.Database
}

// Connect setups the connection to the DB
func Connect(ctx context.Context, host string, port int, user, password string) (Client, error) {
	conn := Client{}
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	connStr := fmt.Sprintf("mongodb://%s:%d", host, port)
	ops := options.Client().ApplyURI(connStr)
	ops.SetAppName("ApiPay")

	if len(user) > 0 {
		creds := options.Credential{Username: user}
		if len(password) > 0 {
			creds.Password = password
		}
		ops.SetAuth(creds)
	}
	client, err := mongo.Connect(ctx, ops)
	if err != nil {
		return conn, err
	}
	conn.mongoClient = client
	conn.db = client.Database(defaultDB)

	return conn, nil
}

// IsErrorNoDBResults is just a small helper to check if the error is just
// that we didn't find anything on DB
func IsErrorNoDBResults(err error) bool {
	return err == mongo.ErrNoDocuments
}

// Close closes DB connection and it is not usable once this method is run
func (cl *Client) Close(ctx context.Context) error {
	err := cl.mongoClient.Disconnect(ctx)
	if err != nil {
		return err
	}
	cl.mongoClient = nil
	return nil
}

// UseDatabase changes the DB used in mongo
func (cl *Client) UseDatabase(dbName string) {

	cl.db = cl.mongoClient.Database(dbName)

}

// DropDatabase deletes the current DB. WARNING, use with care.
func (cl *Client) DropDatabase(ctx context.Context) error {

	return cl.db.Drop(ctx)
}
