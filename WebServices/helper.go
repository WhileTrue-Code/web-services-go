package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func decodeBody(r io.Reader, i int) (Group, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var group Group
	var cfg Config

	if i == 0 {
		if err := dec.Decode(&cfg); err != nil {
			return group, err
		}
	} else {
		if err := dec.Decode(&group); err != nil {
			return group, err
		}
	}

	if len(group.Configs) == 0 {
		group.Configs = append(group.Configs, cfg)
	}

	if len(group.Configs) < 1 {
		return group, fmt.Errorf("configuration list is empty")
	}

	return group, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId() string {
	return uuid.New().String()
}
