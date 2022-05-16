package main

import (
	"errors"
	"mime"
	"net/http"

	"github.com/gorilla/mux"
)

type Service struct {
	configs map[string][]Config
	groups  map[string][]Group
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

	if rt.Configs[0].Id == "" {
		id := createId()
		rt.Configs[0].Id = id
		ts.configs[id] = append(ts.configs[id], rt.Configs[0])
	} else {
		id := rt.Configs[0].Id
		ts.configs[id] = append(ts.configs[id], rt.Configs[0])
	}

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

	for i := range rt.Configs {
		rt.Configs[i].Id = createId()
	}

	idgroup := createId()
	group := Group{
		Id:      idgroup,
		Configs: rt.Configs,
	}
	ts.groups[idgroup] = append(ts.groups[idgroup], group)
	renderJSON(w, group)
}

//test
func (ts *Service) getConfigsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []Config{}
	for _, v := range ts.configs {
		allTasks = append(allTasks, v...)
	}

	renderJSON(w, allTasks)
}

//test
func (ts *Service) getGroupsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []Group{}
	for _, v := range ts.groups {
		for _, v1 := range v {
			allTasks = append(allTasks, v1)
		}
	}

	renderJSON(w, allTasks)
}

func (ts *Service) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	returnConfig := Config{}
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	var isExists bool = false
	for _, v := range ts.configs {
		for _, v1 := range v {
			if id == v1.Id && version == v1.Version {
				returnConfig = v1
				delete(ts.configs, id)
				isExists = true
				break
			}
		}
	}
	if !isExists {
		renderJSON(w, "Ne postoji ta konfiguracija!")
	} else {
		renderJSON(w, returnConfig)
	}
}

func (ts *Service) delConfigGroupsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	returnGroup := Group{}
	var isExists bool = false
	for _, v := range ts.groups {
		for _, v1 := range v {
			if id == v1.Id && version == v1.Version {
				isExists = true
				returnGroup = v1
				delete(ts.groups, id)
				break
			}
		}
	}
	if !isExists {
		renderJSON(w, "Ne postoji ta konfiguraciona grupa!")
	} else {
		renderJSON(w, returnGroup)
	}

	// if v, ok := ts.groups[id]; ok {
	// 	delete(ts.groups, id)
	// 	renderJSON(w, v)
	// } else {
	// 	err := errors.New("key not found")
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// }
}

func (ts *Service) viewConfigHandler(w http.ResponseWriter, req *http.Request) {
	returnConfig := Config{}
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	var isExists bool = false
	for _, v := range ts.configs {
		for _, v1 := range v {
			if id == v1.Id && version == v1.Version {
				returnConfig = v1
				isExists = true
				break
			}
		}
	}
	if !isExists {
		renderJSON(w, "Ne postoji ta konfiguracija!")
	} else {
		renderJSON(w, returnConfig)
	}

}

func (ts *Service) viewGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	returnGroup := Group{}
	var isExists bool = false
	for _, v := range ts.groups {
		for _, v1 := range v {
			if id == v1.Id && version == v1.Version {
				isExists = true
				returnGroup = v1
				break
			}
		}
	}
	if !isExists {
		renderJSON(w, "Ne postoji ta konfiguraciona grupa!")
	} else {
		renderJSON(w, returnGroup)
	}
}

//TODO change..
// func (ts *Service) updateConfigHandler(w http.ResponseWriter, req *http.Request) {
// 	id := mux.Vars(req)["id"]
// 	version := mux.Vars(req)["version"]
// 	group := ts.groups[id]

// 	if len(group.Configs) == 0 {
// 		renderJSON(w, "Ne mozete dodati novu konfiguraciju!")
// 	} else {
// 		rt, err := decodeBody(req.Body, 0)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		idConfig := createId()
// 		rt[0].Id = idConfig
// 		group.Configs = append(group.Configs, rt[0])
// 		ts.groups[id] = group
// 		renderJSON(w, rt)
// 	}

// }
