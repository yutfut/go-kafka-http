package conf

import (
	"encoding/json"
	"os"
)

type Conf struct {
	Kafka struct {
		Host		string	`json:"Host"`
		Port		int		`json:"Port"`
	} `json:"Kafka"`
	Main struct {
		HTTPPort	int		`json:"HTTPPort"`
		Host		string	`json:"Host"`
	} `json:"Main"`
}

func ReadConf(path string) (*Conf, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return &Conf{}, err
	}

	response := &Conf{}

	if err = json.Unmarshal(data, response); err != nil {
		return &Conf{}, err
	}

	return response, nil
}