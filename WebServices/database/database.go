package database

import (
	"WebServices/tracer"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

type Database struct {
	cli *api.Client
}

//aca_lukas<3belo123
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

	spanBase := tracer.StartSpanFromContext(ctxBase, "Get one config from database")
	pair, _, err := kv.Get(constructKey(id, version, ""), nil)

	if pair == nil {
		tracer.LogError(spanBase, fmt.Errorf("ne postoji konfiguracija."))
		return nil, fmt.Errorf("ne postoji konfiguracija.")
	}

	if err != nil {
		tracer.LogError(spanBase, err)
		return nil, err
	}
	spanBase.Finish()
	config := &Config{}
	err = json.Unmarshal(pair.Value, config)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	tracer.LogString("database_getConfigs", "Successful reading from database")
	return config, nil
}

func (db *Database) DeleteConfig(ctx context.Context, id string, version string) (map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "deleteConfigFromDatabase")
	defer span.Finish()

	ctxBase := tracer.ContextWithSpan(ctx, span)

	kv := db.cli.KV()

	config, err1 := db.Get(ctx, id, version)
	if err1 != nil {
		tracer.LogError(span, fmt.Errorf("Config does not exist!"))
		return nil, err1
	}
	spanBase := tracer.StartSpanFromContext(ctxBase, "Delete one config from database")
	_, err := kv.Delete(constructConfigKey(config.Id, config.Version), nil)

	if err != nil {
		tracer.LogError(spanBase, err)
		return nil, err
	}
	spanBase.Finish()
	tracer.LogString("database_deleteConfig", "Successfully deleted from database")

	return map[string]string{"Deleted": id}, nil
}

func (db *Database) DeleteConfigGroup(ctx context.Context, id string, version string) (map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "deleteGroupFromDatabase")
	defer span.Finish()

	ctxBase := tracer.ContextWithSpan(ctx, span)
	group, err1 := db.GetGroup(ctx, id, version)
	if err1 != nil {
		tracer.LogError(span, fmt.Errorf("Group does not exist!"))
		return nil, err1
	}

	kv := db.cli.KV()

	spanBase := tracer.StartSpanFromContext(ctxBase, "Delete one group from database")
	_, err := kv.DeleteTree(constructGroupKey(group.Id, group.Version), nil)

	if err != nil {
		tracer.LogError(spanBase, err)
		return nil, err
	}
	spanBase.Finish()
	return map[string]string{"Deleted": id}, nil
}

func (db *Database) IdempotencyKey(ctx context.Context, ideKey *string) (*string, error) {
	span := tracer.StartSpanFromContext(ctx, "DB create idempotency-key")
	defer span.Finish()

	ctxBase := tracer.ContextWithSpan(context.Background(), span)

	kv := db.cli.KV()

	dbIdeKey := constructKey(*ideKey, "", "")
	// fmt.Println("ISKONSTRUISANI DBIDEKEY IZGLEDA: " + dbIdeKey)

	byteIdeKey := []byte(*ideKey)

	spanF := tracer.StartSpanFromContext(ctxBase, "Put idempotency-key in DB")
	iKey := &api.KVPair{Key: dbIdeKey, Value: byteIdeKey}

	_, err := kv.Put(iKey, nil)

	if err != nil {
		tracer.LogError(spanF, err)
		return nil, err
	}
	spanF.Finish()
	tracer.LogString("database-IdeKeySave", "Idempotency-key is saved.")
	return ideKey, nil

}

func (db *Database) GetIdempotencyKey(ctx context.Context, ideKey *string) (*string, error) {
	span := tracer.StartSpanFromContext(ctx, "DB get idempotency-key")
	defer span.Finish()

	ctxDB := tracer.ContextWithSpan(context.Background(), span)
	kv := db.cli.KV()

	spanF := tracer.StartSpanFromContext(ctxDB, "Put idempotency-key in DB")
	pair, _, err := kv.Get(constructKey(*ideKey, "", ""), nil)

	if err != nil {
		tracer.LogError(spanF, err)
		return nil, err
	}

	if pair == nil {
		tracer.LogError(spanF, err)
		return nil, nil
	}
	spanF.Finish()
	iK := string(pair.Value)

	span.LogFields(
		tracer.LogString("database-IdeKeyGet", "Idempotency-key is taken."),
	)
	return &iK, nil

}

