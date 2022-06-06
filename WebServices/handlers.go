package main

import (
	"WebServices/tracer"
	"context"
	"errors"
	"fmt"
	"mime"
	"net/http"

	"github.com/gorilla/mux"
)

func (ts *Service) createConfHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("createConfigHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling create config at %s", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	contentType := req.Header.Get("Content-Type")
	ideKeyID := req.Header.Get("Idempotency-key")

	if ideKeyID == "" {
		tracer.LogError(span, fmt.Errorf("idempotency-key not represented"))
		renderJSON(ctx, w, "Idempotency-key not represented")
		return
	}

	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json content-type")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, _, err := decodeBody(ctx, req.Body, 0)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ideKey, err := ts.db.GetIdempotencyKey(ctx, &ideKeyID)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ideKey == nil {
		_, err := ts.db.IdempotencyKey(ctx, &ideKeyID)
		if err != nil {
			tracer.LogError(span, err)
			renderJSON(ctx, w, "error occured")
			return
		}
		conf, err := ts.db.Config(ctx, &rt.Configs[0])
		if err != nil {
			tracer.LogError(span, err)
			renderJSON(ctx, w, "error occured")
			return
		}

		renderJSON(ctx, w, conf)
	} else {
		renderJSON(ctx, w, "Saved.")
	}

}

func (ts *Service) createConfGroupHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("createConfigGroupHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling create config group at %s", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	contentType := req.Header.Get("Content-Type")
	ideKeyID := req.Header.Get("Idempotency-key")

	if ideKeyID == "" {
		tracer.LogError(span, fmt.Errorf("idempotency-key not represented"))
		renderJSON(ctx, w, "Idempotency-key not represented")
		return
	}

	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json content-type")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, v, err := decodeBody(ctx, req.Body, 1)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ideKey, err := ts.db.GetIdempotencyKey(ctx, &ideKeyID)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ideKey == nil {
		_, err := ts.db.IdempotencyKey(ctx, &ideKeyID)
		if err != nil {
			tracer.LogError(span, err)
			renderJSON(ctx, w, "error occured")
			return
		}

		rt.Version = v
		group, err := ts.db.Group(ctx, &rt)
		if err != nil {
			tracer.LogError(span, err)
			renderJSON(ctx, w, "error occured")
			return
		}

		renderJSON(ctx, w, group)
	} else {
		renderJSON(ctx, w, "Saved.")
	}

}

func (ts *Service) getConfigsHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigsHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(tracer.LogString("handler", fmt.Sprintf("Starting: handling get configs at %s\n", req.URL.Path)))

	ctx := tracer.ContextWithSpan(context.Background(), span)
	allTasks, error := ts.db.GetAllConfigs(ctx)

	if error != nil {
		renderJSON(ctx, w, "ERROR!")
		return
	}
	renderJSON(ctx, w, allTasks)
}

func (ts *Service) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("deleteConfigHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(tracer.LogString("handler", fmt.Sprintf("Starting: handling delete config at %s\n", req.URL.Path)))
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	msg, err := ts.db.DeleteConfig(ctx, id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		tracer.LogError(span, err)
		return
	}
	renderJSON(ctx, w, msg)
}

func (ts *Service) delConfigGroupsHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("deleteGroupHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(tracer.LogString("handler", fmt.Sprintf("Starting: handling delete group at %s\n", req.URL.Path)))
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	msg, err := ts.db.DeleteConfigGroup(ctx, id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		tracer.LogError(span, err)
		return
	}
	renderJSON(ctx, w, msg)
}

func (ts *Service) viewConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigsHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Starting: handling get config at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	returnConfig, error := ts.db.Get(ctx, id, version)
	//to do for base

	if error != nil {
		renderJSON(ctx, w, "Error!")
		return
	}
	if returnConfig.Id == "" {
		renderJSON(ctx, w, "Error!")
		return
	}
	renderJSON(ctx, w, returnConfig)
}

func (ts *Service) viewGroupHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getGroupHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Starting: handling get group at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	returnGroup, error := ts.db.GetGroup(ctx, id, version)

	if error != nil {
		renderJSON(ctx, w, "Error!")
		return
	}

	if len(returnGroup.Configs) == 0 {
		renderJSON(ctx, w, "Group doesn't exists!")
		tracer.LogError(span, fmt.Errorf("Group doesn't exists!"))
		return
	}

	renderJSON(ctx, w, returnGroup)
}

// func (ts *Service) viewGroupLabelHandler(w http.ResponseWriter, req *http.Request) {

// 	id := mux.Vars(req)["id"]
// 	version := mux.Vars(req)["version"]
// 	label := mux.Vars(req)["label"]
// 	list := strings.Split(label, ";")
// 	sort.Strings(list)
// 	sortedLabel := ""
// 	for _, v := range list {
// 		sortedLabel += v + ";"
// 	}
// 	sortedLabel = sortedLabel[:len(sortedLabel)-1]
// 	returnConfigs, error := ts.db.GetConfigsFromGroup(id, version, sortedLabel)

// 	if error != nil {
// 		renderJSON(w, "Error!")
// 	}
// 	renderJSON(w, returnConfigs)
// }

// func (ts *Service) updateConfigHandler(w http.ResponseWriter, req *http.Request) {
// 	id := mux.Vars(req)["id"]
// 	version := mux.Vars(req)["version"]
// 	rt, _, err := decodeBody(req.Body, 0)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	group, err := ts.db.AddConfigsToGroup(id, version, rt.Configs[0])
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	renderJSON(w, group)
// }
