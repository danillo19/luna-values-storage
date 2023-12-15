package domain

type Variable struct {
	ID   string
	Name string
	Type string
	
	ValueID *string
	Value   *interface{}
}
