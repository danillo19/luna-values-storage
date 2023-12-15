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

type VariableRepository struct {
	collection *mongo.Collection
}

func NewVariableRepository(db *mongo.Database) VariableRepository {
	return VariableRepository{collection: db.Collection("variables")}
}

func (v VariableRepository) Get(ctx context.Context, id string) (*domain.Variable, error) {
	oidFromHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse value oid: %v", err)
	}

	raw := new(Variable)

	if err := v.collection.FindOne(ctx, bson.M{"_id": oidFromHex}).Decode(raw); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}

		return nil, fmt.Errorf("failed to get variable with id: %s: %v", id, err)
	}

	return raw.IntoDomainType(), nil
}

func (v VariableRepository) Set(ctx context.Context, variable *domain.Variable) (*domain.Variable, error) {
	raw := rawVarFromDomainType(variable)
	if raw.ID == nil {
		oid := primitive.NewObjectID()
		raw.ID = &oid
	}

	opts := mongoopts.Update().SetUpsert(true)
	_, err := v.collection.UpdateOne(ctx, bson.M{"_id": raw.ID}, bson.M{"$set": raw}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to set variable: %v", err)
	}

	return raw.IntoDomainType(), nil
}

func (v VariableRepository) Delete(ctx context.Context, id string) error {
	oidFromHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to parse variable oid: %v", err)
	}

	_, err = v.collection.DeleteOne(ctx, bson.M{"_id": oidFromHex})
	if err != nil {
		return fmt.Errorf("failed to delete variable with id: %s: %v", id, err)
	}

	return nil
}
