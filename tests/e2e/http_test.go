package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hedgehog-pm/service-tracker/internal/generated/api/common"
)

type Client struct {
	url  string
	http *http.Client
}

func newHTTPClient(t testing.TB) *http.Client {
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	t.Cleanup(client.CloseIdleConnections)

	return client
}

type Response[T any] struct {
	StatusCode int
	Body       []byte
	method     string
	path       string
}

func (r Response[T]) String() string {
	return fmt.Sprintf("%s %s returned %d\n%s", r.method, r.path, r.StatusCode, r.Body)
}

func (r Response[T]) Decode() (T, error) {
	var value T

	if err := json.Unmarshal(r.Body, &value); err != nil {
		return value, fmt.Errorf("decode JSON from %s: %w", r, err)
	}

	return value, nil
}

func (r Response[T]) DecodeError() (common.Error, error) {
	var value common.Error

	if err := json.Unmarshal(r.Body, &value); err != nil {
		return value, fmt.Errorf("decode error from %s: %w", r, err)
	}

	return value, nil
}

func requireResponse[T any](t testing.TB, wantStatus int, response Response[T], err error) T {
	t.Helper()

	require.NoError(t, err)
	require.Equal(t, wantStatus, response.StatusCode, "%s", response)

	value, err := response.Decode()
	require.NoError(t, err)

	return value
}

func request[T any](
	ctx context.Context, client *Client, method, path string, body any,
) (Response[T], error) {
	var content io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return Response[T]{}, fmt.Errorf("encode %s %s: %w", method, path, err)
		}

		content = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, client.url+path, content)
	if err != nil {
		return Response[T]{}, fmt.Errorf("prepare %s %s: %w", method, path, err)
	}

	req.Header.Set("Accept", "application/json")

	if content != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.http.Do(req)
	if err != nil {
		return Response[T]{}, fmt.Errorf("send %s %s: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response[T]{}, fmt.Errorf("read %s %s response: %w", method, path, err)
	}

	return Response[T]{
		StatusCode: resp.StatusCode,
		Body:       data,
		method:     method,
		path:       path,
	}, nil
}
