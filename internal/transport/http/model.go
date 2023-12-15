package http

import "luna-values-storage/internal/domain"

type Value struct {
	ID   string      `json:"id"`
	Name string      `json:"name"`
	Type string      `json:"type"`
	Val  interface{} `json:"value"`
}

func (v *Value) IntoDomainType() *domain.Value {
	return &domain.Value{
		ID:    v.ID,
		Type:  v.Type,
		Name:  v.Name,
		Value: v.Val,
	}
}

func valueDomainToHttp(v *domain.Value) Value {
	return Value{
		ID:   v.ID,
		Name: v.Name,
		Type: v.Type,
		Val:  v.Value,
	}
}

type Variable struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	ValueID *string      `json:"value_id,omitempty"`
	Value   *interface{} `json:"value,omitempty"`
}

func (v *Variable) IntoDomainType() *domain.Variable {
	return &domain.Variable{
		ID:      v.ID,
		Name:    v.Name,
		Type:    v.Type,
		ValueID: v.ValueID,
		Value:   v.Value,
	}
}

func variableDomainToHttp(v *domain.Variable) Variable {
	return Variable{
		ID:      v.ID,
		Name:    v.Name,
		Type:    v.Type,
		ValueID: v.ValueID,
		Value:   v.Value,
	}
}
