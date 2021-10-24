package handlers

import (
	"log"
	"net/http"
)

type Check struct {
	log *log.Logger
}

func (c Check) Readiness(res http.ResponseWriter, req *http.Request) {
	StatusOK(res)
}
