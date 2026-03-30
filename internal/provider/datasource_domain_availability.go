package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"gitlab.turnbull.uk/awxgit/terraform-provider-netearthone/internal/client"
)

var _ datasource.DataSource = &DomainAvailabilityDataSource{}

type DomainAvailabilityDataSource struct {
	client *client.Client
}

type DomainAvailabilityDataSourceModel struct {
	Domains []types.String        `tfsdk:"domains"`
	TLDs    []types.String        `tfsdk:"tlds"`
	Results types.Map             `tfsdk:"results"`
}

func NewDomainAvailabilityDataSource() datasource.DataSource {
	return &DomainAvailabilityDataSource{}
}

func (d *DomainAvailabilityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_availability"
}

func (d *DomainAvailabilityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Checks whether one or more domain names are available for registration. " +
			"Returns a map of \"domain.tld\" to status (\"available\", \"regthroughus\", \"regthroughothers\", or \"unknown\").",
		Attributes: map[string]schema.Attribute{
			"domains": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of second-level domain names to check (without TLD, e.g. [\"example\", \"mysite\"]).",
			},
			"tlds": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of TLD extensions to check against (e.g. [\"com\", \"net\", \"uk\"]).",
			},
			"results": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Map of \"domain.tld\" to availability status string.",
			},
		},
	}
}

func (d *DomainAvailabilityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data type",
			fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *DomainAvailabilityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DomainAvailabilityDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domains := make([]string, len(config.Domains))
	for i, v := range config.Domains {
		domains[i] = v.ValueString()
	}
	tlds := make([]string, len(config.TLDs))
	for i, v := range config.TLDs {
		tlds[i] = v.ValueString()
	}

	availability, err := d.client.CheckDomainAvailability(domains, tlds)
	if err != nil {
		resp.Diagnostics.AddError("Failed to check domain availability", err.Error())
		return
	}

	resultElems := make(map[string]attr.Value, len(availability))
	for domain, status := range availability {
		resultElems[domain] = types.StringValue(status)
	}

	resultsMap, diags := types.MapValue(types.StringType, resultElems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Results = resultsMap
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
