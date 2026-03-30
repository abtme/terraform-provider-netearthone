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

var _ datasource.DataSource = &DomainsDataSource{}

type DomainsDataSource struct {
	client *client.Client
}

type DomainsDataSourceModel struct {
	DomainName   types.String `tfsdk:"domain_name"`
	Status       types.List   `tfsdk:"status"`
	ProductKey   types.List   `tfsdk:"product_key"`
	NoOfRecords  types.Int64  `tfsdk:"no_of_records"`
	PageNo       types.Int64  `tfsdk:"page_no"`
	TotalRecords types.Int64  `tfsdk:"total_records"`
	Domains      types.List   `tfsdk:"domains"`
}

var domainObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"order_id":       types.StringType,
		"domain_name":    types.StringType,
		"status":         types.StringType,
		"product_key":    types.StringType,
		"expiry_date":    types.StringType,
		"creation_date":  types.StringType,
	},
}

func NewDomainsDataSource() datasource.DataSource {
	return &DomainsDataSource{}
}

func (d *DomainsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domains"
}

func (d *DomainsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists and filters NetearthOne domain registration orders.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by domain name substring.",
			},
			"status": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Filter by order status. Valid values: InActive, Active, Suspended, Pending Delete Restorable, Deleted, Archived.",
			},
			"product_key": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Filter by TLD product key (e.g. \"dotcom\", \"dotnet\", \"dotuk\").",
			},
			"no_of_records": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of records per page (default 50, max 500).",
			},
			"page_no": schema.Int64Attribute{
				Optional:    true,
				Description: "Page number for pagination (default 1).",
			},
			"total_records": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of matching records in the API.",
			},
			"domains": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of matching domain orders.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"order_id": schema.StringAttribute{
							Computed:    true,
							Description: "Order ID of the domain registration.",
						},
						"domain_name": schema.StringAttribute{
							Computed:    true,
							Description: "Fully-qualified domain name.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Current order status.",
						},
						"product_key": schema.StringAttribute{
							Computed:    true,
							Description: "TLD product key.",
						},
						"expiry_date": schema.StringAttribute{
							Computed:    true,
							Description: "Domain expiry timestamp.",
						},
						"creation_date": schema.StringAttribute{
							Computed:    true,
							Description: "Domain creation timestamp.",
						},
					},
				},
			},
		},
	}
}

func (d *DomainsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DomainsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DomainsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := client.DomainSearchParams{
		NoOfRecords: int(config.NoOfRecords.ValueInt64()),
		PageNo:      int(config.PageNo.ValueInt64()),
		DomainName:  config.DomainName.ValueString(),
	}

	if !config.Status.IsNull() && !config.Status.IsUnknown() {
		var statuses []types.String
		resp.Diagnostics.Append(config.Status.ElementsAs(ctx, &statuses, false)...)
		for _, s := range statuses {
			p.Status = append(p.Status, s.ValueString())
		}
	}

	if !config.ProductKey.IsNull() && !config.ProductKey.IsUnknown() {
		var keys []types.String
		resp.Diagnostics.Append(config.ProductKey.ElementsAs(ctx, &keys, false)...)
		for _, k := range keys {
			p.ProductKey = append(p.ProductKey, k.ValueString())
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	results, total, err := d.client.SearchDomains(p)
	if err != nil {
		resp.Diagnostics.AddError("Failed to search domains", err.Error())
		return
	}

	config.TotalRecords = types.Int64Value(int64(total))

	domainElems := make([]attr.Value, len(results))
	for i, r := range results {
		orderIDStr := fmt.Sprintf("%v", r.OrderID)
		obj, diags := types.ObjectValue(domainObjectType.AttrTypes, map[string]attr.Value{
			"order_id":      types.StringValue(orderIDStr),
			"domain_name":   types.StringValue(r.DomainName),
			"status":        types.StringValue(r.CurrentStatus),
			"product_key":   types.StringValue(r.ProductKey),
			"expiry_date":   types.StringValue(r.ExpiryDate),
			"creation_date": types.StringValue(r.CreationDate),
		})
		resp.Diagnostics.Append(diags...)
		domainElems[i] = obj
	}
	if resp.Diagnostics.HasError() {
		return
	}

	domainsList, diags := types.ListValue(domainObjectType, domainElems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Domains = domainsList
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
