package dataset

import (
	"encoding/json"
	"os"
)

type TechDataset struct {
	Languages []string `json:"languages"`
	Databases []string `json:"databases"`
	DevOps    []string `json:"devops"`
	Auxillary    []string `json:"auxillary"`
}

func LoadTechDataset(path string) (map[string]bool, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var dataset TechDataset
	if err := json.Unmarshal(file, &dataset); err != nil {
		return nil, err
	}

	techSet := make(map[string]bool)
	for _, tech := range append(dataset.Languages, dataset.Databases...) {
		techSet[tech] = true
	}
	for _, tech := range append(dataset.DevOps, dataset.Auxillary...) {
		techSet[tech] = true
	}

	return techSet, nil
}
