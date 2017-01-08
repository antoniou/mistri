package domain

type Pipeline interface {
	Create([]string) error
	Delete(interface{}) error
	Read(interface{}) interface{}
}

type PipelineFactory func(conf map[string]string) (Pipeline, error)
