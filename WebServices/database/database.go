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

	config := &Config{}
	err = json.Unmarshal(pair.Value, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (db *Database) DeleteConfig(id string, version string) (map[string]string, error) {
	kv := db.cli.KV()
	_, err := kv.Delete(constructConfigKey(id, version), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (db *Database) DeleteConfigGroup(id string, version string) (map[string]string, error) {
	kv := db.cli.KV()
	_, err := kv.Delete(constructGroupKey(id, version), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"Deleted": id}, nil
}

func (db *Database) Config(config *Config) (*Config, error) {
	kv := db.cli.KV()

	dbkey, id := generateKey(config.Id, config.Version)
	config.Id = id

	data, err := json.Marshal(config)
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
