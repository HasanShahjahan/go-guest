package api

import (
	"fmt"
	"github.com/HasanShahjahan/go-guest/api/config"
	"github.com/HasanShahjahan/go-guest/api/utils"
	"os"
	"runtime/debug"

	"github.com/HasanShahjahan/go-guest/api/controllers"
	"github.com/HasanShahjahan/go-guest/api/seed"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const (
	logTag = "[Server]"
)

var server = controllers.Server{}

func Run() {
	if err := config.LoadJSONConfig(config.Config); err != nil {
		logging.Fatal(logTag, "unable to load configuration. error=%v", err)
	}
	logging.Info(logTag, "configuration file loaded")

	logging.SetLogLevel(config.Config.LogLevel)

	err := godotenv.Load()
	if err != nil {
		logging.Error(logTag, "Error getting env, not coming through %v", err)
	} else {
		logging.Info(logTag, "We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	seed.Load(server.DB)
	server.Run(":8080")
}

func DoAPIPanicRecovery() {
	if r := recover(); r != nil {
		logMessage := fmt.Sprintf("API failed with error %s %s",
			r, string(debug.Stack()),
		)
		logging.Fatal(logTag, logMessage)
	}
}
