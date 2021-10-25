package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

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
	reg, _ := regexp.Compile(`\/\w*\d*`)
	mainPath := reg.FindString(req.URL.String())
	fmt.Println(mainPath)
	StatusNotFound(res)
}

func starshipGroup(a *web.App) {
	a.Get("/starship", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("yo bro this is a get"))
	})
	a.Post("/starship", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("yo bro this is a post"))
	})
}

/**
setup the Mux and direct network requests
*/
func CreateApp(log *log.Logger) http.Handler {
	app := web.NewApp()
	// Checks
	check := Check{log}
	app.HandleFunc("/test", check.Readiness)
	starshipGroup(app)
	// Handlers
	app.Handle("/user", &UserHandler{log})
	app.Handle("/", apiHandler{})
	return app
}
