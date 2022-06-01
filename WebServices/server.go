package main

import (
	"WebServices/database"
	"errors"
	"mime"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

type Service struct {
	db *database.Database
}

//TO-DO
func (ts *Service) createConfHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	ideKeyID := req.Header.Get("Idempotency-key")

	if ideKeyID == "" {
		renderJSON(w, "Idempotency-key not represented")
		return
	}

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

	rt, _, err := decodeBody(req.Body, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ideKey, err := ts.db.GetIdempotencyKey(&ideKeyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ideKey == nil {
		_, err := ts.db.IdempotencyKey(&ideKeyID)
		if err != nil {
			renderJSON(w, "error occured")
			return
		}
		conf, err := ts.db.Config(&rt.Configs[0])
		if err != nil {
			renderJSON(w, "error occured")
			return
		}

		renderJSON(w, conf)
	} else {
		renderJSON(w, "Saved.")
	}

}

func (ts *Service) createConfGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	ideKeyID := req.Header.Get("Idempotency-key")

	if ideKeyID == "" {
		renderJSON(w, "Idempotency-key not represented")
		return
	}

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

	rt, v, err := decodeBody(req.Body, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ideKey, err := ts.db.GetIdempotencyKey(&ideKeyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ideKey == nil {
		_, err := ts.db.IdempotencyKey(&ideKeyID)
		if err != nil {
			renderJSON(w, "error occured")
			return
		}

		rt.Version = v
		group, err := ts.db.Group(&rt)
		if err != nil {
			renderJSON(w, "error occured")
		}

		renderJSON(w, group)
	} else {
		renderJSON(w, "Saved.")
	}

}

// //test
func (ts *Service) getConfigsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, error := ts.db.GetAllConfigs()
	if error != nil {
		renderJSON(w, "ERROR!")
		return
	}
	renderJSON(w, allTasks)
}

//test
func (ts *Service) getGroupsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	label := mux.Vars(req)["label"]
	allTasks, error := ts.db.GetConfigsFromGroup(id, version, label)
	if error != nil {
		renderJSON(w, "ERROR!")
	}
	renderJSON(w, allTasks)
}

func (ts *Service) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	msg, err := ts.db.DeleteConfig(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, msg)
}

func (ts *Service) delConfigGroupsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	msg, err := ts.db.DeleteConfigGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, msg)
}

func (ts *Service) viewConfigHandler(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	returnConfig, error := ts.db.Get(id, version)
	if error != nil {
		renderJSON(w, "Error!")
		return
	}
	if returnConfig.Id == "" {
		renderJSON(w, "Error!")
		return
	}
	renderJSON(w, returnConfig)
}

func (ts *Service) viewGroupHandler(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	returnGroup, error := ts.db.GetGroup(id, version)

	if error != nil {
		renderJSON(w, "Error!")
		return
	}

	if len(returnGroup.Configs) == 0 {
		renderJSON(w, "Group doesn't exists!")
		return
	}

	renderJSON(w, returnGroup)
}

func (ts *Service) viewGroupLabelHandler(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	label := mux.Vars(req)["label"]
	list := strings.Split(label, ";")
	sort.Strings(list)
	sortedLabel := ""
	for _, v := range list {
		sortedLabel += v + ";"
	}
	sortedLabel = sortedLabel[:len(sortedLabel)-1]
	returnConfigs, error := ts.db.GetConfigsFromGroup(id, version, sortedLabel)

	if error != nil {
		renderJSON(w, "Error!")
	}
	renderJSON(w, returnConfigs)
}

func (ts *Service) updateConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	rt, _, err := decodeBody(req.Body, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group, err := ts.db.AddConfigsToGroup(id, version, rt.Configs[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}
