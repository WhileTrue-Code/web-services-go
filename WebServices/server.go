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

	conf, err := ts.db.Config(&rt.Configs[0])
	if err != nil {
		renderJSON(w, "error occured")
	}

	renderJSON(w, conf)
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

	rt, v, err := decodeBody(req.Body, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rt.Version = v
	group, err := ts.db.Group(&rt)
	if err != nil {
		renderJSON(w, "error occured")
	}

	renderJSON(w, group)
}

// //test
func (ts *Service) getConfigsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, error := ts.db.GetAllConfigs()
	if error != nil {
		renderJSON(w, "ERROR!")
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
	}
	renderJSON(w, returnConfig)
}

func (ts *Service) viewGroupHandler(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	returnGroup, error := ts.db.GetGroup(id, version)

	if error != nil {
		renderJSON(w, "Error!")
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

// //TODO change..
// func (ts *Service) updateConfigHandler(w http.ResponseWriter, req *http.Request) {
// 	id := mux.Vars(req)["id"]
// 	version := mux.Vars(req)["version"]
// 	groupList := ts.groups[id]
// 	var group Group
// 	var index int

// 	isExist := false
// 	for i, v := range groupList {
// 		if v.Id == id && v.Version == version {
// 			isExist = true
// 			group = v
// 			index = i
// 			break
// 		}
// 	}

// 	if len(groupList) == 0 {
// 		renderJSON(w, "Ne mozete dodati novu konfiguraciju!")
// 	} else {
// 		rt, _, err := decodeBody(req.Body, 0)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		if isExist {
// 			group.Configs = append(group.Configs, rt.Configs[0])
// 			ts.groups[id] = append(ts.groups[id], group)
// 			ts.groups[id] = removeGroup(ts.groups[id], index)
// 			renderJSON(w, ts.groups[id])

// 		} else {
// 			renderJSON(w, "Group does not exist")
// 		}

// 	}

// }

// func (ts *Service) versionControl(g Group) {
// 	if groups, ok := ts.groups[g.Id]; ok {
// 		for k, v := range groups {
// 			if v.Version == g.Version {
// 				groups[k] = g
// 				return
// 			}
// 		}
// 	}

// 	ts.groups[g.Id] = append(ts.groups[g.Id], g)
// }

// func removeGroup(groups []Group, i int) []Group {
// 	return append(groups[:i], groups[i+1:]...)
// }

// func removeConfig(configs []Config, i int) []Config {
// 	return append(configs[:i], configs[i+1:]...)
// }
