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
	_ datasource.DataSource              = &DevOpsDataSource{}
	_ datasource.DataSourceWithConfigure = &DevOpsDataSource{}
)

// NewDevOpsDataSource is a helper function to simplify the provider implementation.
func NewDevOpsDataSource() datasource.DataSource {
	return &DevOpsDataSource{}
}

// DevOpsDataSource is the data source implementation.
type DevOpsDataSource struct {
	client *client.Client
}

func (d *DevOpsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"fuggg",
			"oopsie",
		)

		return
	}

	d.client = c
}

// Metadata returns the data source type name.
func (d *DevOpsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer"
}

// Schema defines the schema for the data source.
func (d *DevOpsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"engineers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
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
func (d *DevOpsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
		engineerState := engineerModel{
			ID:    types.StringValue(engineer.ID),
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
	Engineers []engineerModel `tfsdk:"engineers"`
}

// engineerModel maps engineers schema data.
type engineerModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

// EngineersInfoModel maps engineers info data
type EngineersInfoModel struct {
	ID types.Int64 `tfsdk:"id"`
}
