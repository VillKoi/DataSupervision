package main

import (
	"os"

	"gopkg.in/yaml.v3"

	"datasupervision/internal/controller"
)

type config struct {
	Server controller.Config `yaml:"server"`
}

func configure(fileName string) (*config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var cnf config

	if errUnmarshal := yaml.Unmarshal(data, &cnf); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return &cnf, nil
}
