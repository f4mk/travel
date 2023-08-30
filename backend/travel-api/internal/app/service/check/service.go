package check

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Service struct {
	build string
	store *sqlx.DB
	log   *zerolog.Logger
}

func NewService(b string, l *zerolog.Logger, s *sqlx.DB) *Service {
	return &Service{
		build: b,
		store: s,
		log:   l,
	}
}

func (cs *Service) Readiness(w http.ResponseWriter, r *http.Request) {

	status := "ok"
	statusCode := http.StatusOK
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	if err := database.StatusCheck(ctx, cs.store); err != nil {
		status = fmt.Sprintf("db not ready: %s", err.Error())
		statusCode = http.StatusInternalServerError
	}
	res := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}
	if err := web.Respond(ctx, w, res, statusCode); err != nil {
		cs.log.Err(err).Msg("readiness: failed to respond:")
	}
	cs.log.Info().Msgf(
		" %s : readiness : remoteAddr: %s path[%s] statusCode [%d]",
		web.GetTraceID(ctx),
		r.RemoteAddr,
		r.URL.Path,
		statusCode,
	)
}

func (cs *Service) Liveness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}
	info := struct {
		Status string `json:"status,omitempty"`
		Build  string `json:"build,omitempty"`
		Host   string `json:"host,omitempty"`
	}{
		Status: "up",
		Build:  cs.build,
		Host:   host,
	}
	if err := web.Respond(ctx, w, info, http.StatusOK); err != nil {
		cs.log.Err(err).Msg("liveness: failed to respond:")
	}
	cs.log.Info().Msgf(
		" %s : liveness  : remoteAddr: %s path[%s] statusCode [%d]",
		web.GetTraceID(ctx),
		r.RemoteAddr,
		r.URL.Path,
		http.StatusOK,
	)
}
