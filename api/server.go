package api

import (
	"github.com/HasanShahjahan/go-guest/api/utils"
	"os"

	"github.com/HasanShahjahan/go-guest/api/controllers"
	"github.com/HasanShahjahan/go-guest/api/seed"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const (
	logTag = "server"
)

var server = controllers.Server{}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		logging.Error(logTag, "Error getting env, not coming through %v", err)
	} else {
		logging.Info(logTag, "We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	seed.Load(server.DB)
	server.Run(":8080")

}
