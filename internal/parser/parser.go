package parser

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ParseInsomniaCollection(path string) (*InsomniaCollection, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var collection InsomniaCollection
	err = yaml.Unmarshal(data, &collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}
