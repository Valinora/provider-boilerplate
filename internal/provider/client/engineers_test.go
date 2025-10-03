package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEngineers(t *testing.T) {
	engineers := []Engineer{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/engineers" {
			t.Errorf("expected path /engineers, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(engineers)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.GetEngineers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 engineers, got %d", len(result))
	}
	if result[0].ID != "1" {
		t.Errorf("expected first engineer ID to be 1, got %s", result[0].ID)
	}
}

func TestGetEngineer(t *testing.T) {
	engineer := Engineer{
		ID:    "1",
		Name:  "Alice",
		Email: "alice@example.com",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/engineers/id/1" {
			t.Errorf("expected path /engineers/id/1, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(engineer)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.GetEngineer("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "1" {
		t.Errorf("expected engineer ID to be 1, got %s", result.ID)
	}
	if result.Name != "Alice" {
		t.Errorf("expected engineer name to be 'Alice', got %s", result.Name)
	}
	if result.Email != "alice@example.com" {
		t.Errorf("expected engineer email to be 'alice@example.com', got %s", result.Email)
	}
}

func TestCreateEngineer(t *testing.T) {
	engineer := Engineer{
		Name:  "Charlie",
		Email: "charlie@example.com",
	}

	expectedEngineer := Engineer{
		ID:    "123",
		Name:  "Charlie",
		Email: "charlie@example.com",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/engineers" {
			t.Errorf("expected path /engineers, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		var receivedEngineer Engineer
		json.NewDecoder(r.Body).Decode(&receivedEngineer)
		if receivedEngineer.Name != engineer.Name {
			t.Errorf("expected name %s, got %s", engineer.Name, receivedEngineer.Name)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedEngineer)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.CreateEngineer(engineer)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "123" {
		t.Errorf("expected ID 123, got %s", result.ID)
	}
	if result.Name != "Charlie" {
		t.Errorf("expected name 'Charlie', got %s", result.Name)
	}
}

func TestUpdateEngineer(t *testing.T) {
	engineer := Engineer{
		Name:  "Alice Updated",
		Email: "alice.updated@example.com",
	}

	updatedEngineer := Engineer{
		ID:    "1",
		Name:  "Alice Updated",
		Email: "alice.updated@example.com",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/engineers/1" {
			t.Errorf("expected path /engineers/1, got %s", r.URL.Path)
		}
		if r.Method != "PUT" {
			t.Errorf("expected PUT method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedEngineer)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.UpdateEngineer("1", engineer)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "1" {
		t.Errorf("expected ID 1, got %s", result.ID)
	}
	if result.Name != "Alice Updated" {
		t.Errorf("expected name 'Alice Updated', got %s", result.Name)
	}
}

func TestDeleteEngineer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/engineers/1" {
			t.Errorf("expected path /engineers/1, got %s", r.URL.Path)
		}
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "resource deleted"}`))
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	err := client.DeleteEngineer("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteEngineerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "unexpected response"}`))
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	err := client.DeleteEngineer("1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
