package main

import (
	"errors"
	"mime"
	"net/http"
)

type Service struct {
	configs map[string]Config
	groups  map[string]Group
}

func (ts *Service) createConfHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json content-type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := createId()
	rt[0].Id = id
	ts.configs[id] = rt[0]
	renderJSON(w, rt)
}
func (ts *Service) createConfGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json content-type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range rt {
		rt[i].Id = createId()
	}

	idgroup := createId()
	group := Group{
		Id:      idgroup,
		Configs: rt,
	}
	ts.groups[idgroup] = group
	renderJSON(w, group)
}

//test
func (ts *Service) getConfigsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []Config{}
	for _, v := range ts.configs {
		allTasks = append(allTasks, v)
	}

	renderJSON(w, allTasks)
}

//test
func (ts *Service) getGroupsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []Group{}
	for _, v := range ts.groups {
		allTasks = append(allTasks, v)
	}

	renderJSON(w, allTasks)
}
