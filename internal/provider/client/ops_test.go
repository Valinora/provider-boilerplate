package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOps(t *testing.T) {
	ops := []Ops{
		{ID: "1", Name: "Ops Team 1", Engineers: []Engineer{{ID: "e1", Name: "Alice", Email: "alice@example.com"}}},
		{ID: "2", Name: "Ops Team 2", Engineers: []Engineer{{ID: "e2", Name: "Bob", Email: "bob@example.com"}}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/op" {
			t.Errorf("expected path /op, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ops)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.GetOps()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 ops, got %d", len(result))
	}
	if result[0].ID != "1" {
		t.Errorf("expected first ops ID to be 1, got %s", result[0].ID)
	}
}

func TestGetOp(t *testing.T) {
	op := Ops{
		ID:        "1",
		Name:      "Ops Team 1",
		Engineers: []Engineer{{ID: "e1", Name: "Alice", Email: "alice@example.com"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/op/id/1" {
			t.Errorf("expected path /op/id/1, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(op)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.GetOp("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "1" {
		t.Errorf("expected ops ID to be 1, got %s", result.ID)
	}
	if result.Name != "Ops Team 1" {
		t.Errorf("expected ops name to be 'Ops Team 1', got %s", result.Name)
	}
}

func TestCreateOps(t *testing.T) {
	op := Ops{
		Name:      "New Ops Team",
		Engineers: []Engineer{{ID: "e1"}},
	}

	expectedOp := Ops{
		ID:        "123",
		Name:      "New Ops Team",
		Engineers: []Engineer{{ID: "e1", Name: "Alice", Email: "alice@example.com"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/op" {
			t.Errorf("expected path /op, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		var receivedOp Ops
		json.NewDecoder(r.Body).Decode(&receivedOp)
		if receivedOp.Name != op.Name {
			t.Errorf("expected name %s, got %s", op.Name, receivedOp.Name)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedOp)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.CreateOps(op)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "123" {
		t.Errorf("expected ID 123, got %s", result.ID)
	}
	if result.Name != "New Ops Team" {
		t.Errorf("expected name 'New Ops Team', got %s", result.Name)
	}
}

func TestUpdateOps(t *testing.T) {
	op := Ops{
		Name:      "Updated Ops Team",
		Engineers: []Engineer{{ID: "e1"}},
	}

	updatedOp := Ops{
		ID:        "1",
		Name:      "Updated Ops Team",
		Engineers: []Engineer{{ID: "e1", Name: "Alice", Email: "alice@example.com"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/op/1" {
			t.Errorf("expected path /op/1, got %s", r.URL.Path)
		}
		if r.Method != "PUT" {
			t.Errorf("expected PUT method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedOp)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	result, err := client.UpdateOps("1", op)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "1" {
		t.Errorf("expected ID 1, got %s", result.ID)
	}
	if result.Name != "Updated Ops Team" {
		t.Errorf("expected name 'Updated Ops Team', got %s", result.Name)
	}
}

func TestDeleteOps(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/op/1" {
			t.Errorf("expected path /op/1, got %s", r.URL.Path)
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

	err := client.DeleteOps("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteOpsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "unexpected response"}`))
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	err := client.DeleteOps("1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
