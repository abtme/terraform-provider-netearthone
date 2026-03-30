package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"gitlab.turnbull.uk/awxgit/terraform-provider-netearthone/internal/client"
)

var _ resource.Resource = &DomainPrivacyResource{}

type DomainPrivacyResource struct {
	client *client.Client
}

type DomainPrivacyModel struct {
	ID               types.String `tfsdk:"id"`
	OrderID          types.Int64  `tfsdk:"order_id"`
	PrivacyProtected types.Bool   `tfsdk:"privacy_protected"`
	Reason           types.String `tfsdk:"reason"`
}

func NewDomainPrivacyResource() resource.Resource {
	return &DomainPrivacyResource{}
}

func (r *DomainPrivacyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_privacy"
}

func (r *DomainPrivacyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages WHOIS privacy protection for a NetearthOne domain registration order.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The order ID as a string (resource identifier).",
			},
			"order_id": schema.Int64Attribute{
				Required:    true,
				Description: "The NetearthOne order ID of the domain registration.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"privacy_protected": schema.BoolAttribute{
				Required:    true,
				Description: "Set to true to enable WHOIS privacy protection, false to disable.",
			},
			"reason": schema.StringAttribute{
				Required:    true,
				Description: "Reason for changing the privacy protection status (required by the API).",
			},
		},
	}
}

func (r *DomainPrivacyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DomainPrivacyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DomainPrivacyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.ModifyPrivacyProtection(
		int(plan.OrderID.ValueInt64()),
		plan.PrivacyProtected.ValueBool(),
		plan.Reason.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Failed to set privacy protection", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(plan.OrderID.ValueInt64(), 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DomainPrivacyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DomainPrivacyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	details, err := r.client.GetDomainDetails(int(state.OrderID.ValueInt64()), []string{"DomainStatus"})
	if err != nil {
		resp.Diagnostics.AddError("Failed to read domain privacy status", err.Error())
		return
	}

	state.PrivacyProtected = types.BoolValue(details.PrivacyEnabled())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DomainPrivacyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DomainPrivacyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.ModifyPrivacyProtection(
		int(plan.OrderID.ValueInt64()),
		plan.PrivacyProtected.ValueBool(),
		plan.Reason.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Failed to update privacy protection", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(plan.OrderID.ValueInt64(), 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DomainPrivacyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Removing this resource from state does not automatically disable privacy —
	// the operator should set privacy_protected = false before destroying.
}
