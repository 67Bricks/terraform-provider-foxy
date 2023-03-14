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
	_ resource.Resource                = &emailTemplateResource{}
	_ resource.ResourceWithConfigure   = &emailTemplateResource{}
	_ resource.ResourceWithImportState = &emailTemplateResource{}
)

// NewEmailTemplateResource is a helper function to simplify the provider implementation.
func NewEmailTemplateResource() resource.Resource {
	return &emailTemplateResource{}
}

// emailTemplateResource is the resource implementation.
type emailTemplateResource struct {
	client *foxyclient.Foxy
}

func (r *emailTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*foxyclient.Foxy)
}

// Metadata returns the resource type name.
func (r *emailTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_template"
}

// Schema defines the schema for the resource.
func (r *emailTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an email template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the email template.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the template.",
				Required:    true,
			},
			"subject": schema.StringAttribute{
				Description: "Subject of the template.",
				Required:    true,
			},
			"content_html": schema.StringAttribute{
				Description: "HTML content of the template.",
				Optional:    true,
			},
			"content_html_url": schema.StringAttribute{
				Description: "Public URL from which the HTML content can be retrieved",
				Optional:    true,
			},
			"content_text": schema.StringAttribute{
				Description: "Text content of the template.",
				Optional:    true,
			},
			"content_text_url": schema.StringAttribute{
				Description: "Public URL from which the text content can be retrieved",
				Optional:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *emailTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan emailTemplateModel
	// "Get populates the struct passed as `target` with the entire plan."
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	emailTemplate := foxyclient.EmailTemplate{
		Id:             plan.Id.ValueString(),
		Description:    plan.Description.ValueString(),
		Subject:        plan.Subject.ValueString(),
		ContentHtml:    plan.ContentHtml.ValueString(),
		ContentHtmlUrl: plan.ContentHtmlUrl.ValueString(),
		ContentText:    plan.ContentText.ValueString(),
		ContentTextUrl: plan.ContentTextUrl.ValueString(),
	}

	id, err := r.client.EmailTemplates.Add(emailTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating email_template",
			"Could not create email_template, unexpected error: "+err.Error(),
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

func (r *emailTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state emailTemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	emailTemplate, err := r.client.EmailTemplates.Get(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading email_template",
			"Could not read email_template ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Id = nullableString(emailTemplate.Id)
	state.Description = nullableString(emailTemplate.Description)
	state.Subject = nullableString(emailTemplate.Subject)
	state.ContentHtml = nullableString(emailTemplate.ContentHtml)
	state.ContentHtmlUrl = nullableString(emailTemplate.ContentHtmlUrl)
	state.ContentText = nullableString(emailTemplate.ContentText)
	state.ContentTextUrl = nullableString(emailTemplate.ContentTextUrl)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *emailTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan emailTemplateModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	emailTemplate := foxyclient.EmailTemplate{
		Id:             plan.Id.ValueString(),
		Description:    plan.Description.ValueString(),
		Subject:        plan.Subject.ValueString(),
		ContentHtml:    plan.ContentHtml.ValueString(),
		ContentHtmlUrl: plan.ContentHtmlUrl.ValueString(),
		ContentText:    plan.ContentText.ValueString(),
		ContentTextUrl: plan.ContentTextUrl.ValueString(),
	}

	// Update existing emailTemplate
	_, err := r.client.EmailTemplates.Update(plan.Id.ValueString(), emailTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating email_template",
			"Could not update email_template, unexpected error: "+err.Error(),
		)
		return
	}

	updatedEmailTemplate, err := r.client.EmailTemplates.Get(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading email_template",
			"Could not read email_template ID "+plan.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Id = nullableString(updatedEmailTemplate.Id)
	plan.Description = nullableString(updatedEmailTemplate.Description)
	plan.Subject = nullableString(updatedEmailTemplate.Subject)
	plan.ContentHtml = nullableString(updatedEmailTemplate.ContentHtml)
	plan.ContentHtmlUrl = nullableString(updatedEmailTemplate.ContentHtmlUrl)
	plan.ContentText = nullableString(updatedEmailTemplate.ContentText)
	plan.ContentTextUrl = nullableString(updatedEmailTemplate.ContentTextUrl)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *emailTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state emailTemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.EmailTemplates.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting email_template",
			"Could not delete email_template, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *emailTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type emailTemplateModel struct {
	Id types.String `tfsdk:"id"`

	Description    types.String `tfsdk:"description"`
	Subject        types.String `tfsdk:"subject"`
	ContentHtml    types.String `tfsdk:"content_html"`
	ContentHtmlUrl types.String `tfsdk:"content_html_url"`
	ContentText    types.String `tfsdk:"content_text"`
	ContentTextUrl types.String `tfsdk:"content_text_url"`
}
