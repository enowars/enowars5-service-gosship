package handler

import (
	"checker/pkg/checker"
	"errors"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) PutFlag(message *checker.TaskMessage) (*checker.ResultMessage, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) GetFlag(message *checker.TaskMessage) (*checker.ResultMessage, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) PutNoise(message *checker.TaskMessage) (*checker.ResultMessage, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) GetNoise(message *checker.TaskMessage) (*checker.ResultMessage, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) Havoc(message *checker.TaskMessage) (*checker.ResultMessage, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) GetServiceInfo() *checker.InfoMessage {
	return &checker.InfoMessage{
		ServiceName:   "gosship",
		FlagVariants:  1,
		NoiseVariants: 1,
		HavocVariants: 1,
	}
}
