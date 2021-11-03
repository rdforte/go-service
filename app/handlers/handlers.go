package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rdforte/go-service/foundation/mtang"
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

type testHandler struct {
	log *log.Logger
}

func (t *testHandler) Serve(ctx mtang.Context, res http.ResponseWriter, req *http.Request) {
	t.log.Println("hit the serve")
	res.Write([]byte("test serve"))
}

/**
setup the Mux and direct network requests
*/
func CreateApp(log *log.Logger) http.Handler {
	app := mtang.NewRouter()
	app.Get("/users/:id/spaceship/:type", func(ctx mtang.Context, res http.ResponseWriter, req *http.Request) {
		query := ctx.GetQuery("key")
		userId := ctx.GetParam("id")
		spaceship := ctx.GetParam("type")
		res.Write([]byte(userId + " " + spaceship + " " + query))
	})
	app.HandleFunc("/dog", func(ctx mtang.Context, res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("dog"))
	})

	app.HandleFunc("/", func(ctx mtang.Context, res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("home"))
	})

	app.Handle("/play", &testHandler{log})

	app.NotFound(func(ctx mtang.Context, res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("cant find the route"))
	})

	return app
}
