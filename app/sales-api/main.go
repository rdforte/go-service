package main

import (
	"log"
	"os"
)

func main() {
	// just write to stdout and it is ops responsibility to direct the logs somewhere.
	log := log.New(os.Stdout, "SALES: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(log); err != nil {
		log.Println("main: error", err)
		os.Exit(1)
	}

}

func run(log *log.Logger) error {
	return nil
}
