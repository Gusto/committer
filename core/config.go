package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Tasks []Task
}

func NewConfig(filename string) (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	config := Config{}

	yamlLocation := dir + "/" + filename
	yamlFile, err := ioutil.ReadFile(yamlLocation)
	if err != nil {
		log.Printf("Error opening config file: %v", err)
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error parsing YAML from %v: %v", yamlLocation, err)
		return nil, err
	}

	return &config, nil
}
