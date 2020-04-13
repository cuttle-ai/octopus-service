// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package db has the utilities required for the db interaction for the service
package db

import (
	"errors"

	"github.com/cuttle-ai/octopus-service/config"
	"github.com/cuttle-ai/octopus-service/datastores"
	"github.com/cuttle-ai/octopus/interpreter"
)

//Exec will execute a query and return the result
func Exec(a config.AppContext, q interpreter.Query) (interface{}, error) {
	/*
	 * Based on the number of tables avaiable in the query, we will execute the same
	 */

	if len(q.Tables) == 0 {
		return nil, errors.New("Couldn't find the table to be queried from")
	}

	if len(q.Tables) == 1 {
		return SingleTableMode(a, q)
	}

	return nil, errors.New("multiple table join query not supported yet")
}

//SingleTableMode execute the given query in a single table mode. So the query is expected not to have any joins or so
func SingleTableMode(a config.AppContext, q interpreter.Query) (interface{}, error) {
	/*
	 * We will convert the query into sql
	 * Then we will get the table from which query has to happen
	 * Then we will get the service corresponding to the table
	 * Will connect to it
	 * Then will execute the query
	 */

	//convert the query
	qs, err := q.ToSQL()
	if err != nil {
		//error while converting the interpreter query to sql
		a.Log.Error("error while converting the interpreter query to sql")
		return nil, err
	}

	//getting the table
	var t *interpreter.TableNode
	for _, v := range q.Tables {
		t = &v
	}
	if t == nil {
		//couldn't find the table
		return nil, errors.New("couldn't find the table in single query mode")
	}

	//getting the datastore service
	ser, err := datastores.GetService(a, t.DatastoreID)
	if err != nil {
		//error while getting the datastore service
		a.Log.Error("error while getting the datastore service")
		return nil, err
	}

	//execute the query
	return ser.Exec(qs.Query, qs.Args...)
}
