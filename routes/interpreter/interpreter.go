// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package interpreter has the implementation of the interpreter api for the server
package interpreter

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cuttle-ai/octopus-service/config"
	"github.com/cuttle-ai/octopus-service/routes"
	"github.com/cuttle-ai/octopus-service/routes/response"
	"github.com/cuttle-ai/octopus/interpreter"
)

//Query struct as input for interperter
type Query struct {
	//NL is the natural language query
	NL string `json:"nl,omitempty"`
}

//Interpret will interpret a given natural language query
func Interpret(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	/*
	 * First we will get the app context
	 * Then we will parse the request payload query
	 * Then we will tokenize the query
	 * Then we will interpret the query
	 * Then we will write the response
	 */
	//getting the app context
	appCtx := ctx.Value(routes.AppContextKey).(*config.AppContext)
	appCtx.Log.Info("Got a request to interpret a query by", appCtx.Session.User.ID)

	//parsing the query
	rq := &Query{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(rq)
	if err != nil {
		//error while decoding the request param
		response.WriteError(w, response.Error{Err: err.Error()}, http.StatusBadRequest)
		return
	}

	//tokenizing the query
	toks, err := interpreter.Tokenize(strconv.Itoa(int(appCtx.Session.User.ID)), []rune(rq.NL))
	if err != nil {
		//error while tokenizing the user query
		response.WriteError(w, response.Error{Err: err.Error()}, http.StatusBadRequest)
		return
	}

	//interpreting the query
	ins, err := interpreter.Interpret(toks)
	if err != nil {
		//error while interpreting the user query
		response.WriteError(w, response.Error{Err: err.Error()}, http.StatusBadRequest)
		return
	}

	//writing the response
	response.Write(w, response.Message{Message: "successfully interpreted the query", Data: ins})
}

func init() {
	routes.AddRoutes(
		routes.Route{
			Version:     "v1",
			Pattern:     "/interpret",
			HandlerFunc: Interpret,
		},
	)
}
