package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Collector struct {
	log    *zerolog.Logger
	client *http.Client
	host   string
}

func New(l *zerolog.Logger, host string) *Collector {
	return &Collector{
		log: l,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		host: host,
	}
}

func (c *Collector) Collect() (map[string]any, error) {
	resp, err := c.client.Get(fmt.Sprintf("http://%s/debug/vars", c.host))
	if err != nil {
		c.log.Err(err).Msg("failed to fetch metrics")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("received non-200 response: %d", resp.StatusCode)
		c.log.Error().Msg(err.Error())
		return nil, err
	}

	raw := make(map[string]any)
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		c.log.Err(err).Msg("failed to decode metrics response")
		return nil, err
	}

	return raw, nil
}
