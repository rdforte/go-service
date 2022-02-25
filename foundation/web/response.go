package web

import (
	"context"
	"encoding/json"
	"net/http"
)

type OK struct {
	Status string `json:"status"`
}

// Respond converts a Go value to JSON andd send it to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {

	// Set the status code for the request logger middleware.
	// If this fails we still want to process the response.
	SetStatusCode(ctx, statusCode)

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling is successful.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

// RespondOk is the default for responding with a http status ok and should be used for all routes
// where all we want to do is notify the client that the request was successful.
func RespondOk(ctx context.Context, w http.ResponseWriter) error {
	status := OK{
		Status: "OK",
	}

	statusCode := http.StatusOK

	return Respond(ctx, w, status, statusCode)
}
