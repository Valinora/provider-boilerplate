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

func TestDevResource_Schema(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/dev":
			var dev client.Dev
			json.NewDecoder(r.Body).Decode(&dev)
			dev.ID = "test-id-1"
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(dev)
		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/dev/id/"):
			dev := client.Dev{
				ID:        "test-id-1",
				Name:      "Dev Team Alpha",
				Engineers: []client.Engineer{{ID: "e1"}, {ID: "e2"}},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(dev)
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
				Config:       testDevResourceConfigWithHost(server.URL, "Dev Team Alpha", []string{"e1", "e2"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops_dev.test", "name", "Dev Team Alpha"),
					resource.TestCheckResourceAttr("devops_dev.test", "engineers.#", "2"),
					resource.TestCheckResourceAttr("devops_dev.test", "engineers.0", "e1"),
					resource.TestCheckResourceAttr("devops_dev.test", "engineers.1", "e2"),
					resource.TestCheckResourceAttrSet("devops_dev.test", "id"),
				),
			},
		},
	})
}

func TestDevResource_Update(t *testing.T) {
	currentDev := client.Dev{
		ID:        "test-id-1",
		Name:      "Dev Team Alpha",
		Engineers: []client.Engineer{{ID: "e1"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/dev":
			json.NewDecoder(r.Body).Decode(&currentDev)
			currentDev.ID = "test-id-1"
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(currentDev)
		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/dev/id/"):
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(currentDev)
		case r.Method == "PUT":
			json.NewDecoder(r.Body).Decode(&currentDev)
			currentDev.ID = "test-id-1"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(currentDev)
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
				Config:       testDevResourceConfigWithHost(server.URL, "Dev Team Alpha", []string{"e1"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops_dev.test", "name", "Dev Team Alpha"),
					resource.TestCheckResourceAttr("devops_dev.test", "engineers.#", "1"),
				),
			},
			{
				RefreshState: false,
				Config:       testDevResourceConfigWithHost(server.URL, "Dev Team Beta", []string{"e3", "e4"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops_dev.test", "name", "Dev Team Beta"),
					resource.TestCheckResourceAttr("devops_dev.test", "engineers.#", "2"),
				),
			},
		},
	})
}

func testDevResourceConfigWithHost(host string, name string, engineers []string) string {
	engineersJSON, _ := json.Marshal(engineers)
	return `
provider "devops" {
  host = "` + host + `"
}

resource "devops_dev" "test" {
  name      = "` + name + `"
  engineers = ` + string(engineersJSON) + `
}
`
}
