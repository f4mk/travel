package imageconverter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	width      int16
	height     int16
	imgType    string
	httpClient *http.Client
}

type Config struct {
	baseURL string
	width   int16
	height  int16
	imgType string
	timeout time.Duration
}

func NewClient(c Config) *Client {
	if c.width == 0 || c.height == 0 {
		c.width = 1280
		c.height = 720
	}
	if c.imgType == "" {
		c.imgType = "webp"
	}
	return &Client{
		baseURL:    c.baseURL,
		width:      c.width,
		height:     c.height,
		imgType:    c.imgType,
		httpClient: &http.Client{Timeout: c.timeout},
	}
}

// make sure to close the returned io.Reader when processed
func (c *Client) Convert(ctx context.Context, input io.Reader) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("%s/convert?width=%d&height=%d&type=%s", c.baseURL, c.width, c.height, c.imgType)
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