func (db *Database) Config(ctx context.Context, config *Config) (*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "Save config DB")
	defer span.Finish()

	ctxDB := tracer.ContextWithSpan(context.Background(), span)
	kv := db.cli.KV()

	dbkey, id := generateKey(config.Id, config.Version, "")
	config.Id = id

	data, err := json.Marshal(config)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	c := &api.KVPair{Key: dbkey, Value: data}
	spanF := tracer.StartSpanFromContext(ctxDB, "Save config in DB")
	_, err = kv.Put(c, nil)

	if err != nil {
		tracer.LogError(spanF, err)
		return nil, err
	}
	spanF.Finish()
	span.LogFields(
		tracer.LogString("ConfigDB", "Successful saving configuration to database."),
	)
	return config, nil

}

func (db *Database) Group(ctx context.Context, group *Group) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "Save group DB")
	defer span.Finish()

	ctxDB := tracer.ContextWithSpan(context.Background(), span)

	kv := db.cli.KV()

	if group.Id == "" {
		group.Id = uuid.New().String()
	}

	spanF := tracer.StartSpanFromContext(ctxDB, "Database kv.Put")

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
			tracer.LogError(span, err)
			return nil, err
		}
		databaseKey := dbkey + uuid.New().String()
		fmt.Println("OVDE JE DBKEY I ON GLASI : " + databaseKey)
		g := &api.KVPair{Key: databaseKey, Value: data}
		_, err = kv.Put(g, nil)
		if err != nil {
			tracer.LogError(spanF, err)
			return nil, err
		}
	}
	spanF.Finish()
	span.LogFields(
		tracer.LogString("GroupDB", "Successful saving group to database."),
	)
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
	spanBase := tracer.StartSpanFromContext(ctxBase, "List method for group")
	data, _, err := kv.List(cKey, nil)

	if err != nil {
		tracer.LogError(spanBase, err)
		return nil, err
	}
	spanBase.Finish()
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

	if len(configs) == 0 {
		tracer.LogError(span, fmt.Errorf("Group doesn't exists!"))
	} else {
		span.LogFields(
			tracer.LogString("database_getConfigs", "Successful reading from database"),
		)
	}
	return group, nil
}

func (ps *Database) GetAllConfigs(ctx context.Context) ([]*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllConfigs")
	defer span.Finish()

	ctxBase := tracer.ContextWithSpan(ctx, span)

	kv := ps.cli.KV()
	spanBase := tracer.StartSpanFromContext(ctxBase, "List configs from database")
	data, _, err := kv.List(allConfigs, nil)

	if err != nil {
		tracer.LogError(spanBase, err)
		return nil, err
	}
	spanBase.Finish()
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
	span.LogFields(
		tracer.LogString("database_getConfigs", "Successful reading from database"),
	)

	return configs, nil
}

func (ps *Database) GetConfigsFromGroup(ctx context.Context, id string, version string, label string) ([]*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConfigsFromGroup DB")
	defer span.Finish()

	ctxF := tracer.ContextWithSpan(ctx, span)

	kv := ps.cli.KV()
	spanF := tracer.StartSpanFromContext(ctxF, "List configs from group with label")
	data, _, err := kv.List(constructKey(id, version, label), nil)

	if err != nil {
		tracer.LogError(spanF, err)
		return nil, err
	}
	spanF.Finish()
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

	span.LogFields(
		tracer.LogString("DB viewConfigsFromGroup", "Successful reading configs with label from database"),
	)
	return configs, nil
}

func (ps *Database) AddConfigsToGroup(ctx context.Context, id string, version string, config Config) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "putConfigsToGroupFromDatabase")
	defer span.Finish()

	ctxBase := tracer.ContextWithSpan(ctx, span)

	group, error := ps.GetGroup(ctx, id, version)
	if error != nil {
		tracer.LogError(span, error)
		return nil, error
	}
	if len(group.Configs) < 1 {
		return nil, errors.New("Group doesn't exists!")
	}

	groupW := Group{}
	groupW.Id = id
	groupW.Version = version
	groupW.Configs = append(groupW.Configs, config)

	spanBase := tracer.StartSpanFromContext(ctxBase, "Put config into group - database")
	_, err := ps.Group(ctx, &groupW)

	if err != nil {
		tracer.LogError(spanBase, err)
		return nil, err
	}
	spanBase.Finish()
	ret, err := ps.GetGroup(ctx, id, version)
	if err != nil {
		tracer.LogError(spanBase, err)
		return nil, err
	}

	return ret, nil
}
