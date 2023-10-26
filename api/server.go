package api

import (
	"fmt"

	"nepse-backend/api/controllers"
	"nepse-backend/nepse/onlinekhabar"

	"github.com/gorilla/mux"
)

func New() *controllers.Server {
	s := &controllers.Server{}
	onlinestock := onlinekhabar.NewOkStock()
	s.Router = mux.NewRouter()
	s.OkStock = onlinestock
	return s
}

func Run() {
	s := New()
	fmt.Println("Server is running on Port 8080")
	s.InitRoutes()
	s.Run(":8080")
}
