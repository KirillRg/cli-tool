package parser

import (
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Структура регулярного выражения для переменных среды
var envVarRegex = regexp.MustCompile(`{{\s*\_\.\s*([a-zA-Z0-9_-]+)\s*}}`)

// Функция подстановки переменных
func SubstituteEnvVariables(input string, env map[string]string) string {
	return envVarRegex.ReplaceAllStringFunc(input, func(s string) string {
		match := envVarRegex.FindStringSubmatch(s)
		if len(match) == 2 {
			if val, exists := env[match[1]]; exists {
				return val
			}
		}
		return s
	})
}

// Функция парсинга Insomnia-коллекции с учетом переменных среды
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

	// Подстановка переменных после парсинга
	for i, req := range collection.Collection {
		req.URL = SubstituteEnvVariables(req.URL, collection.Environments.Data)

		for j, hdr := range req.Headers {
			req.Headers[j].Name = SubstituteEnvVariables(hdr.Name, collection.Environments.Data)
			req.Headers[j].Value = SubstituteEnvVariables(hdr.Value, collection.Environments.Data)
		}

		for k, param := range req.Parameters {
			req.Parameters[k].Name = SubstituteEnvVariables(param.Name, collection.Environments.Data)
			req.Parameters[k].Value = SubstituteEnvVariables(param.Value, collection.Environments.Data)
		}

		req.Body.Text = SubstituteEnvVariables(req.Body.Text, collection.Environments.Data)

		collection.Collection[i] = req
	}

	return &collection, nil
}
