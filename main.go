package main

import (
	"fmt"
	"log"
	"os"

	controllers "nepse-backend/api"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	if err != nil {
		fmt.Println("err", err)
		os.Exit(0)
	}
	if err := controllers.InitRouter().Run(); err != nil {
		fmt.Println("unable to init routes", err)
	}
}
