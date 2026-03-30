package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"gitlab.turnbull.uk/awxgit/terraform-provider-netearthone/internal/client"
)

var _ resource.Resource = &DomainNameserversResource{}

type DomainNameserversResource struct {
	client *client.Client
}

type DomainNameserversModel struct {
	ID          types.String `tfsdk:"id"`
	OrderID     types.Int64  `tfsdk:"order_id"`
	Nameservers types.List   `tfsdk:"nameservers"`
}

func NewDomainNameserversResource() resource.Resource {
	return &DomainNameserversResource{}
}

func (r *DomainNameserversResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_nameservers"
}

func (r *DomainNameserversResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the nameservers for a NetearthOne domain registration order. " +
			"This resource updates the nameservers on an existing domain — it does not register or delete domains.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The order ID as a string (used as the resource identifier).",
			},
			"order_id": schema.Int64Attribute{
				Required:    true,
				Description: "The NetearthOne order ID of the domain registration.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"nameservers": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of nameservers to assign to the domain (minimum 2, maximum 13).",
			},
		},
	}
}

func (r *DomainNameserversResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DomainNameserversResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DomainNameserversModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ns, diags := toStringSlice(ctx, plan.Nameservers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.ModifyNameservers(int(plan.OrderID.ValueInt64()), ns); err != nil {
		resp.Diagnostics.AddError("Failed to set nameservers", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(plan.OrderID.ValueInt64(), 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DomainNameserversResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DomainNameserversModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ns, err := r.client.GetDomainNameservers(int(state.OrderID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read nameservers", err.Error())
		return
	}

	nsList, diags := types.ListValueFrom(ctx, types.StringType, ns)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Nameservers = nsList
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DomainNameserversResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DomainNameserversModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ns, diags := toStringSlice(ctx, plan.Nameservers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.ModifyNameservers(int(plan.OrderID.ValueInt64()), ns); err != nil {
		resp.Diagnostics.AddError("Failed to update nameservers", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(plan.OrderID.ValueInt64(), 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DomainNameserversResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Nameservers cannot be "deleted" — removing this resource from Terraform
	// state simply stops managing them. No API call is made.
}

func toStringSlice(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	var elems []types.String
	diags := list.ElementsAs(ctx, &elems, false)
	if diags.HasError() {
		return nil, diags
	}
	result := make([]string, len(elems))
	for i, e := range elems {
		result[i] = e.ValueString()
	}
	return result, diags
}
