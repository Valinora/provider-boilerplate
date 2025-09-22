package provider

import (
	"context"

	"terraform-provider-devops/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &EngineerDataSource{}
)

// NewEngineerDataSource is a helper function to simplify the provider implementation.
func NewEngineerDataSource() datasource.DataSource {
	return &EngineerDataSource{}
}

// EngineerDataSource is the data source implementation.
type EngineerDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *EngineerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer"
}

// Schema defines the schema for the data source.
func (d *EngineerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"engineers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *EngineerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EngineerDataSourceModel

	engineers, err := d.client.GetEngineers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Engineers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, engineer := range engineers {
		engineerState := engineersModel{
			ID:    types.Int64Value(int64(engineer.ID)),
			Name:  types.StringValue(engineer.Name),
			Email: types.StringValue(engineer.Email),
		}

		state.Engineers = append(state.Engineers, engineerState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// EngineerDataSourceModel maps the data source schema data.
type EngineerDataSourceModel struct {
	Engineers []engineersModel `tfsdk:"engineers"`
}

// engineersModel maps engineers schema data.
type engineersModel struct {
	ID    types.Int64  `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

// EngineersInfoModel maps engineers info data
type EngineersInfoModel struct {
	ID types.Int64 `tfsdk:"id"`
}
