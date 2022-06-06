package main

import (
	"WebServices/database"
	"WebServices/tracer"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func decodeBody(ctx context.Context, r io.Reader, i int) (database.Group, string, error) {
	span := tracer.StartSpanFromContext(ctx, "decodeBody")
	defer span.Finish()

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var group database.Group
	var cfg database.Config
	isCfg := false
	if i == 0 {
		isCfg = true
		if err := dec.Decode(&cfg); err != nil {
			tracer.LogError(span, err)
			return database.Group{}, "", err
		}
	} else {
		if err := dec.Decode(&group); err != nil {
			tracer.LogError(span, err)
			return database.Group{}, "", err
		}
	}

	if len(group.Configs) == 0 && isCfg {
		group.Configs = append(group.Configs, cfg)
		return group, "", nil
	}

	fmt.Println(len(group.Configs))
	if len(group.Configs) < 1 {
		tracer.LogError(span, fmt.Errorf("configuration list is empty"))
		return database.Group{}, "", fmt.Errorf("configuration list is empty")
	}

	return group, group.Version, nil
}

func renderJSON(ctx context.Context, w http.ResponseWriter, v interface{}) {
	span := tracer.StartSpanFromContext(ctx, "renderJSON")
	js, err := json.Marshal(v)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId() string {
	return uuid.New().String()
}
