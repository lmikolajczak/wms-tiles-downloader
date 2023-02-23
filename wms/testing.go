package wms

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientWithServer(t *testing.T) (*Client, *http.ServeMux, func()) {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("err = %v; want: nil", err)
	}

	return client, mux, func() {
		server.Close()
	}
}

func TestAuthClientWithServer(t *testing.T, credentials string) (*Client, *http.ServeMux, func()) {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client, err := NewClient(server.URL, WithBasicAuth(credentials))
	if err != nil {
		t.Fatalf("err = %v; want: nil", err)
	}

	return client, mux, func() {
		server.Close()
	}
}
