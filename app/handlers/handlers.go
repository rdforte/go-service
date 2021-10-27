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

type apiHandler struct {
	app *web.App
}

func starshipGroup(a *web.App) {
	a.Get("/starship", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("yo bro this is a get"))
	})
	a.Post("/starship", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("yo bro this is a post"))
	})
}

func (a *apiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	reg, _ := regexp.Compile(`\/\w*\d*`)
	mainPath := reg.FindString(req.URL.String())
	// const handler = a.app.Router[mainPath]
	fmt.Println(mainPath)
	StatusNotFound(res)
}

/**
setup the Mux and direct network requests
*/
func CreateApp(log *log.Logger) http.Handler {
	app := web.NewApp()
	// Checks
	// check := Check{log}
	// app.HandleFunc("/test", check.Readiness)
	// starshipGroup(app)
	// app.Get("^\\/$", func(res http.ResponseWriter, req *http.Request) {
	// 	res.Write([]byte("this is the home route"))
	// })
	app.Get("/users/:id/spaceship", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("this is the user route"))
	})
	// // Handlers
	// app.Handle("/user", &UserHandler{log})
	// app.Handle("/", &apiHandler{app})
	return app
}
