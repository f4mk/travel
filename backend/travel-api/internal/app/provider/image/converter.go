package image

import (
	"context"
	"io"
	"sync"

	"github.com/rs/zerolog"
)

type ImgConverter interface {
	Convert(ctx context.Context, input io.Reader) (io.ReadCloser, error)
}
type Converter struct {
	client ImgConverter
	log    *zerolog.Logger
	sim    chan (struct{})
}

// TODO: find something appropriate for m
func NewConverter(l *zerolog.Logger, c ImgConverter, m int16) *Converter {
	return &Converter{client: c, log: l, sim: make(chan struct{}, m)}
}

func (c *Converter) Convert(ctx context.Context, images []io.Reader) ([]io.ReadCloser, error) {
	results := make(chan convResult, len(images))
	converted := make([]io.ReadCloser, len(images))
	var firstError error
	var wg sync.WaitGroup
	for i, img := range images {
		wg.Add(1)
		go func(idx int, img io.Reader) {
			defer wg.Done()

			select {
			case c.sim <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() {
				<-c.sim
			}()

			convImg, err := c.client.Convert(ctx, img)
			select {
			case results <- convResult{index: idx, img: convImg, err: err}:
			case <-ctx.Done():
				if convImg != nil {
					convImg.Close()
				}
				return
			}
		}(i, img)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for {
		select {
		case r, ok := <-results:
			if !ok {
				if firstError == nil {
					return converted, nil
				}
				closeConvertedImages(converted)
				return nil, firstError
			}
			if r.err != nil {
				c.log.Err(r.err).Msg("error converting image")
				if firstError == nil {
					firstError = r.err
				}
			} else {
				converted[r.index] = r.img
			}
		case <-ctx.Done():
			closeConvertedImages(converted)
			return nil, ctx.Err()
		}
	}
}

func closeConvertedImages(images []io.ReadCloser) {
	for _, img := range images {
		if img != nil {
			img.Close()
		}
	}
}

type convResult struct {
	index int
	img   io.ReadCloser
	err   error
}
