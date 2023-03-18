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
	_ resource.Resource                = &storeInfoResource{}
	_ resource.ResourceWithConfigure   = &storeInfoResource{}
	_ resource.ResourceWithImportState = &storeInfoResource{}
)

// NewStoreInfoResource is a helper function to simplify the provider implementation.
func NewStoreInfoResource() resource.Resource {
	return &storeInfoResource{}
}

// storeInfoResource is the resource implementation.
type storeInfoResource struct {
	client *foxyclient.Foxy
}

func (r *storeInfoResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*foxyclient.Foxy)
}

// Metadata returns the resource type name.
func (r *storeInfoResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_store_info"
}

// Schema defines the schema for the resource.
func (r *storeInfoResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages store info.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the store info",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"store_name": schema.StringAttribute{
				Required: true,
			},
			"store_domain": schema.StringAttribute{
				Required: true,
			},
			"use_remote_domain": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			"store_url": schema.StringAttribute{
				Required: true,
			},
			"receipt_continue_url": schema.StringAttribute{
				Optional: true,
			},
			"store_email": schema.StringAttribute{
				Required: true,
			},
			"from_email": schema.StringAttribute{
				Optional: true,
			},
			"use_email_dns": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			"bcc_on_receipt_email": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(true),
				},
				Computed: true,
			},
			"smtp_config": schema.StringAttribute{
				Optional: true,
			},
			"postal_code": schema.StringAttribute{
				Required: true,
			},
			"region": schema.StringAttribute{
				Required: true,
			},
			"country": schema.StringAttribute{
				Required: true,
			},
			"locale_code": schema.StringAttribute{
				Required: true,
			},
			"timezone": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringDefault("America/Los_Angeles"),
				},
				Computed: true,
			},
			"hide_currency_symbol": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			"hide_decimal_characters": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			"use_international_currency_symbol": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			"language": schema.StringAttribute{
				Optional: true,
			},
			"logo_url": schema.StringAttribute{
				Optional: true,
			},
			"checkout_type": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringDefault("default_account"),
				},
				Computed: true,
			},
			"use_webhook": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			"webhook_url": schema.StringAttribute{
				Optional: true,
			},
			"webhook_key": schema.StringAttribute{
				Optional: true,
			},
			"use_cart_validation": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			"use_single_sign_on": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			"single_sign_on_url": schema.StringAttribute{
				Optional: true,
			},
			"customer_password_hash_type": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringDefault("phpass"),
				},
				Computed: true,
			},
			"customer_password_hash_config": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringDefault("8"),
				},
				Computed: true,
			},
			"features_multiship": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			"products_require_expires_property": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			"app_session_time": schema.Int64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64Default(604800),
				},
				Computed: true,
			},
			"shipping_address_type": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringDefault("residential"),
				},
				Computed: true,
			},
			"require_signed_shipping_rates": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(true),
				},
				Computed: true,
			},
			"unified_order_entry_password": schema.StringAttribute{
				Optional: true,
			},
			//"custom_display_id_config": schema.StringAttribute{
			//	Optional:    true,
			//},
			"is_maintenance_mode": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			"is_active": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolDefault(false),
				},
				Computed: true,
			},
			//"first_payment_date": schema.StringAttribute{
			//	Optional:    true,
			//},
			//"features": schema.StringAttribute{
			//	Optional:    true,
			//},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *storeInfoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

}

