package domain

type Pipeline interface {
	Create(interface{}) error
	Delete(interface{}) error
	Read(interface{}) interface{}
}
