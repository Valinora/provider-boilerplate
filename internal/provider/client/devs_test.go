package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDevs(t *testing.T) {
	devs := []Dev{
		{ID: "1", Name: "Dev Team 1", Engineers: []Engineer{{ID: "e1", Name: "Alice", Email: "alice@example.com"}}},
		{ID: "2", Name: "Dev Team 2", Engineers: []Engineer{{ID: "e2", Name: "Bob", Email: "bob@example.com"}}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dev" {
			t.Errorf("expected path /dev, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(devs)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.GetDevs()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 devs, got %d", len(result))
	}
	if result[0].ID != "1" {
		t.Errorf("expected first dev ID to be 1, got %s", result[0].ID)
	}
}

func TestGetDev(t *testing.T) {
	dev := Dev{
		ID:        "1",
		Name:      "Dev Team 1",
		Engineers: []Engineer{{ID: "e1", Name: "Alice", Email: "alice@example.com"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dev/id/1" {
			t.Errorf("expected path /dev/id/1, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dev)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.GetDev("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "1" {
		t.Errorf("expected dev ID to be 1, got %s", result.ID)
	}
	if result.Name != "Dev Team 1" {
		t.Errorf("expected dev name to be 'Dev Team 1', got %s", result.Name)
	}
}

func TestCreateDev(t *testing.T) {
	dev := Dev{
		Name:      "New Dev Team",
		Engineers: []Engineer{{ID: "e1"}},
	}

	expectedDev := Dev{
		ID:        "123",
		Name:      "New Dev Team",
		Engineers: []Engineer{{ID: "e1", Name: "Alice", Email: "alice@example.com"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dev" {
			t.Errorf("expected path /dev, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		var receivedDev Dev
		json.NewDecoder(r.Body).Decode(&receivedDev)
		if receivedDev.Name != dev.Name {
			t.Errorf("expected name %s, got %s", dev.Name, receivedDev.Name)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedDev)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.CreateDev(dev)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "123" {
		t.Errorf("expected ID 123, got %s", result.ID)
	}
	if result.Name != "New Dev Team" {
		t.Errorf("expected name 'New Dev Team', got %s", result.Name)
	}
}

func TestUpdateDev(t *testing.T) {
	dev := Dev{
		Name:      "Updated Dev Team",
		Engineers: []Engineer{{ID: "e1"}},
	}

	updatedDev := Dev{
		ID:        "1",
		Name:      "Updated Dev Team",
		Engineers: []Engineer{{ID: "e1", Name: "Alice", Email: "alice@example.com"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dev/1" {
			t.Errorf("expected path /dev/1, got %s", r.URL.Path)
		}
		if r.Method != "PUT" {
			t.Errorf("expected PUT method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedDev)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.UpdateDev("1", dev)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "1" {
		t.Errorf("expected ID 1, got %s", result.ID)
	}
	if result.Name != "Updated Dev Team" {
		t.Errorf("expected name 'Updated Dev Team', got %s", result.Name)
	}
}

func TestDeleteDev(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dev/1" {
			t.Errorf("expected path /dev/1, got %s", r.URL.Path)
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

	err := client.DeleteDev("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteDevError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "unexpected response"}`))
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	err := client.DeleteDev("1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
