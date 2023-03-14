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
	_ resource.Resource              = &checkoutTemplateResource{}
	_ resource.ResourceWithConfigure = &checkoutTemplateResource{}
)

// NewCheckoutTemplateResource is a helper function to simplify the provider implementation.
func NewCheckoutTemplateResource() resource.Resource {
	return &checkoutTemplateResource{}
}

// checkoutTemplateResource is the resource implementation.
type checkoutTemplateResource struct {
	client *foxyclient.Foxy
}

func (r *checkoutTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*foxyclient.Foxy)
}

// Metadata returns the resource type name.
func (r *checkoutTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_checkout_template"
}

// Schema defines the schema for the resource.
func (r *checkoutTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a checkout template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the checkout template.",
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
func (r *checkoutTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan checkoutTemplateModel
	// "Get populates the struct passed as `target` with the entire plan."
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkoutTemplate := foxyclient.CheckoutTemplate{
		Id:          plan.Id.ValueString(),
		Description: plan.Description.ValueString(),
		Content:     plan.Content.ValueString(),
		ContentUrl:  plan.ContentUrl.ValueString(),
	}

	id, err := r.client.CheckoutTemplates.Add(checkoutTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating checkout_template",
			"Could not create checkout_template, unexpected error: "+err.Error(),
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

func (r *checkoutTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state checkoutTemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkoutTemplate, err := r.client.CheckoutTemplates.Get(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading checkout_template",
			"Could not read checkout_template ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Id = nullableString(checkoutTemplate.Id)
	state.Description = nullableString(checkoutTemplate.Description)
	state.Content = nullableString(checkoutTemplate.Content)
	state.ContentUrl = nullableString(checkoutTemplate.ContentUrl)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *checkoutTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan checkoutTemplateModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkoutTemplate := foxyclient.CheckoutTemplate{
		Id:          plan.Id.ValueString(),
		Description: plan.Description.ValueString(),
		Content:     plan.Content.ValueString(),
		ContentUrl:  plan.ContentUrl.ValueString(),
	}

	// Update existing checkoutTemplate
	_, err := r.client.CheckoutTemplates.Update(plan.Id.ValueString(), checkoutTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating checkout_template",
			"Could not update checkout_template, unexpected error: "+err.Error(),
		)
		return
	}

	updatedCheckoutTemplate, err := r.client.CheckoutTemplates.Get(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading checkout_template",
			"Could not read checkout_template ID "+plan.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Id = nullableString(updatedCheckoutTemplate.Id)
	plan.Description = nullableString(updatedCheckoutTemplate.Description)
	plan.Content = nullableString(updatedCheckoutTemplate.Content)
	plan.ContentUrl = nullableString(updatedCheckoutTemplate.ContentUrl)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *checkoutTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state checkoutTemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CheckoutTemplates.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting checkout_template",
			"Could not delete checkout_template, unexpected error: "+err.Error(),
		)
		return
	}
}

type checkoutTemplateModel struct {
	Id types.String `tfsdk:"id"`

	Description types.String `tfsdk:"description"`
	Content     types.String `tfsdk:"content"`
	ContentUrl  types.String `tfsdk:"content_url"`
}
