// Copyright 2019 Melvin Davis<melvinodsa@gmail.com>. All rights reserved.
// Use of this source code is governed by a Melvin Davis<melvinodsa@gmail.com>
// license that can be found in the LICENSE file.

//Octopus Service Interpreter service for octopus language
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/cuttle-ai/octopus-service/config"
	"github.com/cuttle-ai/octopus-service/log"
	"github.com/cuttle-ai/octopus-service/routes"
)

/*
 * This file contains the main start point of the application
 */

func main() {
	/*
	 * Create a new Server mux
	 * Create a default server
	 * Init the routes
	 * Now listen and serve
	 * Listen to the os signals for exit
	 * Graceful exit when command comes
	 */
	//creating a new server mux
	m := http.NewServeMux()

	//created the default server
	s := &http.Server{
		Addr:           ":" + config.Port,
		Handler:        m,
		ReadTimeout:    config.RequestRTimeout,
		WriteTimeout:   config.ResponseWTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	//inited the routes
	routes.InitRoutes(m)

	//listen and serve to the server
	go func() {
		log.Info("Starting the server at :" + config.Port)
		log.Error(s.ListenAndServe())
	}()

	//listening for syscalls
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, os.Interrupt)
	sig := <-gracefulStop

	//gracefulling exiting when request comes in
	log.Info("Received the interrupt", sig)
	log.Info("Shutting down the server")
	err := s.Shutdown(context.Background())
	if err != nil {
		log.Error("Couldn't end the server gracefully")
	}
}
