// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package datastores has the utilities to access the  datastore services
package datastores

import (
	"sync"

	authConfig "github.com/cuttle-ai/auth-service/config"
	"github.com/cuttle-ai/brain/appctx"
	toolkit "github.com/cuttle-ai/db-toolkit"
	"github.com/cuttle-ai/db-toolkit/datastores/services"
	gDatastores "github.com/cuttle-ai/go-sdk/services/datastores"
	"github.com/cuttle-ai/octopus-service/config"
)

type datastores struct {
	d  map[uint]*services.Service
	ds map[uint]toolkit.Datastore
	m  sync.Mutex
}

var ds datastores

func init() {
	ds = datastores{d: map[uint]*services.Service{}, ds: map[uint]toolkit.Datastore{}}
}

//GetService will return the service whose id has been
func GetService(a config.AppContext, serviceID uint) (toolkit.Datastore, error) {
	/*
	 * We will add the lock
	 * Then will try to get the service from the cache
	 * If not available, will get it from the datastore service
	 * Then will connect to the datastore service
	 * Then will add the same to the cache
	 */
	ds.m.Lock()
	defer ds.m.Unlock()
	_, ok := ds.d[serviceID]
	if ok {
		st, _ := ds.ds[serviceID]
		return st, nil
	}
	s, err := gDatastores.GetDatastore(appctx.WithAccessToken(a, authConfig.MasterAppDetails.AccessToken), serviceID)
	if err != nil {
		a.Log.Error("error while fetching the datastore details from the data integration services")
		return nil, err
	}
	st, err := s.Datastore()
	if err != nil {
		a.Log.Error("error while connecting to the datastore", serviceID)
		return nil, err
	}
	ds.d[serviceID] = s
	ds.ds[serviceID] = st
	return st, nil
}
