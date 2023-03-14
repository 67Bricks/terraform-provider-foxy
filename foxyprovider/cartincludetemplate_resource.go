package foxyprovider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-foxycart/foxyclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &cartIncludeTemplateResource{}
	_ resource.ResourceWithConfigure   = &cartIncludeTemplateResource{}
	_ resource.ResourceWithImportState = &cartIncludeTemplateResource{}
)

// NewCartIncludeTemplateResource is a helper function to simplify the provider implementation.
func NewCartIncludeTemplateResource() resource.Resource {
	return &cartIncludeTemplateResource{}
}

// cartIncludeTemplateResource is the resource implementation.
type cartIncludeTemplateResource struct {
	client *foxyclient.Foxy
}

func (r *cartIncludeTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*foxyclient.Foxy)
}

// Metadata returns the resource type name.
func (r *cartIncludeTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cart_include_template"
}

// Schema defines the schema for the resource.
func (r *cartIncludeTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a cart_include template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the cart_include template.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the template.",
				Required:    true,
			},
			"content": schema.StringAttribute{
				Description: "HTML content of the template.",
				Optional:    true,
			},
			"content_url": schema.StringAttribute{
				Description: "Public URL from which the content can be retrieved",
				Optional:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *cartIncludeTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cartIncludeTemplateModel
	// "Get populates the struct passed as `target` with the entire plan."
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cartIncludeTemplate := foxyclient.CartIncludeTemplate{
		Id:          plan.Id.ValueString(),
		Description: plan.Description.ValueString(),
		Content:     plan.Content.ValueString(),
		ContentUrl:  plan.ContentUrl.ValueString(),
	}

	id, err := r.client.CartIncludeTemplates.Add(cartIncludeTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cart_include_template",
			"Could not create cart_include_template, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(id)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *cartIncludeTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cartIncludeTemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cartIncludeTemplate, err := r.client.CartIncludeTemplates.Get(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading cart_include_template",
			"Could not read cart_include_template ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Id = nullableString(cartIncludeTemplate.Id)
	state.Description = nullableString(cartIncludeTemplate.Description)
	state.Content = nullableString(cartIncludeTemplate.Content)
	state.ContentUrl = nullableString(cartIncludeTemplate.ContentUrl)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *cartIncludeTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan cartIncludeTemplateModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cartIncludeTemplate := foxyclient.CartIncludeTemplate{
		Id:          plan.Id.ValueString(),
		Description: plan.Description.ValueString(),
		Content:     plan.Content.ValueString(),
		ContentUrl:  plan.ContentUrl.ValueString(),
	}

	// Update existing cartIncludeTemplate
	_, err := r.client.CartIncludeTemplates.Update(plan.Id.ValueString(), cartIncludeTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating cart_include_template",
			"Could not update cart_include_template, unexpected error: "+err.Error(),
		)
		return
	}

	updatedCartIncludeTemplate, err := r.client.CartIncludeTemplates.Get(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading cart_include_template",
			"Could not read cart_include_template ID "+plan.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Id = nullableString(updatedCartIncludeTemplate.Id)
	plan.Description = nullableString(updatedCartIncludeTemplate.Description)
	plan.Content = nullableString(updatedCartIncludeTemplate.Content)
	plan.ContentUrl = nullableString(updatedCartIncludeTemplate.ContentUrl)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *cartIncludeTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cartIncludeTemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CartIncludeTemplates.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting cart_include_template",
			"Could not delete cart_include_template, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *cartIncludeTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type cartIncludeTemplateModel struct {
	Id types.String `tfsdk:"id"`

	Description types.String `tfsdk:"description"`
	Content     types.String `tfsdk:"content"`
	ContentUrl  types.String `tfsdk:"content_url"`
}
