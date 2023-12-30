package watcher

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type SearchCriteria map[string]string

func NewSearchCriteria(fileName string) (SearchCriteria, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	rawCriteria := make(map[string]any)
	if err := yaml.Unmarshal(bytes, &rawCriteria); err != nil {
		return nil, err
	}

	// Convert from map[string]any to map[string]string
	criteria := make(map[string]string, len(rawCriteria))
	for key, value := range rawCriteria {
		criteria[key], err = asString(value)
		if err != nil {
			return nil, err
		}
	}
	return criteria, nil
}

func asString(val any) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	case int:
		return fmt.Sprintf("%d", v), nil
	default:
		return "", fmt.Errorf("unsupported type %T: %v", v, v)
	}
}
