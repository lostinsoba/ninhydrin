package controller

import (
	"lostinsoba/ninhydrin/internal/storage"
)

type Controller struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Controller {
	return &Controller{storage: storage}
}
