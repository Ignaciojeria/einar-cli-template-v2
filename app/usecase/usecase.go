package usecase

import ioc "github.com/Ignaciojeria/einar-ioc"

var _ = ioc.Registry(NewUsecaseStruct)

type IUsecase interface {
	Execute() error
}

type usecaseStruct struct {
}

func NewUsecaseStruct() IUsecase {
	return usecaseStruct{}
}

func (u usecaseStruct) Execute() error {
	return nil
}
