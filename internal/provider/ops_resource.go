package provider

import (
	"context"
	"fmt"
	"terraform-provider-devops/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &OpsResource{}
	_ resource.ResourceWithConfigure = &OpsResource{}
)

// NewOpsResource is a helper function to simplify the provider implementation.
func NewOpsResource() resource.Resource {
	return &OpsResource{}
}

// OpsResource is the resource implementation.
type OpsResource struct {
	client *client.Client
}

type opsResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Engineers types.List   `tfsdk:"engineers"`
}

// Metadata returns the resource type name.
func (r *OpsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ops"
}

// Schema defines the schema for the resource.
func (r *OpsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"engineers": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *OpsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan devResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	var EngineerIDs []string
	diags = plan.Engineers.ElementsAs(ctx, &EngineerIDs, false)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	engineers := make([]client.Engineer, len(EngineerIDs))

	for i, engiID := range EngineerIDs {
		engineers[i] = client.Engineer{ID: engiID}
	}

	var dev = client.Ops{
		Name:      plan.Name.ValueString(),
		Engineers: engineers,
	}

	createdOps, err := r.client.CreateOps(dev)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Ops",
			"Could not create order, unexpected error: "+err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(createdOps.ID)
	plan.Name = types.StringValue(createdOps.Name)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *OpsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state devResourceModel

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dev, err := r.client.GetOp(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Ops Resource",
			"Could not read Ops: "+state.ID.ValueString()+": "+err.Error(),
		)

		return
	}

	if dev == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(dev.ID)
	state.Name = types.StringValue(dev.Name)

	engineerIDs := make([]string, 0, len(dev.Engineers))

	for _, eng := range dev.Engineers {
		engineerIDs = append(engineerIDs, eng.ID)
	}

	engList, diags2 := types.ListValueFrom(ctx, types.StringType, engineerIDs)

	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		// fugg
		return
	}

	state.Engineers = engList

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *OpsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan devResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var engineerIDs []string
	diags = plan.Engineers.ElementsAs(ctx, &engineerIDs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	engs := make([]client.Engineer, len(engineerIDs))
	for i, engiID := range engineerIDs {
		engs[i] = client.Engineer{ID: engiID}
	}

	dev := client.Ops{
		Name:      plan.Name.ValueString(),
		Engineers: engs,
	}

	_, err := r.client.UpdateOps(plan.ID.ValueString(), dev)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Ops",
			"Could not update dev ID: "+plan.ID.ValueString()+", error: "+err.Error(),
		)

		return
	}

	devResp, err := r.client.GetOp(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Ops",
			"Could not read Ops ID: "+plan.ID.ValueString()+": "+err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(devResp.ID)
	plan.Name = types.StringValue(devResp.Name)

	respEngineerIDs := make([]string, 0, len(devResp.Engineers))
	for _, eng := range devResp.Engineers {
		respEngineerIDs = append(respEngineerIDs, eng.ID)
	}

	engList, diags2 := types.ListValueFrom(ctx, types.StringType, respEngineerIDs)
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Engineers = engList

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *OpsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state devResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOps(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Ops Resource",
			"Could not delete Ops with ID: "+state.ID.ValueString()+" error: "+err.Error(),
		)
		return
	}
}

func (r *OpsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}
