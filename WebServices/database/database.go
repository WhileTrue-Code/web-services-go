package database

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

type Database struct {
	cli *api.Client
}

func New() (*Database, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Database{
		cli: client,
	}, nil
}

func (ps *Database) Get(id string, version string) (*Config, error) {
	kv := ps.cli.KV()

	pair, _, err := kv.Get(constructKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(pair.Value, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
