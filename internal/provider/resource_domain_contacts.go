package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/abtme/terraform-provider-netearthone/internal/client"
)

var _ resource.Resource = &DomainContactsResource{}
var _ resource.ResourceWithImportState = &DomainContactsResource{}

type DomainContactsResource struct {
	client *client.Client
}

type DomainContactsModel struct {
	ID               types.String `tfsdk:"id"`
	OrderID          types.Int64  `tfsdk:"order_id"`
	RegContactID     types.Int64  `tfsdk:"registrant_contact_id"`
	AdminContactID   types.Int64  `tfsdk:"admin_contact_id"`
	TechContactID    types.Int64  `tfsdk:"tech_contact_id"`
	BillingContactID types.Int64  `tfsdk:"billing_contact_id"`
}

func NewDomainContactsResource() resource.Resource {
	return &DomainContactsResource{}
}

func (r *DomainContactsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_contacts"
}

func (r *DomainContactsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Assigns registrant, admin, tech, and billing contacts to a NetearthOne domain registration order. " +
			"Use the netearthone_contact resource to manage the contact records themselves.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The order ID as a string (resource identifier).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"order_id": schema.Int64Attribute{
				Required:    true,
				Description: "The NetearthOne order ID of the domain registration.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"registrant_contact_id": schema.Int64Attribute{
				Required:    true,
				Description: "Contact ID to set as the registrant (owner) of the domain.",
			},
			"admin_contact_id": schema.Int64Attribute{
				Required:    true,
				Description: "Contact ID to set as the administrative contact.",
			},
			"tech_contact_id": schema.Int64Attribute{
				Required:    true,
				Description: "Contact ID to set as the technical contact.",
			},
			"billing_contact_id": schema.Int64Attribute{
				Required:    true,
				Description: "Contact ID to set as the billing contact.",
			},
		},
	}
}

func (r *DomainContactsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DomainContactsResource) applyContacts(plan *DomainContactsModel, diag interface{ AddError(string, string) }) {
	if err := r.client.ModifyDomainContacts(
		int(plan.OrderID.ValueInt64()),
		int(plan.RegContactID.ValueInt64()),
		int(plan.AdminContactID.ValueInt64()),
		int(plan.TechContactID.ValueInt64()),
		int(plan.BillingContactID.ValueInt64()),
	); err != nil {
		diag.AddError("Failed to set domain contacts", err.Error())
	}
}

func (r *DomainContactsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DomainContactsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.applyContacts(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(plan.OrderID.ValueInt64(), 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DomainContactsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DomainContactsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reg, admin, tech, billing, err := r.client.GetDomainContactIDs(int(state.OrderID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read domain contact IDs", err.Error())
		return
	}

	state.RegContactID = types.Int64Value(int64(reg))
	state.AdminContactID = types.Int64Value(int64(admin))
	state.TechContactID = types.Int64Value(int64(tech))
	state.BillingContactID = types.Int64Value(int64(billing))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DomainContactsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DomainContactsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.applyContacts(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DomainContactsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by order ID (e.g. terraform import netearthone_domain_contacts.this 124814418)
	orderID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected a numeric order ID, got: "+req.ID)
		return
	}

	reg, admin, tech, billing, err := r.client.GetDomainContactIDs(int(orderID))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read domain contacts during import", err.Error())
		return
	}

	state := DomainContactsModel{
		ID:               types.StringValue(req.ID),
		OrderID:          types.Int64Value(orderID),
		RegContactID:     types.Int64Value(int64(reg)),
		AdminContactID:   types.Int64Value(int64(admin)),
		TechContactID:    types.Int64Value(int64(tech)),
		BillingContactID: types.Int64Value(int64(billing)),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DomainContactsResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Contacts cannot be unassigned from a domain — removing this resource
	// from state simply stops Terraform from managing the contact assignments.
}
