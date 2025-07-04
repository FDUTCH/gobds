package service

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// Service ...
type Service struct {
	Enabled bool
	Url     string
	Key     string
	Closed  bool

	Client *http.Client
	Log    *slog.Logger
}

const (
	MaxRetries            = 3
	RetryDelay            = 300 * time.Millisecond
	RequestTimeout        = 5 * time.Second
	MaxConcurrentRequests = 5
)

// NewService ...
func NewService(log *slog.Logger, c Config) *Service {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   MaxConcurrentRequests * 2,
		MaxConnsPerHost:       MaxConcurrentRequests * 3,
	}

	return &Service{
		Enabled: c.Enabled,
		Url:     c.URL,
		Key:     c.Key,
		Closed:  false,
		Client: &http.Client{
			Timeout:   RequestTimeout,
			Transport: transport,
		},
		Log: log,
	}
}

// ErrorIsTemporary ...
func ErrorIsTemporary(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}
