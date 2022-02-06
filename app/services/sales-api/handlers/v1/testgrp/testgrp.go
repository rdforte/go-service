// Package testgrp has handlers for testing
package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/rdforte/go-service/business/sys/validate"
	"github.com/rdforte/go-service/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints
type Handlers struct {
	Log *zap.SugaredLogger
}

func (h *Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if n := rand.Intn(100); n%2 == 0 {
		// return errors.New("untrusted error")
		return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
		// 	panic("testing panic")
		// return web.NewShutdownError("restart service")
	}

	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	statusCode := http.StatusOK

	return web.Respond(ctx, w, status, statusCode)
}
