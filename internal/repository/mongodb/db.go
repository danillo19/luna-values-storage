package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type DB struct {
	client             *mongo.Client
	VariableRepository VariableRepository
	ValueRepository    ValueRepository
}

func NewDBWrapper(client *mongo.Client, isCacheEnabled bool) *DB {
	db := client.Database("luna")

	return &DB{
		client:             client,
		VariableRepository: NewVariableRepository(db),
		ValueRepository:    NewValueRepository(db),
	}
}

func NewDB(context context.Context, URL string, isCacheEnabled bool) (*DB, func(), error) {
	mongoClient, err := NewMongoClient(context, URL)
	if err != nil {
		return nil, nil, err
	}

	disconnect := func() {
		if err = mongoClient.Disconnect(context); err != nil {
			panic(err)
		}
	}

	db := NewDBWrapper(mongoClient, isCacheEnabled)

	return db, disconnect, nil
}

type TransactionCallback func(sessionContext mongo.SessionContext) (interface{}, error)

func (db *DB) RunTransaction(ctx context.Context, callback func(sessionContext mongo.SessionContext) error) error {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := db.client.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		return nil, callback(sessionContext)
	}, txnOpts)

	return err
}
