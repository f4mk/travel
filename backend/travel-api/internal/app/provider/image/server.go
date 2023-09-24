package image

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	baseURL    string
	httpClient *http.Client
	log        *zerolog.Logger
	sim        chan (struct{})
}

func NewServer(l *zerolog.Logger, b string, t time.Duration, m int16) *Server {
	return &Server{
		baseURL:    b,
		log:        l,
		httpClient: &http.Client{Timeout: t},
		sim:        make(chan struct{}, m),
	}
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
