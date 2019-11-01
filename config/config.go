package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/kravcs/gogo/models"
)

var Configuration models.Config

func init() {
	var err error

	Configuration, err = LoadConfiguration("config.yml")

	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	fmt.Println("Config initialized")
}

func LoadConfiguration(filename string) (config models.Config, err error) {

	configFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	err = yaml.NewDecoder(configFile).Decode(&config)

	return config, err
}
