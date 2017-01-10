package domain

type Actor interface {
	Run(interface{}) error
}
