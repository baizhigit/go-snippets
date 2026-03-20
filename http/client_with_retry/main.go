package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func doRequestWithRetry(ctx context.Context, url string, maxRetries int) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var resp *http.Response
	var err error

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err = client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		defer resp.Body.Close()

		if i < maxRetries-1 {
			backoff := time.Duration(1<<uint(i)) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %d", err)
}
