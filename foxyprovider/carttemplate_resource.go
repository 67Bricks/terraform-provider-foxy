package foxyprovider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-foxycart/foxyclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &cartTemplateResource{}
	_ resource.ResourceWithConfigure = &cartTemplateResource{}
)

// NewCartTemplateResource is a helper function to simplify the provider implementation.
func NewCartTemplateResource() resource.Resource {
	return &cartTemplateResource{}
}

// cartTemplateResource is the resource implementation.
type cartTemplateResource struct {
	client *foxyclient.Foxy
}

func (r *cartTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*foxyclient.Foxy)
}

// Metadata returns the resource type name.
func (r *cartTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cart_template"
}

// Schema defines the schema for the resource.
func (r *cartTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a cart template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the cart template.",
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
				Required:    true,
			},
			"content_url": schema.StringAttribute{
				Description: "Public URL from which the content can be retrieved",
				Required:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *cartTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cartTemplateModel
	// "Get populates the struct passed as `target` with the entire plan."
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cartTemplate := foxyclient.CartTemplate{
		Id:          plan.Id.ValueString(),
		Description: plan.Description.ValueString(),
		Content:     plan.Content.ValueString(),
		ContentUrl:  plan.ContentUrl.ValueString(),
	}

	id, err := r.client.CartTemplates.Add(cartTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cart_template",
			"Could not create cart_template, unexpected error: "+err.Error(),
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

func (r *cartTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cartTemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cartTemplate, err := r.client.CartTemplates.Get(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading cart_template",
			"Could not read cart_template ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Id = nullableString(cartTemplate.Id)
	state.Description = nullableString(cartTemplate.Description)
	state.Content = nullableString(cartTemplate.Content)
	state.ContentUrl = nullableString(cartTemplate.ContentUrl)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *cartTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan cartTemplateModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cartTemplate := foxyclient.CartTemplate{
		Id:          plan.Id.ValueString(),
		Description: plan.Description.ValueString(),
		Content:     plan.Content.ValueString(),
		ContentUrl:  plan.ContentUrl.ValueString(),
	}

	// Update existing cartTemplate
	_, err := r.client.CartTemplates.Update(plan.Id.ValueString(), cartTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating cart_template",
			"Could not update cart_template, unexpected error: "+err.Error(),
		)
		return
	}

	updatedCartTemplate, err := r.client.CartTemplates.Get(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading cart_template",
			"Could not read cart_template ID "+plan.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Id = nullableString(updatedCartTemplate.Id)
	plan.Description = nullableString(updatedCartTemplate.Description)
	plan.Content = nullableString(updatedCartTemplate.Content)
	plan.ContentUrl = nullableString(updatedCartTemplate.ContentUrl)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *cartTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cartTemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CartTemplates.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting cart_template",
			"Could not delete cart_template, unexpected error: "+err.Error(),
		)
		return
	}
}

type cartTemplateModel struct {
	Id types.String `tfsdk:"id"`

	Description types.String `tfsdk:"description"`
	Content     types.String `tfsdk:"content"`
	ContentUrl  types.String `tfsdk:"content_url"`
}
