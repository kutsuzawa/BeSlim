package controller

import (
	"io"

	"github.com/kutsuzawa/slim-load-recorder/interactor"
)

type Parser interface {
	Parse() (io.Reader, error)
}

type RequestReceive struct {
	parser Parser
}

func (rr *RequestReceive) GetRequest() (interactor.Request, error) {
	// TODO: infra層でparseメソッドを実装して、返り値のio.Readerをビジネスに
	// 使える形(interactor.Request) に変換して返す
	return interactor.Request{}, nil
}
