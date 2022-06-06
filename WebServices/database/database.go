package database

import (
	"WebServices/tracer"
	"context"
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

func (ps *Database) Get(ctx context.Context, id string, version string) (*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConfigFromDatabase")
	defer span.Finish()

	ctxBase := tracer.ContextWithSpan(ctx, span)

	kv := ps.cli.KV()

	pair, _, err := kv.Get(constructKey(id, version, ""), nil)

	spanBase := tracer.StartSpanFromContext(ctxBase, "Get one config from database")
	defer spanBase.Finish()

	if err != nil {
		tracer.LogError(spanBase, err)
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(pair.Value, config)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	tracer.LogString("database_getConfigs", "Successful reading from database")
	return config, nil
}

// func (db *Database) DeleteConfig(id string, version string) (map[string]string, error) {
// 	kv := db.cli.KV()
// 	_, err := kv.Delete(constructConfigKey(id, version), nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return map[string]string{"Deleted": id}, nil
// }

// func (db *Database) DeleteConfigGroup(id string, version string) (map[string]string, error) {
// 	kv := db.cli.KV()
// 	_, err := kv.DeleteTree(constructGroupKey(id, version), nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return map[string]string{"Deleted": id}, nil
// }

func (db *Database) IdempotencyKey(ideKey *string) (*string, error) {
	kv := db.cli.KV()

	dbIdeKey := constructKey(*ideKey, "", "")
	fmt.Println("ISKONSTRUISANI DBIDEKEY IZGLEDA: " + dbIdeKey)

	byteIdeKey := []byte(*ideKey)

	iKey := &api.KVPair{Key: dbIdeKey, Value: byteIdeKey}
	_, err := kv.Put(iKey, nil)
	if err != nil {
		return nil, err
	}
	return ideKey, nil

}

func (db *Database) GetIdempotencyKey(ideKey *string) (*string, error) {
	kv := db.cli.KV()

	pair, _, err := kv.Get(constructKey(*ideKey, "", ""), nil)

	if err != nil {
		return nil, err
	}

	if pair == nil {
		return nil, nil
	}

	iK := string(pair.Value)

	return &iK, nil

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
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		databaseKey := dbkey + uuid.New().String()
		fmt.Println("OVDE JE DBKEY I ON GLASI : " + databaseKey)
		g := &api.KVPair{Key: databaseKey, Value: data}
		_, err = kv.Put(g, nil)
		if err != nil {
			return nil, err
		}
	}

	return group, nil
}

func (ps *Database) GetGroup(ctx context.Context, id string, version string) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "GetGroupFromDataBase")
	defer span.Finish()

	ctxBase := tracer.ContextWithSpan(context.Background(), span)

	kv := ps.cli.KV()
	cKey := constructKey(id, version, "1")
	cKey = cKey[:len(cKey)-3]
	fmt.Println("OVDE JE CKEY I ON GLASI: " + cKey)
	data, _, err := kv.List(cKey, nil)

	spanBase := tracer.StartSpanFromContext(ctxBase, "List method for group")
	defer spanBase.Finish()

	if err != nil {
		tracer.LogError(spanBase, err)
		return nil, err
	}

	configs := []Config{}
	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		configs = append(configs, *config)
	}
	group := &Group{}
	group.Id = id
	group.Version = version
	group.Configs = configs

	tracer.LogString("database_getConfigs", "Successful reading from database")
	return group, nil
}

func (ps *Database) GetAllConfigs(ctx context.Context) ([]*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllConfigs")
	defer span.Finish()

	ctxBase := tracer.ContextWithSpan(ctx, span)

	kv := ps.cli.KV()
	data, _, err := kv.List(allConfigs, nil)

	spanBase := tracer.StartSpanFromContext(ctxBase, "Get configs from base")
	defer spanBase.Finish()

	if err != nil {
		tracer.LogError(spanBase, err)
		return nil, err
	}

	configs := []*Config{}
	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		configs = append(configs, config)
	}
	tracer.LogString("database_getConfigs", "Successful reading from database")
	return configs, nil
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

// func (ps *Database) AddConfigsToGroup(id string, version string, config Config) (*Group, error) {

// 	group, error := ps.GetGroup(id, version)
// 	if error != nil {
// 		return nil, error
// 	}
// 	if len(group.Configs) < 1 {
// 		return nil, errors.New("Group doesn't exists!")
// 	}

// 	groupW := Group{}
// 	groupW.Id = id
// 	groupW.Version = version
// 	groupW.Configs = append(groupW.Configs, config)

// 	_, err := ps.Group(&groupW)

// 	if err != nil {
// 		return nil, err
// 	}

// 	ret, err := ps.GetGroup(id, version)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ret, nil
// }
