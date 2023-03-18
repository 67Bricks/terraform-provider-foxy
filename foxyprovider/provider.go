package foxyprovider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"reflect"
	"strings"
	"terraform-provider-foxycart/foxyclient"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &foxyProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &foxyProvider{}
}

// foxyProvider is the provider implementation.
type foxyProvider struct{}

// Metadata returns the provider type name.
func (p *foxyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "foxy"
}

// Schema defines the provider-level schema for configuration data.
func (p *foxyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Foxy",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Description: "Base URL for Foxy - typically https://api.foxycart.com. May also be provided via FOXY_BASEURL environment variable.",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "Client ID for accessing Foxy with OAuth, obtained from the Foxy admin UI. May also be provided via FOXY_CLIENTID environment variable.",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Client secret for accessing Foxy with OAuth, obtained from the Foxy admin UI. May also be provided via FOXY_CLIENTSECRET environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"refresh_token": schema.StringAttribute{
				Description: "Refresh token for accessing Foxy with OAuth, obtained from the Foxy admin UI. May also be provided via FOXY_REFRESHTOKEN environment variable.",
				Optional:    true,
			},
		},
	}
}

func (p *foxyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Foxy client")

	// Retrieve provider data from configuration
	var config foxyProviderModel
	// "Diagnostics" is an array of "Diagnostic", with additional methods like "Contains"
	// Diagnostic has helper methods like Equal
	var diags diag.Diagnostics = req.Config.Get(ctx, &config) // a "variadic function" i.e. varargs
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		// resp is a pointer to the passed-in one; so although we return nothing, resp has been mutated in-place
		return
	}

	// --- My edits
	type FieldRef[U string, T any] struct {
		Name  U
		Value T
	}
	fieldsToCheck := []FieldRef[string, types.String]{
		{"BaseUrl", config.BaseUrl}, {"ClientId", config.ClientId}, {"ClientSecret", config.ClientSecret}, {"RefreshToken", config.RefreshToken},
	}
	for _, f := range fieldsToCheck {
		if f.Value.IsUnknown() {
			fieldDefinition, _ := reflect.TypeOf(config).FieldByName(f.Name)
			terraformName := reflect.StructTag.Get(fieldDefinition.Tag, "tfsdk")
			envName := "FOXY_" + strings.ToUpper(f.Name)
			resp.Diagnostics.AddAttributeError(
				path.Root(terraformName),
				"Unknown Foxy API "+f.Name,
				fmt.Sprintf("The provider cannot create the Foxy API client as there is an unknown configuration value for the Foxy API %s. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the %s environment variable.", terraformName, envName),
			)
		}
	}
	// ---- end

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	baseUrl := os.Getenv("FOXY_BASEURL")
	clientId := os.Getenv("FOXY_CLIENTID")
	clientSecret := os.Getenv("FOXY_CLIENTSECRET")
	refreshToken := os.Getenv("FOXY_REFRESHTOKEN")

	// golang does not have a ternary operator, or a null-coalescing operator
	// types.StringValue is broadly equivalent to an Option[String] in another language
	if !config.BaseUrl.IsNull() {
		baseUrl = config.BaseUrl.ValueString()
	}
	if !config.ClientId.IsNull() {
		clientId = config.ClientId.ValueString()
	}
	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}
	if !config.RefreshToken.IsNull() {
		refreshToken = config.RefreshToken.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if baseUrl == "" {
		p.reportMissingValue(&resp.Diagnostics, "baseUrl")
	}
	if clientId == "" {
		p.reportMissingValue(&resp.Diagnostics, "clientId")
	}
	if clientSecret == "" {
		p.reportMissingValue(&resp.Diagnostics, "clientSecret")
	}
	if refreshToken == "" {
		p.reportMissingValue(&resp.Diagnostics, "refreshToken")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// "Context" is a core Go thing which is basically an info holder; SetField adds to it and returns the updated version
	// All values in context are displayed when logging
	ctx = tflog.SetField(ctx, "foxy_baseUrl", baseUrl)
	ctx = tflog.SetField(ctx, "foxy_clientId", clientId)
	ctx = tflog.SetField(ctx, "foxy_clientSecret", clientSecret)
	ctx = tflog.SetField(ctx, "foxy_refreshToken", refreshToken)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "foxy_clientSecret")

	tflog.Debug(ctx, "Creating Foxy client")

	foxy, err := foxyclient.New(baseUrl, clientId, clientSecret, refreshToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Foxy API Client",
			"An error occurred when creating the Foxy API client - perhaps the client details are incorrect. Error is : "+err.Error(),
		)
		return
	}

	// Make the Foxy client available during DataSource and Resource type Configure methods.
	// DataSourceData and ResourceData are both "any" - they are there so they are available when passed in to the other functions
	resp.DataSourceData = &foxy
	resp.ResourceData = &foxy

	tflog.Info(ctx, "Configured Foxy client", map[string]any{"success": true})
}

func (p *foxyProvider) reportMissingValue(diags *diag.Diagnostics, missingField string) {
	diags.AddAttributeError(
		path.Root(missingField),
		"Missing Foxy API "+missingField,
		"The provider cannot create the Foxy API client as there is a missing or empty value for the Foxy "+missingField+". "+
			"Set the baseUrl value in the configuration or use the FOXY_"+strings.ToUpper(missingField)+" environment variable. "+
			"If either is already set, ensure the value is not empty.",
	)
}

// DataSources defines the data sources implemented in the provider.
func (p *foxyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	// An array of functions, taking no arguments, each returning a DataSource
	//return nil
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *foxyProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWebhookResource,
		NewCartTemplateResource,
		NewCartIncludeTemplateResource,
		NewCheckoutTemplateResource,
		NewReceiptTemplateResource,
		NewEmailTemplateResource,
		NewStoreInfoResource,
	}
}

// foxyProviderModel maps provider schema data to a Go type.
type foxyProviderModel struct {
	// The tfsdk bit is an (non-typesafe) annotation to the struct, visible when doing reflection
	BaseUrl      types.String `tfsdk:"base_url"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	RefreshToken types.String `tfsdk:"refresh_token"`
}

// @todo Perhaps we could avoid this by changing the JSON serialization in the client to remove omitempty?
func nullableString(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
