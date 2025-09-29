// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"terraform-provider-devops/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var (
	_ provider.Provider = &DevOpsProvider{}
)

// DevOpsProvider defines the provider implementation.
type DevOpsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type DevOpsProviderModel struct {
	HostURL types.String `tfsdk:"host"`
}

func (p *DevOpsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config DevOpsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if config.Endpoint.IsNull() { /* ... */ }

	// Initialize custom API client for data sources and resources
	var endpointPtr *string
	if !config.HostURL.IsNull() && !config.HostURL.IsUnknown() {
		v := config.HostURL.ValueString()
		endpointPtr = &v
	}

	c, err := client.NewClient(endpointPtr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create API client",
			err.Error(),
		)
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c

}

func (p *DevOpsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "devops"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *DevOpsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *DevOpsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEngineerResource,
		NewDevResource,
	}
}

func (p *DevOpsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDevOpsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DevOpsProvider{
			version: version,
		}
	}
}
