package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-devops/internal/provider/client"
)

func TestEngineerResource_Schema(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/engineers":
			var engineer client.Engineer
			json.NewDecoder(r.Body).Decode(&engineer)
			engineer.ID = "test-id-1"
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(engineer)
		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/engineers/id/"):
			engineer := client.Engineer{
				ID:    "test-id-1",
				Name:  "Alice",
				Email: "alice@example.com",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(engineer)
		case r.Method == "DELETE":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "resource deleted"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testEngineerResourceConfigWithHost(server.URL, "Alice", "alice@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops_engineer.test", "name", "Alice"),
					resource.TestCheckResourceAttr("devops_engineer.test", "email", "alice@example.com"),
					resource.TestCheckResourceAttrSet("devops_engineer.test", "id"),
				),
			},
		},
	})
}

func TestEngineerResource_Update(t *testing.T) {
	currentEngineer := client.Engineer{
		ID:    "test-id-1",
		Name:  "Alice",
		Email: "alice@example.com",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/engineers":
			json.NewDecoder(r.Body).Decode(&currentEngineer)
			currentEngineer.ID = "test-id-1"
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(currentEngineer)
		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/engineers/id/"):
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(currentEngineer)
		case r.Method == "PUT":
			json.NewDecoder(r.Body).Decode(&currentEngineer)
			currentEngineer.ID = "test-id-1"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(currentEngineer)
		case r.Method == "DELETE":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "resource deleted"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				RefreshState: false,
				Config:       testEngineerResourceConfigWithHost(server.URL, "Alice", "alice@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops_engineer.test", "name", "Alice"),
					resource.TestCheckResourceAttr("devops_engineer.test", "email", "alice@example.com"),
				),
			},
			{
				RefreshState: false,
				Config:       testEngineerResourceConfigWithHost(server.URL, "Bob", "bob@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops_engineer.test", "name", "Bob"),
					resource.TestCheckResourceAttr("devops_engineer.test", "email", "bob@example.com"),
				),
			},
		},
	})
}

func testEngineerResourceConfigWithHost(host string, name string, email string) string {
	return `
provider "devops" {
  host = "` + host + `"
}

resource "devops_engineer" "test" {
  name  = "` + name + `"
  email = "` + email + `"
}
`
}
