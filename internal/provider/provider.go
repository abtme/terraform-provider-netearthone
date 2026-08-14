package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/abtme/terraform-provider-netearthone/internal/client"
)

const defaultBaseURL = "https://httpapi.com"

var _ provider.Provider = &NetearthOneProvider{}

type NetearthOneProvider struct {
	version string
}

type NetearthOneProviderModel struct {
	BaseURL    types.String `tfsdk:"base_url"`
	AuthUserID types.Int64  `tfsdk:"auth_userid"`
	APIKey     types.String `tfsdk:"api_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NetearthOneProvider{version: version}
	}
}

func (p *NetearthOneProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "netearthone"
	resp.Version = p.version
}

func (p *NetearthOneProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for managing NetearthOne domain resources via the HTTP API.",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "NetearthOne API base URL. Defaults to https://httpapi.com. Can also be set via NETEARTHONE_BASE_URL environment variable.",
			},
			"auth_userid": schema.Int64Attribute{
				Optional:    true,
				Description: "NetearthOne auth-userid for API authentication. Can also be set via NETEARTHONE_AUTH_USERID environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "NetearthOne API key. Can also be set via NETEARTHONE_API_KEY environment variable.",
			},
		},
	}
}

func (p *NetearthOneProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config NetearthOneProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	baseURL := defaultBaseURL
	if v := os.Getenv("NETEARTHONE_BASE_URL"); v != "" {
		baseURL = v
	}
	if !config.BaseURL.IsNull() && !config.BaseURL.IsUnknown() {
		baseURL = config.BaseURL.ValueString()
	}

	var authUserID int
	if v := os.Getenv("NETEARTHONE_AUTH_USERID"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddError("Invalid NETEARTHONE_AUTH_USERID", "Must be an integer: "+err.Error())
			return
		}
		authUserID = id
	}
	if !config.AuthUserID.IsNull() && !config.AuthUserID.IsUnknown() {
		authUserID = int(config.AuthUserID.ValueInt64())
	}
	if authUserID == 0 {
		resp.Diagnostics.AddError("Missing auth_userid", "auth_userid must be set in provider config or NETEARTHONE_AUTH_USERID environment variable.")
		return
	}

	apiKey := os.Getenv("NETEARTHONE_API_KEY")
	if !config.APIKey.IsNull() && !config.APIKey.IsUnknown() {
		apiKey = config.APIKey.ValueString()
	}
	if apiKey == "" {
		resp.Diagnostics.AddError("Missing api_key", "api_key must be set in provider config or NETEARTHONE_API_KEY environment variable.")
		return
	}

	c := client.NewClient(baseURL, authUserID, apiKey)
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *NetearthOneProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDomainNameserversResource,
		NewDomainPrivacyResource,
		NewDomainContactsResource,
		NewChildNameserverResource,
		NewContactResource,
	}
}

func (p *NetearthOneProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDomainDataSource,
		NewContactDataSource,
		NewDomainAvailabilityDataSource,
		NewDomainsDataSource,
		NewDefaultNameserversDataSource,
	}
}