func (r *storeInfoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state storeInfoModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	storeInfo, err := r.client.StoreInfo.Get()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading storeInfo",
			"Could not read storeInfo : "+err.Error(),
		)
		return
	}

	state.StoreName = nullableString(storeInfo.StoreName)
	state.StoreDomain = nullableString(storeInfo.StoreDomain)
	state.UseRemoteDomain = types.BoolValue(storeInfo.UseRemoteDomain)
	state.StoreUrl = nullableString(storeInfo.StoreUrl)
	state.ReceiptContinueUrl = nullableString(storeInfo.ReceiptContinueUrl)
	state.StoreEmail = nullableString(storeInfo.StoreEmail)
	state.FromEmail = nullableString(storeInfo.FromEmail)
	state.UseEmailDns = types.BoolValue(storeInfo.UseEmailDns)
	state.BccOnReceiptEmail = types.BoolValue(storeInfo.BccOnReceiptEmail)
	state.SmtpConfig = nullableString(storeInfo.SmtpConfig)
	state.PostalCode = nullableString(storeInfo.PostalCode)
	state.Region = nullableString(storeInfo.Region)
	state.Country = nullableString(storeInfo.Country)
	state.LocaleCode = nullableString(storeInfo.LocaleCode)
	state.Timezone = nullableString(storeInfo.Timezone)
	state.HideCurrencySymbol = types.BoolValue(storeInfo.HideCurrencySymbol)
	state.HideDecimalCharacters = types.BoolValue(storeInfo.HideDecimalCharacters)
	state.UseInternationalCurrencySymbol = types.BoolValue(storeInfo.UseInternationalCurrencySymbol)
	state.Language = nullableString(storeInfo.Language)
	state.LogoUrl = nullableString(storeInfo.LogoUrl)
	state.CheckoutType = nullableString(storeInfo.CheckoutType)
	state.UseWebhook = types.BoolValue(storeInfo.UseWebhook)
	state.WebhookUrl = nullableString(storeInfo.WebhookUrl)
	state.WebhookKey = nullableString(storeInfo.WebhookKey)
	state.UseCartValidation = types.BoolValue(storeInfo.UseCartValidation)
	state.UseSingleSignOn = types.BoolValue(storeInfo.UseSingleSignOn)
	state.SingleSignOnUrl = nullableString(storeInfo.SingleSignOnUrl)
	state.CustomerPasswordHashType = nullableString(storeInfo.CustomerPasswordHashType)
	state.CustomerPasswordHashConfig = nullableString(storeInfo.CustomerPasswordHashConfig)
	state.FeaturesMultiship = types.BoolValue(storeInfo.FeaturesMultiship)
	state.ProductsRequireExpiresProperty = types.BoolValue(storeInfo.ProductsRequireExpiresProperty)
	state.AppSessionTime = types.Int64Value(int64(storeInfo.AppSessionTime))
	state.ShippingAddressType = nullableString(storeInfo.ShippingAddressType)
	state.RequireSignedShippingRates = types.BoolValue(storeInfo.RequireSignedShippingRates)
	state.UnifiedOrderEntryPassword = nullableString(storeInfo.UnifiedOrderEntryPassword)
	state.IsMaintenanceMode = types.BoolValue(storeInfo.IsMaintenanceMode)
	state.IsActive = types.BoolValue(storeInfo.IsActive)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *storeInfoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan storeInfoModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	storeInfo := foxyclient.StoreInfo{
		StoreName:                      plan.StoreName.ValueString(),
		StoreDomain:                    plan.StoreDomain.ValueString(),
		UseRemoteDomain:                plan.UseRemoteDomain.ValueBool(),
		StoreUrl:                       plan.StoreUrl.ValueString(),
		ReceiptContinueUrl:             plan.ReceiptContinueUrl.ValueString(),
		StoreEmail:                     plan.StoreEmail.ValueString(),
		FromEmail:                      plan.FromEmail.ValueString(),
		UseEmailDns:                    plan.UseEmailDns.ValueBool(),
		BccOnReceiptEmail:              plan.BccOnReceiptEmail.ValueBool(),
		SmtpConfig:                     plan.SmtpConfig.ValueString(),
		PostalCode:                     plan.PostalCode.ValueString(),
		Region:                         plan.Region.ValueString(),
		Country:                        plan.Country.ValueString(),
		LocaleCode:                     plan.LocaleCode.ValueString(),
		Timezone:                       plan.Timezone.ValueString(),
		HideCurrencySymbol:             plan.HideCurrencySymbol.ValueBool(),
		HideDecimalCharacters:          plan.HideDecimalCharacters.ValueBool(),
		UseInternationalCurrencySymbol: plan.UseInternationalCurrencySymbol.ValueBool(),
		Language:                       plan.Language.ValueString(),
		LogoUrl:                        plan.LogoUrl.ValueString(),
		CheckoutType:                   plan.CheckoutType.ValueString(),
		UseWebhook:                     plan.UseWebhook.ValueBool(),
		WebhookUrl:                     plan.WebhookUrl.ValueString(),
		WebhookKey:                     plan.WebhookKey.ValueString(),
		UseCartValidation:              plan.UseCartValidation.ValueBool(),
		UseSingleSignOn:                plan.UseSingleSignOn.ValueBool(),
		SingleSignOnUrl:                plan.SingleSignOnUrl.ValueString(),
		CustomerPasswordHashType:       plan.CustomerPasswordHashType.ValueString(),
		CustomerPasswordHashConfig:     plan.CustomerPasswordHashConfig.ValueString(),
		FeaturesMultiship:              plan.FeaturesMultiship.ValueBool(),
		ProductsRequireExpiresProperty: plan.ProductsRequireExpiresProperty.ValueBool(),
		AppSessionTime:                 int(plan.AppSessionTime.ValueInt64()),
		ShippingAddressType:            plan.ShippingAddressType.ValueString(),
		RequireSignedShippingRates:     plan.RequireSignedShippingRates.ValueBool(),
		UnifiedOrderEntryPassword:      plan.UnifiedOrderEntryPassword.ValueString(),
		IsMaintenanceMode:              plan.IsMaintenanceMode.ValueBool(),
		IsActive:                       plan.IsActive.ValueBool(),
	}

	// Update existing storeInfo
	_, err := r.client.StoreInfo.Update(storeInfo)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating StoreInfo",
			"Could not update storeInfo, unexpected error: "+err.Error(),
		)
		return
	}

	updatedStoreInfo, err := r.client.StoreInfo.Get()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading StoreInfo",
			"Could not read StoreInfo ID : "+err.Error(),
		)
		return
	}

	plan.StoreName = nullableString(updatedStoreInfo.StoreName)
	plan.StoreDomain = nullableString(updatedStoreInfo.StoreDomain)
	plan.UseRemoteDomain = types.BoolValue(updatedStoreInfo.UseRemoteDomain)
	plan.StoreUrl = nullableString(updatedStoreInfo.StoreUrl)
	plan.ReceiptContinueUrl = nullableString(updatedStoreInfo.ReceiptContinueUrl)
	plan.StoreEmail = nullableString(updatedStoreInfo.StoreEmail)
	plan.FromEmail = nullableString(updatedStoreInfo.FromEmail)
	plan.UseEmailDns = types.BoolValue(updatedStoreInfo.UseEmailDns)
	plan.BccOnReceiptEmail = types.BoolValue(updatedStoreInfo.BccOnReceiptEmail)
	plan.SmtpConfig = nullableString(updatedStoreInfo.SmtpConfig)
	plan.PostalCode = nullableString(updatedStoreInfo.PostalCode)
	plan.Region = nullableString(updatedStoreInfo.Region)
	plan.Country = nullableString(updatedStoreInfo.Country)
	plan.LocaleCode = nullableString(updatedStoreInfo.LocaleCode)
	plan.Timezone = nullableString(updatedStoreInfo.Timezone)
	plan.HideCurrencySymbol = types.BoolValue(updatedStoreInfo.HideCurrencySymbol)
	plan.HideDecimalCharacters = types.BoolValue(updatedStoreInfo.HideDecimalCharacters)
	plan.UseInternationalCurrencySymbol = types.BoolValue(updatedStoreInfo.UseInternationalCurrencySymbol)
	plan.Language = nullableString(updatedStoreInfo.Language)
	plan.LogoUrl = nullableString(updatedStoreInfo.LogoUrl)
	plan.CheckoutType = nullableString(updatedStoreInfo.CheckoutType)
	plan.UseWebhook = types.BoolValue(updatedStoreInfo.UseWebhook)
	plan.WebhookUrl = nullableString(updatedStoreInfo.WebhookUrl)
	plan.WebhookKey = nullableString(updatedStoreInfo.WebhookKey)
	plan.UseCartValidation = types.BoolValue(updatedStoreInfo.UseCartValidation)
	plan.UseSingleSignOn = types.BoolValue(updatedStoreInfo.UseSingleSignOn)
	plan.SingleSignOnUrl = nullableString(updatedStoreInfo.SingleSignOnUrl)
	plan.CustomerPasswordHashType = nullableString(updatedStoreInfo.CustomerPasswordHashType)
	plan.CustomerPasswordHashConfig = nullableString(updatedStoreInfo.CustomerPasswordHashConfig)
	plan.FeaturesMultiship = types.BoolValue(updatedStoreInfo.FeaturesMultiship)
	plan.ProductsRequireExpiresProperty = types.BoolValue(updatedStoreInfo.ProductsRequireExpiresProperty)
	plan.AppSessionTime = types.Int64Value(int64(updatedStoreInfo.AppSessionTime))
	plan.ShippingAddressType = nullableString(updatedStoreInfo.ShippingAddressType)
	plan.RequireSignedShippingRates = types.BoolValue(updatedStoreInfo.RequireSignedShippingRates)
	plan.UnifiedOrderEntryPassword = nullableString(updatedStoreInfo.UnifiedOrderEntryPassword)
	plan.IsMaintenanceMode = types.BoolValue(updatedStoreInfo.IsMaintenanceMode)
	plan.IsActive = types.BoolValue(updatedStoreInfo.IsActive)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *storeInfoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *storeInfoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type storeInfoModel struct {
	Id types.String `tfsdk:"id"`

	StoreName                      types.String `tfsdk:"store_name"`
	StoreDomain                    types.String `tfsdk:"store_domain"`
	UseRemoteDomain                types.Bool   `tfsdk:"use_remote_domain"`
	StoreUrl                       types.String `tfsdk:"store_url"`
	ReceiptContinueUrl             types.String `tfsdk:"receipt_continue_url"`
	StoreEmail                     types.String `tfsdk:"store_email"`
	FromEmail                      types.String `tfsdk:"from_email"`
	UseEmailDns                    types.Bool   `tfsdk:"use_email_dns"`
	BccOnReceiptEmail              types.Bool   `tfsdk:"bcc_on_receipt_email"`
	SmtpConfig                     types.String `tfsdk:"smtp_config"`
	PostalCode                     types.String `tfsdk:"postal_code"`
	Region                         types.String `tfsdk:"region"`
	Country                        types.String `tfsdk:"country"`
	LocaleCode                     types.String `tfsdk:"locale_code"`
	Timezone                       types.String `tfsdk:"timezone"`
	HideCurrencySymbol             types.Bool   `tfsdk:"hide_currency_symbol"`
	HideDecimalCharacters          types.Bool   `tfsdk:"hide_decimal_characters"`
	UseInternationalCurrencySymbol types.Bool   `tfsdk:"use_international_currency_symbol"`
	Language                       types.String `tfsdk:"language"`
	LogoUrl                        types.String `tfsdk:"logo_url"`
	CheckoutType                   types.String `tfsdk:"checkout_type"`
	UseWebhook                     types.Bool   `tfsdk:"use_webhook"`
	WebhookUrl                     types.String `tfsdk:"webhook_url"`
	WebhookKey                     types.String `tfsdk:"webhook_key"`
	UseCartValidation              types.Bool   `tfsdk:"use_cart_validation"`
	UseSingleSignOn                types.Bool   `tfsdk:"use_single_sign_on"`
	SingleSignOnUrl                types.String `tfsdk:"single_sign_on_url"`
	CustomerPasswordHashType       types.String `tfsdk:"customer_password_hash_type"`
	CustomerPasswordHashConfig     types.String `tfsdk:"customer_password_hash_config"`
	FeaturesMultiship              types.Bool   `tfsdk:"features_multiship"`
	ProductsRequireExpiresProperty types.Bool   `tfsdk:"products_require_expires_property"`
	AppSessionTime                 types.Int64  `tfsdk:"app_session_time"`
	ShippingAddressType            types.String `tfsdk:"shipping_address_type"`
	RequireSignedShippingRates     types.Bool   `tfsdk:"require_signed_shipping_rates"`
	UnifiedOrderEntryPassword      types.String `tfsdk:"unified_order_entry_password"`
	//CustomDisplayIdConfig          interface{} `tfsdk:"custom_display_id_config"`
	IsMaintenanceMode types.Bool `tfsdk:"is_maintenance_mode"`
	IsActive          types.Bool `tfsdk:"is_active"`
	//FirstPaymentDate               interface{} `tfsdk:"first_payment_date"`
	//Features                       interface{} `tfsdk:"features"`
}
