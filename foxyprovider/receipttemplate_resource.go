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
	_ resource.Resource                = &receiptTemplateResource{}
	_ resource.ResourceWithConfigure   = &receiptTemplateResource{}
	_ resource.ResourceWithImportState = &receiptTemplateResource{}
)

// NewReceiptTemplateResource is a helper function to simplify the provider implementation.
func NewReceiptTemplateResource() resource.Resource {
	return &receiptTemplateResource{}
}

// receiptTemplateResource is the resource implementation.
type receiptTemplateResource struct {
	client *foxyclient.Foxy
}

func (r *receiptTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*foxyclient.Foxy)
}

// Metadata returns the resource type name.
func (r *receiptTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_receipt_template"
}

// Schema defines the schema for the resource.
func (r *receiptTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a receipt template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the receipt template.",
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
func (r *receiptTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan receiptTemplateModel
	// "Get populates the struct passed as `target` with the entire plan."
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	receiptTemplate := foxyclient.ReceiptTemplate{
		Id:          plan.Id.ValueString(),
		Description: plan.Description.ValueString(),
		Content:     plan.Content.ValueString(),
		ContentUrl:  plan.ContentUrl.ValueString(),
	}

	id, err := r.client.ReceiptTemplates.Add(receiptTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating receipt_template",
			"Could not create receipt_template, unexpected error: "+err.Error(),
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

func (r *receiptTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state receiptTemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	receiptTemplate, err := r.client.ReceiptTemplates.Get(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading receipt_template",
			"Could not read receipt_template ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Id = nullableString(receiptTemplate.Id)
	state.Description = nullableString(receiptTemplate.Description)
	state.Content = nullableString(receiptTemplate.Content)
	state.ContentUrl = nullableString(receiptTemplate.ContentUrl)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *receiptTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan receiptTemplateModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	receiptTemplate := foxyclient.ReceiptTemplate{
		Id:          plan.Id.ValueString(),
		Description: plan.Description.ValueString(),
		Content:     plan.Content.ValueString(),
		ContentUrl:  plan.ContentUrl.ValueString(),
	}

	// Update existing receiptTemplate
	_, err := r.client.ReceiptTemplates.Update(plan.Id.ValueString(), receiptTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating receipt_template",
			"Could not update receipt_template, unexpected error: "+err.Error(),
		)
		return
	}

	updatedReceiptTemplate, err := r.client.ReceiptTemplates.Get(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading receipt_template",
			"Could not read receipt_template ID "+plan.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Id = nullableString(updatedReceiptTemplate.Id)
	plan.Description = nullableString(updatedReceiptTemplate.Description)
	plan.Content = nullableString(updatedReceiptTemplate.Content)
	plan.ContentUrl = nullableString(updatedReceiptTemplate.ContentUrl)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *receiptTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state receiptTemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ReceiptTemplates.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting receipt_template",
			"Could not delete receipt_template, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *receiptTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type receiptTemplateModel struct {
	Id types.String `tfsdk:"id"`

	Description types.String `tfsdk:"description"`
	Content     types.String `tfsdk:"content"`
	ContentUrl  types.String `tfsdk:"content_url"`
}
