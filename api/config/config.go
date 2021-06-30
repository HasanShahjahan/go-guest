package config

import (
	"encoding/json"
	"errors"
	"github.com/HasanShahjahan/go-guest/api/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

var Config = &struct {
	LogLevel string `json:"logLevel"`
}{}

const (
	logTag = "[config]"
)

var (
	// errNilConfig is returned when a nil reference is passed in as Un/Marshaler reference
	errNilConfig = errors.New("config object is empty")
)

// LoadJSONConfig gets your config from the json file sitting in the same folder of the app
func LoadJSONConfig(config interface{}) error {
	if config == nil {
		return errNilConfig
	}

	// try to get filename from env variable
	filename := os.Getenv("CONFIG_PATH")
	if filename != "" {
		logging.Info(logTag, "loading config from envVar=%s, file=%s", "CONFIG_PATH", filename)
		return LoadJSONFile(filename, config)
	}

	// try to load config file with the same name as executable
	absFileName, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}

	fileName := absFileName + ".json"
	logging.Info(logTag, "loading config from \"%s\"", fileName)

	return LoadJSONFile(fileName, config)
}

// LoadJSONFile gets your config from the json file, and fills your struct with the option
func LoadJSONFile(filename string, config interface{}) error {
	if config == nil {
		return errNilConfig
	}

	absFileName, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadFile(absFileName)
	if err != nil {
		return err
	}

	return json.Unmarshal(UpdateFromEnv(bytes), config)
}

// UpdateFromEnv is used during testing to update the config from enviromental settings (needed for Go.CD and docker)
func UpdateFromEnv(bytes []byte) []byte {
	configAsString := string(bytes)
	return []byte(configAsString)
}
