package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"gitlab.turnbull.uk/awxgit/terraform-provider-netearthone/internal/client"
)

var _ resource.Resource = &ContactResource{}

type ContactResource struct {
	client *client.Client
}

type ContactModel struct {
	ID           types.String `tfsdk:"id"`
	CustomerID   types.Int64  `tfsdk:"customer_id"`
	Type         types.String `tfsdk:"type"`
	Name         types.String `tfsdk:"name"`
	Company      types.String `tfsdk:"company"`
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
	FaxCC        types.String `tfsdk:"fax_cc"`
	Fax          types.String `tfsdk:"fax"`
}

func NewContactResource() resource.Resource {
	return &ContactResource{}
}

func (r *ContactResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact"
}

func (r *ContactResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a NetearthOne registrant contact record. " +
			"Contacts can be assigned to domain registrations using the netearthone_domain_contacts resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The contact ID assigned by NetearthOne.",
			},
			"customer_id": schema.Int64Attribute{
				Required:    true,
				Description: "The customer ID under which this contact is created.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Contact"),
				Description: "Contact type. One of: Contact, UkContact, EuContact, CnContact, CoContact, CaContact, DeContact, EsContact. Defaults to \"Contact\".",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Full name of the contact (max 255 characters).",
			},
			"company": schema.StringAttribute{
				Required:    true,
				Description: "Company name. Use \"N/A\" for natural persons.",
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "Email address of the contact.",
			},
			"address_line_1": schema.StringAttribute{
				Required:    true,
				Description: "Primary address line (max 64 characters).",
			},
			"address_line_2": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Secondary address line.",
			},
			"address_line_3": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Tertiary address line.",
			},
			"city": schema.StringAttribute{
				Required:    true,
				Description: "City (max 64 characters).",
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "State or province (max 64 characters).",
			},
			"country": schema.StringAttribute{
				Required:    true,
				Description: "ISO 3166-1 alpha-2 country code (e.g. \"GB\", \"US\").",
			},
			"zipcode": schema.StringAttribute{
				Required:    true,
				Description: "Postal/ZIP code.",
			},
			"phone_cc": schema.StringAttribute{
				Required:    true,
				Description: "Telephone country code (1–3 digits, e.g. \"44\" for UK).",
			},
			"phone": schema.StringAttribute{
				Required:    true,
				Description: "Telephone number (4–12 digits).",
			},
			"fax_cc": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Fax country code.",
			},
			"fax": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Fax number.",
			},
		},
	}
}

func (r *ContactResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data type",
			fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData))
		return
	}
	r.client = c
}

func modelToContactParams(m ContactModel) client.ContactParams {
	return client.ContactParams{
		CustomerID:   int(m.CustomerID.ValueInt64()),
		Type:         m.Type.ValueString(),
		Name:         m.Name.ValueString(),
		Company:      m.Company.ValueString(),
		Email:        m.Email.ValueString(),
		AddressLine1: m.AddressLine1.ValueString(),
		AddressLine2: m.AddressLine2.ValueString(),
		AddressLine3: m.AddressLine3.ValueString(),
		City:         m.City.ValueString(),
		State:        m.State.ValueString(),
		Country:      m.Country.ValueString(),
		Zipcode:      m.Zipcode.ValueString(),
		PhoneCC:      m.PhoneCC.ValueString(),
		Phone:        m.Phone.ValueString(),
		FaxCC:        m.FaxCC.ValueString(),
		Fax:          m.Fax.ValueString(),
	}
}

func contactDetailsToModel(d *client.ContactDetails, m *ContactModel) {
	m.Name = types.StringValue(d.Name)
	m.Company = types.StringValue(d.Company)
	m.Type = types.StringValue(d.Type)
	m.Email = types.StringValue(d.Email)
	m.AddressLine1 = types.StringValue(d.AddressLine1)
	m.AddressLine2 = types.StringValue(d.AddressLine2)
	m.AddressLine3 = types.StringValue(d.AddressLine3)
	m.City = types.StringValue(d.City)
	m.State = types.StringValue(d.State)
	m.Country = types.StringValue(d.Country)
	m.Zipcode = types.StringValue(d.Zipcode)
	m.PhoneCC = types.StringValue(d.PhoneCC)
	m.Phone = types.StringValue(d.Phone)
}

func (r *ContactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ContactModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactID, err := r.client.AddContact(modelToContactParams(plan))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create contact", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(contactID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ContactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ContactModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid contact ID in state", err.Error())
		return
	}

	details, err := r.client.GetContact(contactID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read contact", err.Error())
		return
	}

	contactDetailsToModel(details, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ContactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ContactModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactID, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid contact ID in state", err.Error())
		return
	}

	if err := r.client.ModifyContact(contactID, modelToContactParams(plan)); err != nil {
		resp.Diagnostics.AddError("Failed to update contact", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ContactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ContactModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid contact ID in state", err.Error())
		return
	}

	if err := r.client.DeleteContact(contactID); err != nil {
		resp.Diagnostics.AddError("Failed to delete contact", err.Error())
	}
}
