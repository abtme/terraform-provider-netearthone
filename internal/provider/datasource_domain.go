package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/abtme/terraform-provider-netearthone/internal/client"
)

var _ datasource.DataSource = &DomainDataSource{}

type DomainDataSource struct {
	client *client.Client
}

type DomainDataSourceModel struct {
	DomainName          types.String `tfsdk:"domain_name"`
	OrderID             types.Int64  `tfsdk:"order_id"`
	Nameservers         types.List   `tfsdk:"nameservers"`
	RegistrantContactID types.Int64  `tfsdk:"registrant_contact_id"`
	AdminContactID      types.Int64  `tfsdk:"admin_contact_id"`
	TechContactID       types.Int64  `tfsdk:"tech_contact_id"`
	BillingContactID    types.Int64  `tfsdk:"billing_contact_id"`
}

func NewDomainDataSource() datasource.DataSource {
	return &DomainDataSource{}
}

func (d *DomainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (d *DomainDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up a NetearthOne domain by name, returning its order ID and current nameservers. " +
			"Use the order_id output with the netearthone_domain_nameservers resource.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "The fully-qualified domain name to look up (e.g. \"example.com\").",
			},
			"order_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The NetearthOne order ID for this domain registration.",
			},
			"nameservers": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The current nameservers assigned to the domain.",
			},
			"registrant_contact_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Contact ID of the domain registrant (owner).",
			},
			"admin_contact_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Contact ID of the administrative contact.",
			},
			"tech_contact_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Contact ID of the technical contact.",
			},
			"billing_contact_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Contact ID of the billing contact.",
			},
		},
	}
}

func (d *DomainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DomainDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	details, err := d.client.GetDomainByName(config.DomainName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to look up domain", err.Error())
		return
	}

	orderID, err := details.OrderIDInt()
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse order ID", err.Error())
		return
	}

	nsList, diags := types.ListValueFrom(ctx, types.StringType, details.Nameservers())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.OrderID = types.Int64Value(int64(orderID))
	config.Nameservers = nsList

	// Fetch contact IDs in a second call with ContactIds option.
	reg, admin, tech, billing, err := d.client.GetDomainContactIDs(orderID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read domain contact IDs", err.Error())
		return
	}
	config.RegistrantContactID = types.Int64Value(int64(reg))
	config.AdminContactID = types.Int64Value(int64(admin))
	config.TechContactID = types.Int64Value(int64(tech))
	config.BillingContactID = types.Int64Value(int64(billing))

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
