package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongoopts "go.mongodb.org/mongo-driver/mongo/options"
	"luna-values-storage/internal/domain"
)

type ValueRepository struct {
	collection *mongo.Collection
}

func NewValueRepository(db *mongo.Database) ValueRepository {
	return ValueRepository{collection: db.Collection("values")}
}

func (v ValueRepository) Get(ctx context.Context, id string) (*domain.Value, error) {
	oidFromHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse value oid: %v", err)
	}

	valueRaw := new(Value)

	if err := v.collection.FindOne(ctx, bson.M{"_id": oidFromHex}).Decode(valueRaw); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}

		return nil, fmt.Errorf("failed to get value with id: %s: %v", id, err)
	}

	return valueRaw.IntoDomainType(), nil
}

func (v ValueRepository) Set(ctx context.Context, value *domain.Value) (*domain.Value, error) {
	raw := rawValueFromDomainType(value)
	if raw.ID == nil {
		oid := primitive.NewObjectID()
		raw.ID = &oid
	}

	opts := mongoopts.Update().SetUpsert(true)
	_, err := v.collection.UpdateOne(ctx, bson.M{"_id": raw.ID}, bson.M{"$set": raw}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to set value : %v", err)
	}

	return raw.IntoDomainType(), nil
}

func (v ValueRepository) Delete(ctx context.Context, id string) error {
	oidFromHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to parse value oid: %v", err)
	}

	_, err = v.collection.DeleteOne(ctx, bson.M{"_id": oidFromHex})
	if err != nil {
		return fmt.Errorf("failed to delete value with id: %s: %v", id, err)
	}

	return nil
}
