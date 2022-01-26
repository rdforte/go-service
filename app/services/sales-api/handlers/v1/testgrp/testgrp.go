package testgrp

import (
	"context"
	"net/http"

	"github.com/rdforte/go-service/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints
type Handlers struct {
	Log *zap.SugaredLogger
}

func (h *Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	statusCode := http.StatusOK

	return web.Respond(ctx, w, status, statusCode)
}
