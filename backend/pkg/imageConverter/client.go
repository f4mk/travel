package imageconverter

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	baseURL string
	width   int16
	height  int16
	imgType string
}

type Config struct {
	baseURL string
	width   int16
	height  int16
	imgType string
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
		baseURL: c.baseURL,
		width:   c.width,
		height:  c.height,
		imgType: c.imgType,
	}
}

func (c *Client) Convert(input io.Reader) (io.Reader, error) {
	endpoint := fmt.Sprintf("%s/convert?width=%d&height=%d&type=%s", c.baseURL, c.width, c.height, c.imgType)
	res, err := http.Post(endpoint, "image/*", input)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to convert and resize image: %s", res.Status)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, res.Body)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
