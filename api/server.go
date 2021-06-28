package api

import (
	"fmt"
	"log"
	"os"

	"github.com/HasanShahjahan/go-guest/api/controllers"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_NAME"))
	//seed.Load(server.DB)
	server.Run(":8080")

}