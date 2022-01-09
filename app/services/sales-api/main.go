package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	// Automatically set GOMAXPROCS to match Linux container CPU quota.
	_ "go.uber.org/automaxprocs"
)

var build = "develop"

func main() {
	// tells us the number of cpu's that can be executing at the same time / number of goroutines that can run in parallel at the same time.
	g := runtime.GOMAXPROCS(0)

	log.Printf("starting service build[%s] CPU[%d]", build, g)
	defer log.Println("service ended")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM) // SIGTERM is what K8s sends for the shutdown

	<-shutdown
	log.Println("stopping service")
}
