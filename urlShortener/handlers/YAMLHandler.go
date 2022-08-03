package handlers

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

type pathURL struct {
	Path string `yaml:"path"`
	Dest string `yaml:"url"`
}

func YAMLHandler(data []byte, mapFallback http.HandlerFunc) (http.HandlerFunc, error) {

	var paths []pathURL

	err := yaml.Unmarshal(data, &paths)

	if err != nil {
		return nil, err
	}

	yamlMap := map[string]string{}

	for _, path := range paths {
		yamlMap[path.Path] = path.Dest
	}

	return MapHandler(yamlMap, mapFallback), nil
}
