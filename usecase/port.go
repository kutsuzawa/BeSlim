package usecase

import "github.com/kutsuzawa/slim-load-recorder/entity"

// InputPort defines inputPort
type InputPort interface {
	Handle(request Request) (Response, error)
}

// OutputPort defines outputPort
type OutputPort interface {
	Handle(loads []entity.Load) (Response, error)
}
