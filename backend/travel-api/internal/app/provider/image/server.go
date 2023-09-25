package image

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
)

type Server struct {
	log         *zerolog.Logger
	sim         chan (struct{})
	minioClient *minio.Client
	bucketName  string
}

func NewServer(
	l *zerolog.Logger,
	m int16,
	endpoint,
	accessKey,
	secretKey,
	bucketName string,
) (*Server, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &Server{
		log:         l,
		sim:         make(chan struct{}, m),
		minioClient: minioClient,
		bucketName:  bucketName,
	}, nil
}

func (s *Server) ServeFile(ctx context.Context, fileID string) ([]byte, error) {
	object, err := s.minioClient.GetObject(ctx, s.bucketName, fileID, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Server) SaveFiles(ctx context.Context, filesID []string, streams []io.ReadCloser) error {
	for i, fileID := range filesID {
		_, err := s.minioClient.PutObject(ctx, s.bucketName, fileID, streams[i], -1, minio.PutObjectOptions{})
		streams[i].Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) TryDeleteFiles(ctx context.Context, filesID []string) error {
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
