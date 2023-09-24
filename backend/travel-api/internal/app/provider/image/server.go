package image

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	log        *zerolog.Logger
	sim        chan (struct{})
	baseURL    string
	httpClient *http.Client
}

func NewServer(l *zerolog.Logger, t time.Duration, m int16) *Server {
	return &Server{httpClient: &http.Client{Timeout: t}, log: l, sim: make(chan struct{}, m)}
}

func (s *Server) ServeFile(ctx context.Context, fileID string) ([]byte, error) {
	return []byte{}, nil
}
func (s *Server) SaveFiles(ctx context.Context, filesID []string, streams []io.ReadCloser) error {
	return nil
}
func (s *Server) DeleteFiles(ctx context.Context, filesID []string) error {
	return nil
}
