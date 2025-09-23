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
	_ resource.Resource              = &EngineerResource{}
	_ resource.ResourceWithConfigure = &EngineerResource{}
)

// NewEngineerResource is a helper function to simplify the provider implementation.
func NewEngineerResource() resource.Resource {
	return &EngineerResource{}
}

// EngineerResource is the resource implementation.
type EngineerResource struct {
	client *client.Client
}

type engineerResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

// Metadata returns the resource type name.
func (r *EngineerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer"
}

// Schema defines the schema for the resource.
func (r *EngineerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"email": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *EngineerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan engineerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var engineer = client.Engineer{
		Name:  plan.Name.ValueString(),
		Email: plan.Email.ValueString(),
	}

	createdEngineer, err := r.client.CreateEngineer(engineer)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Engineer",
			"Could not create order, unexpected error: "+err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(createdEngineer.ID)
	plan.Name = types.StringValue(createdEngineer.Name)
	plan.Email = types.StringValue(createdEngineer.Email)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *EngineerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state engineerResourceModel

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	engineer, err := r.client.GetEngineer(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Engineers",
			"Could not read Engineer: "+state.ID.ValueString()+": "+err.Error(),
		)

		return
	}

	state.ID = types.StringValue(engineer.ID)
	state.Name = types.StringValue(engineer.Name)
	state.Email = types.StringValue(engineer.Email)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *EngineerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan engineerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var engineer client.Engineer

	_, err := r.client.UpdateEngineer(plan.ID.ValueString(), engineer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Engineer",
			"Could not update engineer ID: "+plan.ID.ValueString()+", error: "+err.Error(),
		)

		return
	}

	engi, err := r.client.GetEngineer(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Engineer",
			"Could not read Engineer ID: "+plan.ID.ValueString()+": "+err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, engi)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *EngineerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state engineerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEngineer(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Engineer Resource",
			"Could not delete Engineer with ID: "+state.ID.ValueString()+" error: "+err.Error(),
		)
		return
	}
}

func (r *EngineerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
