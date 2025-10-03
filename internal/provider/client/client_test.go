package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("creates client with default host", func(t *testing.T) {
		client, err := NewClient(nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if client.HostURL != HostURL {
			t.Errorf("expected HostURL to be %s, got %s", HostURL, client.HostURL)
		}
		if client.HTTPClient.Timeout != 10*time.Second {
			t.Errorf("expected timeout to be 10s, got %v", client.HTTPClient.Timeout)
		}
	})

	t.Run("creates client with custom host", func(t *testing.T) {
		customHost := "http://custom:9090"
		client, err := NewClient(&customHost)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if client.HostURL != customHost {
			t.Errorf("expected HostURL to be %s, got %s", customHost, client.HostURL)
		}
	})
}

func TestDoRequest(t *testing.T) {
	t.Run("successful request returns body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": true}`))
		}))
		defer server.Close()

		client := &Client{
			HostURL:    server.URL,
			HTTPClient: &http.Client{},
		}

		req, _ := http.NewRequest("GET", server.URL, nil)
		body, err := client.doRequest(req)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := `{"success": true}`
		if string(body) != expected {
			t.Errorf("expected body %s, got %s", expected, string(body))
		}
	})

	t.Run("accepts 201 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"created": true}`))
		}))
		defer server.Close()

		client := &Client{
			HostURL:    server.URL,
			HTTPClient: &http.Client{},
		}

		req, _ := http.NewRequest("POST", server.URL, nil)
		body, err := client.doRequest(req)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := `{"created": true}`
		if string(body) != expected {
			t.Errorf("expected body %s, got %s", expected, string(body))
		}
	})

	t.Run("returns error on non-success status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "bad request"}`))
		}))
		defer server.Close()

		client := &Client{
			HostURL:    server.URL,
			HTTPClient: &http.Client{},
		}

		req, _ := http.NewRequest("GET", server.URL, nil)
		_, err := client.doRequest(req)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
