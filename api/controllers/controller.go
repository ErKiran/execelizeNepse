package controllers

import "nepse-backend/api/controllers/onlinekhabar"

type controller struct {
	OKController onlinekhabar.OKController
}

func NewController() *controller {
	return &controller{
		OKController: onlinekhabar.NewOKController(),
	}
}
