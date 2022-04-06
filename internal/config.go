package internal

import (
	"encoding/json"
	"log"
	"os"
	"shopingList/pkg/models"
)

func LoadConfigFromPath(configPath string) models.AppConfig {
	config := models.AppConfig{}

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalln("[ERROR]: can't open config file: ", err)
	}
	defer file.Close() // nolint: errcheck, gosec - not critic here
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatalln("[ERROR]: can't parse config file: ", err) // nolint: gocritic
	}

	return config
}
