package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadYaml(file string) map[interface{}]interface{} {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal(data, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return m
}
