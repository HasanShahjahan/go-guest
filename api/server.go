package api

import (
	"fmt"
	"log"
	"os"

	"github.com/HasanShahjahan/go-guest/api/controllers"
	"github.com/HasanShahjahan/go-guest/api/seed"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not coming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	fmt.Println(os.Getenv("DB_DRIVER"))
	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	seed.Load(server.DB)
	server.Run(":8080")

}
