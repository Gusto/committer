package core

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Tasks []Task
}

func NewConfigFromFile(filename string) (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	yamlLocation := dir + "/" + filename
	yamlFile, err := ioutil.ReadFile(yamlLocation)
	if err != nil {
		log.Printf("Error opening config file: %v", err)
		return nil, err
	}

	return NewConfig(yamlFile)
}

func NewConfig(content []byte) (*Config, error) {
	config := Config{}

	if err := yaml.Unmarshal(content, &config); err != nil {
		log.Printf("Error parsing YAML %v", err)
		return nil, err
	}

	return ValidateConfig(&config)
}

func ValidateConfig(config *Config) (*Config, error) {
	if len(config.Tasks) == 0 {
		return nil, errors.New("There must be a tasks array present in the config file.")

	}

	for _, task := range config.Tasks {
		if task.Name == "" {
			return nil, errors.New("All tasks must have a \"name\" key.")
		}

		if task.Command == "" {
			return nil, errors.New("All tasks must have a \"command\" key.")
		}

		if task.Fix.Command != "" {
			if task.Fix.Output == "" {
				return nil, errors.New("All tasks with a \"fix.command\" must specify a \"fix.output\" regex to show the autocorrect output.")
			}

			if task.Fix.Files == "" {
				return nil, errors.New("All tasks with a \"fix.command\" must specify a \"fix.files\" regex to determine if it should be run on change.")
			}
		}
	}

	return config, nil
}
