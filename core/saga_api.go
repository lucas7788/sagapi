package core

import "github.com/ontio/sagapi/core/nasa"

var DefSagaApi *SagaApi

type SagaApi struct {
	Nasa      *nasa.Nasa
	SagaOrder *SagaOrder
}

func NewSagaApi() *SagaApi {
	return &SagaApi{
		Nasa:      nasa.NewNasa(),
		SagaOrder: NewSagaOrder(),
	}
}
