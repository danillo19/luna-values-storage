package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"luna-values-storage/internal/domain"
)

type Variable struct {
	ID      *primitive.ObjectID `bson:"_id,omitempty"`
	Name    string              `bson:"name"`
	Type    string              `bson:"type"`
	Value   *interface{}        `bson:"value,omitempty"`
	ValueId *string             `bson:"value_id,omitempty"`
}

func (v *Variable) IntoDomainType() *domain.Variable {

	return &domain.Variable{
		ID:      v.ID.Hex(),
		Name:    v.Name,
		Type:    v.Type,
		ValueID: v.ValueId,
		Value:   v.Value,
	}
}

func rawVarFromDomainType(v *domain.Variable) *Variable {
	variable := &Variable{
		Name:    v.Name,
		Type:    v.Type,
		Value:   v.Value,
		ValueId: v.ValueID,
	}

	if v.ID != "" {
		oid, _ := primitive.ObjectIDFromHex(v.ID)
		variable.ID = &oid
	}

	return variable
}

type Value struct {
	ID    *primitive.ObjectID `bson:"_id,omitempty"`
	Name  string              `bson:"name"`
	Type  string              `bson:"type"`
	Value interface{}         `bson:"value"`
}

func (v *Value) IntoDomainType() *domain.Value {
	return &domain.Value{
		ID:    v.ID.Hex(),
		Name:  v.Name,
		Type:  v.Type,
		Value: v.Value,
	}
}

func rawValueFromDomainType(val *domain.Value) *Value {
	value := &Value{
		Name:  val.Name,
		Type:  val.Type,
		Value: val.Value,
	}

	if val.ID != "" {
		oid, _ := primitive.ObjectIDFromHex(val.ID)
		value.ID = &oid
	}

	return value
}
