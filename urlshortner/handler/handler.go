package handler

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type Map struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

type Config struct {
	Mappings []Map `yaml:"mappings"`
}

func MapHandler(pathToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if url, ok := pathToUrls[path]; ok {
			http.Redirect(rw, req, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(rw, req)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	data, err := parseYaml(yml)
	if err != nil {
		log.Fatal(err)
	}
	pathsToUrls := buildMap(data)
	return MapHandler(pathsToUrls, fallback), err
}

func parseYaml(yml []byte) (Config, error) {
	var data Config

	err := yaml.Unmarshal(yml, &data)
	return data, err
}

func buildMap(data Config) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, mapping := range data.Mappings {
		pathsToUrls[mapping.Path] = mapping.URL
	}
	return pathsToUrls
}
