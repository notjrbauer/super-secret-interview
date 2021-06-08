package main

import (
	"flag"
	"log"
	"os"

	"github.com/notjrbauer/interview/super-secret-interview/worker"
	"github.com/notjrbauer/interview/super-secret-interview/worker/cli"
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

	if err := cli.Exec(cfg, os.Args[1:]); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
