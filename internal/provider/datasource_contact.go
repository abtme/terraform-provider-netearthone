package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/abtme/terraform-provider-netearthone/internal/client"
)

var _ datasource.DataSource = &ContactDataSource{}

type ContactDataSource struct {
	client *client.Client
}

type ContactDataSourceModel struct {
	ContactID    types.Int64  `tfsdk:"contact_id"`
	Name         types.String `tfsdk:"name"`
	Company      types.String `tfsdk:"company"`
	Type         types.String `tfsdk:"type"`
	Email        types.String `tfsdk:"email"`
	AddressLine1 types.String `tfsdk:"address_line_1"`
	AddressLine2 types.String `tfsdk:"address_line_2"`
	AddressLine3 types.String `tfsdk:"address_line_3"`
	City         types.String `tfsdk:"city"`
	State        types.String `tfsdk:"state"`
	Country      types.String `tfsdk:"country"`
	Zipcode      types.String `tfsdk:"zipcode"`
	PhoneCC      types.String `tfsdk:"phone_cc"`
	Phone        types.String `tfsdk:"phone"`
	Status       types.String `tfsdk:"status"`
}

func NewContactDataSource() datasource.DataSource {
	return &ContactDataSource{}
}

func (d *ContactDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact"
}

func (d *ContactDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches details of a NetearthOne contact record by its contact ID.",
		Attributes: map[string]schema.Attribute{
			"contact_id": schema.Int64Attribute{
				Required:    true,
				Description: "The NetearthOne contact ID to look up.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Full name of the contact.",
			},
			"company": schema.StringAttribute{
				Computed:    true,
				Description: "Company name.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Contact type (e.g. Contact, UkContact).",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "Email address.",
			},
			"address_line_1": schema.StringAttribute{
				Computed:    true,
				Description: "Primary address line.",
			},
			"address_line_2": schema.StringAttribute{
				Computed:    true,
				Description: "Secondary address line.",
			},
			"address_line_3": schema.StringAttribute{
				Computed:    true,
				Description: "Tertiary address line.",
			},
			"city": schema.StringAttribute{
				Computed:    true,
				Description: "City.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "State or province.",
			},
			"country": schema.StringAttribute{
				Computed:    true,
				Description: "ISO 3166-1 alpha-2 country code.",
			},
			"zipcode": schema.StringAttribute{
				Computed:    true,
				Description: "Postal/ZIP code.",
			},
			"phone_cc": schema.StringAttribute{
				Computed:    true,
				Description: "Telephone country code.",
			},
			"phone": schema.StringAttribute{
				Computed:    true,
				Description: "Telephone number.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current status of the contact (e.g. Active, InActive).",
			},
		},
	}
}

func (d *ContactDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ContactDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ContactDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	details, err := d.client.GetContact(int(config.ContactID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read contact", err.Error())
		return
	}

	contactID, err := strconv.Atoi(fmt.Sprintf("%v", details.ContactID))
	if err == nil {
		config.ContactID = types.Int64Value(int64(contactID))
	}

	config.Name = types.StringValue(details.Name)
	config.Company = types.StringValue(details.Company)
	config.Type = types.StringValue(details.Type)
	config.Email = types.StringValue(details.Email)
	config.AddressLine1 = types.StringValue(details.AddressLine1)
	config.AddressLine2 = types.StringValue(details.AddressLine2)
	config.AddressLine3 = types.StringValue(details.AddressLine3)
	config.City = types.StringValue(details.City)
	config.State = types.StringValue(details.State)
	config.Country = types.StringValue(details.Country)
	config.Zipcode = types.StringValue(details.Zipcode)
	config.PhoneCC = types.StringValue(details.PhoneCC)
	config.Phone = types.StringValue(details.Phone)
	config.Status = types.StringValue(details.CurrentStatus)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
