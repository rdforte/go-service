package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rdforte/go-service/foundation/web"
)

type Status struct {
	Status string `json:"status"`
}

func StatusNotFound(res http.ResponseWriter) {
	res.WriteHeader(http.StatusNotFound)
	json.NewEncoder(res).Encode(Status{"NOT FOUND"})
}

func StatusOK(res http.ResponseWriter) {
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(Status{"OK"})
}

type apiHandler struct{}

func (apiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	StatusNotFound(res)
}

/**
setup the Mux and direct network requests
*/
func CreateApp(log *log.Logger) http.Handler {
	app := web.NewApp()
	// Checks
	check := Check{log}
	app.HandleFunc("/test", check.Readiness)
	// Handlers
	app.Handle("/user", &UserHandler{log})
	app.Handle("/", apiHandler{})
	app.Get("/starship", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("yo bro this is a get"))
	})
	app.Post("/starship", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("yo bro this is a post"))
	})
	app.NotFound(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("this is an error"))
	})
	return app
}
