package imageconverter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	host       string
	width      int16
	height     int16
	imgType    string
	httpClient *http.Client
}

type Config struct {
	Host    string
	Width   int16
	Height  int16
	ImgType string
	Timeout time.Duration
}

func NewClient(c Config) *Client {
	if c.Width == 0 || c.Height == 0 {
		c.Width = 1280
		c.Height = 720
	}
	if c.ImgType == "" {
		c.ImgType = "webp"
	}
	return &Client{
		host:       c.Host,
		width:      c.Width,
		height:     c.Height,
		imgType:    c.ImgType,
		httpClient: &http.Client{Timeout: c.Timeout},
	}
}

// make sure to close the returned io.Reader when processed
func (c *Client) Convert(ctx context.Context, input io.Reader) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("%s/convert?width=%d&height=%d&type=%s", c.host, c.width, c.height, c.imgType)
	req, err := http.NewRequest(http.MethodPost, endpoint, input)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "image/*")
	req = req.WithContext(ctx)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("failed to convert and resize image: %s", res.Status)
	}

	return res.Body, nil
}
