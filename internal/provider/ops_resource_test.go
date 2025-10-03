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

func TestOpsResource_Schema(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/op":
			var op client.Ops
			json.NewDecoder(r.Body).Decode(&op)
			op.ID = "test-id-1"
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(op)
		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/op/id/"):
			op := client.Ops{
				ID:        "test-id-1",
				Name:      "Ops Team Alpha",
				Engineers: []client.Engineer{{ID: "e1"}, {ID: "e2"}},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(op)
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
				Config:       testOpsResourceConfigWithHost(server.URL, "Ops Team Alpha", []string{"e1", "e2"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops_ops.test", "name", "Ops Team Alpha"),
					resource.TestCheckResourceAttr("devops_ops.test", "engineers.#", "2"),
					resource.TestCheckResourceAttr("devops_ops.test", "engineers.0", "e1"),
					resource.TestCheckResourceAttr("devops_ops.test", "engineers.1", "e2"),
					resource.TestCheckResourceAttrSet("devops_ops.test", "id"),
				),
			},
		},
	})
}

func TestOpsResource_Update(t *testing.T) {
	currentOp := client.Ops{
		ID:        "test-id-1",
		Name:      "Ops Team Alpha",
		Engineers: []client.Engineer{{ID: "e1"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/op":
			json.NewDecoder(r.Body).Decode(&currentOp)
			currentOp.ID = "test-id-1"
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(currentOp)
		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/op/id/"):
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(currentOp)
		case r.Method == "PUT":
			json.NewDecoder(r.Body).Decode(&currentOp)
			currentOp.ID = "test-id-1"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(currentOp)
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
				Config:       testOpsResourceConfigWithHost(server.URL, "Ops Team Alpha", []string{"e1"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops_ops.test", "name", "Ops Team Alpha"),
					resource.TestCheckResourceAttr("devops_ops.test", "engineers.#", "1"),
				),
			},
			{
				RefreshState: false,
				Config:       testOpsResourceConfigWithHost(server.URL, "Ops Team Beta", []string{"e3", "e4"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops_ops.test", "name", "Ops Team Beta"),
					resource.TestCheckResourceAttr("devops_ops.test", "engineers.#", "2"),
				),
			},
		},
	})
}

func testOpsResourceConfigWithHost(host string, name string, engineers []string) string {
	engineersJSON, _ := json.Marshal(engineers)
	return `
provider "devops" {
  host = "` + host + `"
}

resource "devops_ops" "test" {
  name      = "` + name + `"
  engineers = ` + string(engineersJSON) + `
}
`
}
