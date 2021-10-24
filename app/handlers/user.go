package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserHandler struct {
	log *log.Logger
}

func (u *UserHandler) handleGetUser(res http.ResponseWriter, req *http.Request) {
	user := struct {
		FName string `json:"fName"`
	}{
		FName: "Ryan",
	}

	usr, err := json.Marshal(user)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("experienced error marshaling user data"))
	}

	u.log.Println("got the user")

	res.WriteHeader(http.StatusOK)
	res.Write(usr)
}

func (u *UserHandler) notFound(res http.ResponseWriter, req *http.Request) {
	StatusNotFound(res)
}

func (u *UserHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		u.handleGetUser(res, req)
	default:
		u.notFound(res, req)
	}
}
