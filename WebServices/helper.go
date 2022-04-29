package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func decodeBody(r io.Reader, i int) ([]Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var cfglist []Config
	var cfg Config

	if i == 0 {
		if err := dec.Decode(&cfg); err != nil {
			return nil, err
		}
	} else {
		if err := dec.Decode(&cfglist); err != nil {
			return nil, err
		}
	}

	if len(cfglist) == 0 {
		cfglist = append(cfglist, cfg)
	}

	if len(cfglist) < 1 {
		return nil, fmt.Errorf("configuration list is empty")
	}

	return cfglist, nil
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
