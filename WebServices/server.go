package main

import (
	"WebServices/database"
	"net/http"

	"github.com/gorilla/mux"
)

type Service struct {
	db database.Database
}

// TO-DO
// func (ts *Service) createConfHandler(w http.ResponseWriter, req *http.Request) {
// 	contentType := req.Header.Get("Content-Type")
// 	mediatype, _, err := mime.ParseMediaType(contentType)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if mediatype != "application/json" {
// 		err := errors.New("expect application/json content-type")
// 		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
// 		return
// 	}

// 	rt, _, err := decodeBody(req.Body, 0)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if rt.Configs[0].Id == "" {
// 		id := createId()
// 		rt.Configs[0].Id = id
// 		ts.configs[id] = append(ts.configs[id], rt.Configs[0])
// 	} else {
// 		id := rt.Configs[0].Id
// 		ts.configs[id] = append(ts.configs[id], rt.Configs[0])
// 	}

// 	renderJSON(w, rt.Configs[0])
// }

// func (ts *Service) createConfGroupHandler(w http.ResponseWriter, req *http.Request) {
// 	contentType := req.Header.Get("Content-Type")
// 	mediatype, _, err := mime.ParseMediaType(contentType)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if mediatype != "application/json" {
// 		err := errors.New("expect application/json content-type")
// 		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
// 		return
// 	}

// 	rt, v, err := decodeBody(req.Body, 1)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	group := rt
// 	group.Version = v
// 	if rt.Id == "" {
// 		idgroup := createId()
// 		group.Id = idgroup
// 	} else {
// 		group.Id = rt.Id
// 	}

// 	ts.versionControl(group)

// 	renderJSON(w, group)
// }

// //test
// func (ts *Service) getConfigsHandler(w http.ResponseWriter, req *http.Request) {
// 	allTasks := []Config{}
// 	for _, v := range ts.configs {
// 		allTasks = append(allTasks, v...)
// 	}

// 	renderJSON(w, allTasks)
// }

// //test
// func (ts *Service) getGroupsHandler(w http.ResponseWriter, req *http.Request) {
// 	allTasks := []Group{}
// 	for _, v := range ts.groups {
// 		allTasks = append(allTasks, v...)
// 	}

// 	renderJSON(w, allTasks)
// }

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

// func (ts *Service) delConfigGroupsHandler(w http.ResponseWriter, req *http.Request) {
// 	id := mux.Vars(req)["id"]
// 	version := mux.Vars(req)["version"]
// 	returnGroup := Group{}
// 	var isExists bool = false
// 	for keyId, v := range ts.groups {
// 		if keyId == id {
// 			for i, g := range v {
// 				if g.Version == version {
// 					isExists = true
// 					returnGroup = g
// 					ts.groups[id] = removeGroup(v, i)
// 					break
// 				}
// 			}
// 		}
// 	}
// 	if !isExists {
// 		renderJSON(w, "Group does not exist")
// 	} else {
// 		renderJSON(w, returnGroup)
// 	}

// }

func (ts *Service) viewConfigHandler(w http.ResponseWriter, req *http.Request) {
	returnConfig := database.Config{}
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	var isExists bool = false
	// for _, v := range db.Get(id) {
	// 	for _, v1 := range v {
	// 		if id == v1.Id && version == v1.Version {
	// 			returnConfig = v1
	// 			isExists = true
	// 			break
	// 		}
	// 	}
	// }
	if !isExists {
		renderJSON(w, "Config does not exist")
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
		renderJSON(w, "Group does not exist")
	} else {
		renderJSON(w, returnGroup)
	}
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
