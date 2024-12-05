package cmd

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Program   string `json:"program"`
	SearchPid string `json:"search_pid"`
	Logpath   string `json:"logpath"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
