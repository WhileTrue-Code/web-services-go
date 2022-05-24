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

func (db *Database) Config(config *Config) (*Config, error) {
	kv := db.cli.KV()

	dbkey, id := generateKey(config.Id, config.Version)
	config.Id = id

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	c := &api.KVPair{Key: dbkey, Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (db *Database) Group(group *Group) (*Group, error) {
	kv := db.cli.KV()

	dbkey, id := generateKey(group.Id, group.Version)
	group.Id = id

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	g := &api.KVPair{Key: dbkey, Value: data}
	_, err = kv.Put(g, nil)
	if err != nil {
		return nil, err
	}

	return group, nil
}
