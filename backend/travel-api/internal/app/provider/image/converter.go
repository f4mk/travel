package image

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

type ImgConverter interface {
	Convert(ctx context.Context, input io.Reader) (io.Reader, error)
}
type Converter struct {
	client ImgConverter
	log    *zerolog.Logger
}

func NewConverter(l *zerolog.Logger, c ImgConverter) *Converter {
	return &Converter{client: c, log: l}
}

func (c Converter) Convert(ctx context.Context, images []io.Reader) ([]io.Reader, error) {
	results := make(chan convResult, len(images))

	for i, img := range images {
		go func(idx int, img io.Reader) {
			convImg, err := c.client.Convert(ctx, img)
			results <- convResult{index: idx, img: convImg, err: err}
		}(i, img)
	}

	converted := make([]io.Reader, len(images))
	var firstError error

	for i := 0; i < len(images); i++ {
		select {
		case r := <-results:
			if r.err != nil && firstError == nil {
				firstError = r.err
			}
			converted[r.index] = r.img
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if firstError != nil {
		return nil, firstError
	}

	return converted, nil
}

type convResult struct {
	index int
	img   io.Reader
	err   error
}
