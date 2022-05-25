package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/google/uuid"
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
	_, err := kv.DeleteTree(constructGroupKey(id, version), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"Deleted": id}, nil
}

func (db *Database) Config(config *Config) (*Config, error) {
	kv := db.cli.KV()

	dbkey, id := generateKey(config.Id, config.Version, "")
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
	if group.Id == "" {
		group.Id = uuid.New().String()
	}
	for _, v := range group.Configs {
		label := ""
		keys := make([]string, 0, len(v.Entries))
		for k, _ := range v.Entries {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			label += k + ":" + v.Entries[k] + ";"
		}

		label = label[:len(label)-1]
		dbkey, _ := generateKey(group.Id, group.Version, label)
		fmt.Println("OVDE JE DBKEY I ON GLASI : " + dbkey)
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		g := &api.KVPair{Key: dbkey, Value: data}
		_, err = kv.Put(g, nil)
		if err != nil {
			return nil, err
		}
	}

	return group, nil
}

func (ps *Database) GetGroup(id string, version string) (*Group, error) {
	kv := ps.cli.KV()
	cKey := constructKey(id, version, "1")
	cKey = cKey[:len(cKey)-2]
	fmt.Println("OVDE JE CKEY I ON GLASI: " + cKey)
	data, _, err := kv.List(cKey, nil)
	if err != nil {
		return nil, err
	}

	configs := []Config{}
	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, *config)
	}
	group := &Group{}
	group.Id = id
	group.Version = version
	group.Configs = configs

	return group, nil
}

func (ps *Database) GetAllConfigs() ([]*Config, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(allConfigs, nil)
	if err != nil {
		return nil, err
	}

	posts := []*Config{}
	for _, pair := range data {
		post := &Config{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (ps *Database) GetConfigsFromGroup(id string, version string, label string) ([]*Config, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(constructKey(id, version, label), nil)
	if err != nil {
		return nil, err
	}

	configs := []*Config{}
	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}
