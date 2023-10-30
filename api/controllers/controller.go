package controllers

import (
	"nepse-backend/api/controllers/cricket"
	"nepse-backend/api/controllers/onlinekhabar"
)

type controller struct {
	OKController       onlinekhabar.OKController
	CricInfoController cricket.CrickInfoController
}

func NewController() *controller {
	return &controller{
		OKController:       onlinekhabar.NewOKController(),
		CricInfoController: cricket.NewCricInfoController(),
	}
}
