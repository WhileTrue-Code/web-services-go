package database

import (
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

func (db *Database) DeleteConfig(id string, version string) (map[string]string, error) {
	kv := db.cli.KV()
	_, err := kv.Delete(constructKey(id, version), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}
