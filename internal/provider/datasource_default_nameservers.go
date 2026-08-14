package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/abtme/terraform-provider-netearthone/internal/client"
)

var _ datasource.DataSource = &DefaultNameserversDataSource{}

type DefaultNameserversDataSource struct {
	client *client.Client
}

type DefaultNameserversDataSourceModel struct {
	CustomerID  types.Int64 `tfsdk:"customer_id"`
	Nameservers types.List  `tfsdk:"nameservers"`
}

func NewDefaultNameserversDataSource() datasource.DataSource {
	return &DefaultNameserversDataSource{}
}

func (d *DefaultNameserversDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_nameservers"
}

func (d *DefaultNameserversDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the default nameservers configured for a NetearthOne customer.",
		Attributes: map[string]schema.Attribute{
			"customer_id": schema.Int64Attribute{
				Required:    true,
				Description: "The NetearthOne customer ID whose default nameservers to retrieve.",
			},
			"nameservers": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The customer's configured default nameservers.",
			},
		},
	}
}

func (d *DefaultNameserversDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DefaultNameserversDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DefaultNameserversDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ns, err := d.client.GetDefaultNameservers(int(config.CustomerID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch default nameservers", err.Error())
		return
	}

	nsList, diags := types.ListValueFrom(ctx, types.StringType, ns)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Nameservers = nsList
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
