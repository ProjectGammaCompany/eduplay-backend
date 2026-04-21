package emailClient

import (
	"bytes"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Client struct {
	log     *slog.Logger
	api     *http.Client
	address string
}

func New(log *slog.Logger, address string) *Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		MaxConnsPerHost:       10,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    false,
		DisableKeepAlives:     false,
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	return &Client{
		log:     log,
		api:     client,
		address: address,
	}
}

func (c *Client) SendEmail(data []byte) error {
	op := "SendEmail.Client"

	req, err := http.NewRequest(
		"POST",
		c.address,
		bytes.NewReader(data),
	)
	if err != nil {
		c.log.Error("failed to create http-request", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.api.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.log.Error("Failed to close response body", slog.String("error", err.Error()))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		c.log.Error("failed to send email", slog.String("status code", resp.Status))
		return fmt.Errorf("%s: failed to send email - status code %s", op, resp.Status)
	}

	return nil
}
