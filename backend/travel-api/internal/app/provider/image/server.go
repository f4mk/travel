package image

import (
	"context"
	"io"
	"sync"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
)

type Server struct {
	log         *zerolog.Logger
	sem         chan (struct{})
	minioClient *minio.Client
	bucketName  string
}

type ServerConfig struct {
	Log        *zerolog.Logger
	WRConns    int16
	Host       string
	AccessKey  string
	SecretKey  string
	BucketName string
}

func NewServer(
	s ServerConfig,
) (*Server, error) {
	minioClient, err := minio.New(s.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(s.AccessKey, s.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &Server{
		log:         s.Log,
		sem:         make(chan struct{}, s.WRConns),
		minioClient: minioClient,
		bucketName:  s.BucketName,
	}, nil
}

func (s *Server) ServeFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
	ctx, span := web.AddSpan(ctx, "provider.image.server.serve-file")
	defer span.End()
	object, err := s.minioClient.GetObject(ctx, s.bucketName, fileID, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (s *Server) SaveFiles(ctx context.Context, filesID []string, streams []io.ReadCloser) error {
	ctx, span := web.AddSpan(ctx, "provider.image.server.save-files")
	defer span.End()
	var wg sync.WaitGroup
	errChan := make(chan error, len(filesID))
	for i, fileID := range filesID {
		wg.Add(1)
		go func(fID string, stream io.ReadCloser) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
			case s.sem <- struct{}{}:
				_, err := s.minioClient.PutObject(ctx, s.bucketName, fID, stream, -1, minio.PutObjectOptions{})
				stream.Close()
				<-s.sem
				if err != nil {
					errChan <- err
				}
			}
		}(fileID, streams[i])
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()

	var firstErr error
	for err := range errChan {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

func (s *Server) TryDeleteFiles(ctx context.Context, filesID []string) error {
	ctx, span := web.AddSpan(ctx, "provider.image.server.try-delete-files")
	defer span.End()
	var errEncountered error
	for _, fileID := range filesID {
		err := s.minioClient.RemoveObject(ctx, s.bucketName, fileID, minio.RemoveObjectOptions{})
		if err != nil {
			if e, ok := err.(minio.ErrorResponse); ok {
				if e.Code == "NoSuchKey" {
					s.log.Warn().Msgf("attempted to delete non-existing file: %s", fileID)
				} else {
					s.log.Err(err).Msgf("minio error encountered on file: %s", fileID)
					errEncountered = err
				}
			} else {
				s.log.Err(err).Msgf("failed to delete file: %s", fileID)
				errEncountered = err
			}
		}
	}

	return errEncountered
}
