// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package dict has the implementation of the dictionary api for the server
package dict

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cuttle-ai/octopus-service/config"
	"github.com/cuttle-ai/octopus-service/routes"
	"github.com/cuttle-ai/octopus-service/routes/response"
	"github.com/cuttle-ai/octopus/interpreter"
)

//Dict is the dto for the API
type Dict struct {
	//ID of the dicytionary
	ID uint
}

//GetDict will return the dictionary being used
func GetDict(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	/*
	 * First we will get the app context
	 * Then we will get the dictionary
	 * Then we will write the response
	 */
	//getting the app context
	appCtx := ctx.Value(routes.AppContextKey).(*config.AppContext)
	appCtx.Log.Info("Got a request to get fetch the dictionary by", appCtx.Session.User.ID)

	//getting the dictionary
	req := interpreter.DICTRequest{ID: strconv.Itoa(int(appCtx.Session.User.ID)), Type: interpreter.DICTGet, Out: make(chan interpreter.DICTRequest)}
	go interpreter.SendDICTToChannel(interpreter.DICTInputChannel, req)
	res := <-req.Out

	//writing the response
	response.Write(w, response.Message{Message: "successfully fetched the dictionary", Data: res.DICT})
}

//RemoveDict will remove a dictionary from cache
func RemoveDict(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	/*
	 * First we will get the app context
	 * Then we will remove the dict
	 * Then we will write the response
	 */
	//getting the app context
	appCtx := ctx.Value(routes.AppContextKey).(*config.AppContext)
	appCtx.Log.Info("Got a request to remove the dictionary by", appCtx.Session.User.ID)

	//removing the dictionary
	req := interpreter.DICTRequest{ID: strconv.Itoa(int(appCtx.Session.User.ID)), Type: interpreter.DICTRemove}
	go interpreter.SendDICTToChannel(interpreter.DICTInputChannel, req)

	//writing the response
	response.Write(w, response.Message{Message: "successfully removed the dictionary"})
}

//UpdateDict will update teh dict in the cache
func UpdateDict(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	/*
	 * First we will get the app context
	 * Then we will update the dict
	 * Then we will write the response
	 */
	//getting the app context
	appCtx := ctx.Value(routes.AppContextKey).(*config.AppContext)
	appCtx.Log.Info("Got a request to update the dictionary by", appCtx.Session.User.ID)

	//updating the dictionary
	req := interpreter.DICTRequest{ID: strconv.Itoa(int(appCtx.Session.User.ID)), Type: interpreter.DICTUpdate}
	go interpreter.SendDICTToChannel(interpreter.DICTInputChannel, req)

	//writing the response
	response.Write(w, response.Message{Message: "successfully updated the dictionary"})
}

func init() {
	routes.AddRoutes(
		routes.Route{
			Version:     "v1",
			Pattern:     "/dict",
			HandlerFunc: GetDict,
		},
		routes.Route{
			Version:     "v1",
			Pattern:     "/dict/remove",
			HandlerFunc: RemoveDict,
		},
		routes.Route{
			Version:     "v1",
			Pattern:     "/dict/update",
			HandlerFunc: UpdateDict,
		},
	)
}
