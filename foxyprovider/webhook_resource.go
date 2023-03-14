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
	_ resource.Resource              = &webhookResource{}
	_ resource.ResourceWithConfigure = &webhookResource{}
)

// NewWebhookResource is a helper function to simplify the provider implementation.
func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

// webhookResource is the resource implementation.
type webhookResource struct {
	client *foxyclient.Foxy
}

func (r *webhookResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*foxyclient.Foxy)
}

// Metadata returns the resource type name.
func (r *webhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

// Schema defines the schema for the resource.
func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a webhook.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the webhook.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"format": schema.StringAttribute{
				Description: "Format of the webhook.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the webhook.",
				Required:    true,
			},
			"url": schema.StringAttribute{
				Description: "URL that should be called by the webhook",
				Required:    true,
			},
			"query": schema.StringAttribute{
				Description: "Query used in the webhook.",
				Optional:    true,
			},
			"encryption_key": schema.StringAttribute{
				Description: "Encryption key for the webhook.",
				Optional:    true,
			},
			"event_resource": schema.StringAttribute{
				Description: "Event resource for the webhook.",
				Required:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan webhookModel
	// "Get populates the struct passed as `target` with the entire plan."
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhook := foxyclient.Webhook{
		Id:            plan.Id.ValueString(),
		Format:        plan.Format.ValueString(),
		Version:       2,
		Name:          plan.Name.ValueString(),
		Url:           plan.Url.ValueString(),
		Query:         plan.Query.ValueString(),
		EncryptionKey: plan.EncryptionKey.ValueString(),
		EventResource: plan.EventResource.ValueString(),
	}

	id, err := r.client.Webhooks.Add(webhook)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating webhook",
			"Could not create webhook, unexpected error: "+err.Error(),
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

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state webhookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhook, err := r.client.Webhooks.Get(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading webhook",
			"Could not read webhook ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = nullableString(webhook.Name)
	state.Id = nullableString(webhook.Id)
	state.Format = nullableString(webhook.Format)
	state.Url = nullableString(webhook.Url)
	state.Query = nullableString(webhook.Query)
	state.EncryptionKey = nullableString(webhook.EncryptionKey)
	state.EventResource = nullableString(webhook.EventResource)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan webhookModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhook := foxyclient.Webhook{
		Id:            plan.Id.ValueString(),
		Format:        plan.Format.ValueString(),
		Version:       2,
		Name:          plan.Name.ValueString(),
		Url:           plan.Url.ValueString(),
		Query:         plan.Query.ValueString(),
		EncryptionKey: plan.EncryptionKey.ValueString(),
		EventResource: plan.EventResource.ValueString(),
	}

	// Update existing webhook
	_, err := r.client.Webhooks.Update(plan.Id.ValueString(), webhook)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Webhook",
			"Could not update webhook, unexpected error: "+err.Error(),
		)
		return
	}

	updatedWebhook, err := r.client.Webhooks.Get(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Webhook",
			"Could not read Webhook ID "+plan.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Name = nullableString(updatedWebhook.Name)
	plan.Id = nullableString(updatedWebhook.Id)
	plan.Format = nullableString(updatedWebhook.Format)
	plan.Url = nullableString(updatedWebhook.Url)
	plan.Query = nullableString(updatedWebhook.Query)
	plan.EncryptionKey = nullableString(updatedWebhook.EncryptionKey)
	plan.EventResource = nullableString(updatedWebhook.EventResource)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webhookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Webhooks.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Webhook",
			"Could not delete webhook, unexpected error: "+err.Error(),
		)
		return
	}
}

type webhookModel struct {
	Id types.String `tfsdk:"id"`

	Format        types.String `tfsdk:"format"`
	Name          types.String `tfsdk:"name"`
	Url           types.String `tfsdk:"url"`
	Query         types.String `tfsdk:"query"`
	EncryptionKey types.String `tfsdk:"encryption_key"`
	EventResource types.String `tfsdk:"event_resource"`
}
