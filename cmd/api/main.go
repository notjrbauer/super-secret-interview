package main

import (
	"flag"
	"log"
	"os"

	"github.com/notjrbauer/interview/super-secret-interview/worker"
	"github.com/notjrbauer/interview/super-secret-interview/worker/api"
)

func main() {
	var cfgPath string

	flag.StringVar(&cfgPath, "config", "worker.conf", "path to worker config")
	flag.Parse()

	// call with "" to define props on the config object instead of passing a filepath
	// to a config.
	cfg, err := worker.LoadConfig(cfgPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// todo: handle graceful shutdown

	if err := api.ListenAndServe(cfg); err != nil {
		log.Fatalf("server error server %v", err)
	}
}
