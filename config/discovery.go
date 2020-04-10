// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	aConfig "github.com/cuttle-ai/auth-service/config"
	aLog "github.com/cuttle-ai/auth-service/log"
	"github.com/hashicorp/consul/api"
)

/*
 * This file contains the discovery service init
 */

//OctopusServiceID is the service id to be used with the discovery service
var OctopusServiceID = "Brain-Octopus-Service"

//OctopusServiceRPCID is the rpc service id to be used with the discovery service
var OctopusServiceRPCID = "Brain-Octopus-Service-RPC"

func init() {
	/*
	 * We will communicate with the consul client
	 * Will prepare the service instance for the http and rpc service
	 * Then will register the application with consul
	 * Then we will register the rpc service with the consul agent
	 */
	//Registering the db with the discovery api
	// Get a new client
	log.Println("Going to register with the discovery service")
	dConfig := api.DefaultConfig()
	dConfig.Address = DiscoveryURL
	dConfig.Token = DiscoveryToken
	client, err := api.NewClient(dConfig)
	if err != nil {
		log.Fatal("Error while initing the discovery service client", err.Error())
		return
	}

	//service instances for the http service
	log.Println("Connected with discovery service")
	appInstance := &api.AgentServiceRegistration{
		Name:    OctopusServiceID,
		Port:    IntPort,
		Address: ServiceDomain,
		Tags:    []string{OctopusServiceID},
	}

	//registering the service with the agent
	log.Println("Going to register with the discovery service")
	err = client.Agent().ServiceRegister(appInstance)
	if err != nil {
		log.Fatal("Error while registering with the discovery agent", err.Error())
	}

	//service instance for rpc service
	rpcInstance := &api.AgentServiceRegistration{
		Name:    OctopusServiceRPCID,
		Port:    RPCIntPort,
		Address: ServiceDomain,
		Tags:    []string{OctopusServiceRPCID},
		Meta:    map[string]string{"RPCService": "yes"},
	}
	log.Println("Going to register the rpc service with the discovery service")
	err = client.Agent().ServiceRegister(rpcInstance)
	if err != nil {
		log.Fatal("Error while registering the rpc service with the discovery agent", err.Error())
	}

	log.Println("Successfully registered with the discovery service")
}

func init() {
	//we will init the auth service
	l := aLog.NewLogger(0)
	err := aConfig.InitAuthState(l)
	if err != nil {
		log.Fatal("Error while registering with the auth service agent", err.Error())
	}
}

//StartRPC service will start the rpc service. It helps the services to communicate between each other
func StartRPC() {
	/*
	 * Will register the user auth rpc with rpc package
	 * We will listen to the http with rpc of auth module
	 * Then we will start listening to the rpc port
	 */
	//Registering the auth model with the rpc package
	rpc.Register(new(aConfig.RPCAuth))

	//registering the handler with http
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+RPCPort)
	if e != nil {
		log.Fatal("Error while listening to the rpc port", e.Error())
	}
	go http.Serve(l, nil)
}
