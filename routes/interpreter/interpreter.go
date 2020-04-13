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
	"github.com/cuttle-ai/octopus-service/db"
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
		appCtx.Log.Error("error while parsing the request body", err)
		response.WriteError(w, response.Error{Err: err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//tokenizing the query
	toks, err := interpreter.Tokenize(strconv.Itoa(int(appCtx.Session.User.ID)), []rune(rq.NL))
	if err != nil {
		//error while tokenizing the user query
		appCtx.Log.Error("error while tokenizing the query", err)
		response.WriteError(w, response.Error{Err: "Unable to interpret your query"}, http.StatusInternalServerError)
		return
	}

	//interpreting the query
	ins, err := interpreter.Interpret(toks)
	if err != nil {
		//error while interpreting the user query
		appCtx.Log.Error("error while interpreting the query", err)
		response.WriteError(w, response.Error{Err: "Unable to interpret your query"}, http.StatusInternalServerError)
		return
	}

	//writing the response
	response.Write(w, response.Message{Message: "successfully interpreted the query", Data: ins})
}

//Search will interpret and find the result the given natural language query
func Search(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	/*
	 * First we will get the app context
	 * Then we will parse the request payload query
	 * Then we will tokenize the query
	 * Then we will interpret the query
	 * Then we will execute the query
	 * Then we will write the response
	 */
	//getting the app context
	appCtx := ctx.Value(routes.AppContextKey).(*config.AppContext)
	appCtx.Log.Info("Got a request to search a query by", appCtx.Session.User.ID)

	//parsing the query
	rq := &Query{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(rq)
	if err != nil {
		//error while decoding the request param
		appCtx.Log.Error("error while parsing the request body", err)
		response.WriteError(w, response.Error{Err: err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//tokenizing the query
	toks, err := interpreter.Tokenize(strconv.Itoa(int(appCtx.Session.User.ID)), []rune(rq.NL))
	if err != nil {
		//error while tokenizing the user query
		appCtx.Log.Error("error while tokenizing the query", err)
		response.WriteError(w, response.Error{Err: "Unable to interpret your query"}, http.StatusInternalServerError)
		return
	}

	//interpreting the query
	ins, err := interpreter.Interpret(toks)
	if err != nil {
		//error while interpreting the user query
		appCtx.Log.Error("error while interpreting the query", err)
		response.WriteError(w, response.Error{Err: "Unable to interpret your query"}, http.StatusInternalServerError)
		return
	}

	//executing the query
	ins.Result, err = db.Exec(*appCtx, *ins)
	if err != nil {
		//error while interpreting the user query
		appCtx.Log.Error("error while executing the query", err)
		response.WriteError(w, response.Error{Err: "Unable to interpret your query"}, http.StatusInternalServerError)
		return
	}

	//writing the response
	response.Write(w, response.Message{Message: "successfully search the query", Data: ins})
}

func init() {
	routes.AddRoutes(
		routes.Route{
			Version:     "v1",
			Pattern:     "/interpret",
			HandlerFunc: Interpret,
		},
		routes.Route{
			Version:     "v1",
			Pattern:     "/search",
			HandlerFunc: Search,
		},
	)
}
