package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &MemberInviteResource{}
	_ resource.ResourceWithImportState = &MemberInviteResource{}
)

func NewMemberInviteResource() resource.Resource {
	return &MemberInviteResource{}
}

type MemberInviteResource struct {
	client *client.Client
}

type MemberInviteResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Email       types.String `tfsdk:"email"`
	Phone       types.String `tfsdk:"phone"`
	CountryCode types.String `tfsdk:"country_code"`
	MemberName  types.String `tfsdk:"member_name"`
	RefID       types.String `tfsdk:"ref_id"`
	RoleIDs     types.List   `tfsdk:"role_ids"`
	MemberID    types.Int64  `tfsdk:"member_id"`
}

func (r *MemberInviteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_member_invite"
}

func (r *MemberInviteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Invites a member to Flashduty. The invited member will be created and can login to activate their account.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the invited member.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The email address of the member. At least one of `email` or `phone` must be provided.",
			},
			"phone": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The phone number of the member. At least one of `email` or `phone` must be provided.",
			},
			"country_code": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("CN"),
				MarkdownDescription: "The country code for the phone number. Defaults to `CN`.",
			},
			"member_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the member. Required when using phone number.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ref_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An external reference ID for the member.",
			},
			"role_ids": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "A list of role IDs to assign to the member.",
			},
			"member_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The Flashduty member ID assigned after invitation.",
			},
		},
	}
}

func (r *MemberInviteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (r *MemberInviteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MemberInviteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Email.IsNull() && data.Phone.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"At least one of 'email' or 'phone' must be provided.",
		)
		return
	}

	if !data.Phone.IsNull() && data.MemberName.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"'member_name' is required when 'phone' is provided.",
		)
		return
	}

	member := client.MemberInvite{}

	if !data.Email.IsNull() {
		member.Email = data.Email.ValueString()
	}
	if !data.Phone.IsNull() {
		member.Phone = data.Phone.ValueString()
	}
	if !data.CountryCode.IsNull() {
		member.CountryCode = data.CountryCode.ValueString()
	}
	if !data.MemberName.IsNull() {
		member.MemberName = data.MemberName.ValueString()
	}
	if !data.RefID.IsNull() {
		member.RefID = data.RefID.ValueString()
	}

	if !data.RoleIDs.IsNull() {
		resp.Diagnostics.Append(data.RoleIDs.ElementsAs(ctx, &member.RoleIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, "Inviting member to Flashduty", map[string]interface{}{
		"email":       member.Email,
		"phone":       member.Phone,
		"member_name": member.MemberName,
	})

	result, err := r.client.InviteMembers(ctx, []client.MemberInvite{member})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Inviting Member",
			fmt.Sprintf("Unable to invite member, got error: %s", err),
		)
		return
	}

	if result == nil || len(result.Items) == 0 {
		resp.Diagnostics.AddError(
			"Error Inviting Member",
			"No member was returned from the API response.",
		)
		return
	}

	invitedMember := result.Items[0]
	data.ID = types.StringValue(strconv.FormatInt(invitedMember.MemberID, 10))
	data.MemberID = types.Int64Value(invitedMember.MemberID)
	data.MemberName = types.StringValue(invitedMember.MemberName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MemberInviteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MemberInviteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	memberID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Member ID", err.Error())
		return
	}

	member, err := r.client.GetMemberByID(ctx, memberID)
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Member", err.Error())
		return
	}

	data.MemberID = types.Int64Value(memberID)
	data.MemberName = types.StringValue(member.MemberName)
	if member.Email != "" {
		data.Email = types.StringValue(member.Email)
	}
	if member.Phone != "" {
		data.Phone = types.StringValue(member.Phone)
	}
	if member.CountryCode != "" {
		data.CountryCode = types.StringValue(member.CountryCode)
	}
	if member.RefID != "" {
		data.RefID = types.StringValue(member.RefID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MemberInviteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MemberInviteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	memberID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Member ID", err.Error())
		return
	}

	updateReq := &client.UpdateMemberRequest{
		MemberID: memberID,
		Updates:  client.MemberUpdatePayload{},
	}

	if !data.MemberName.IsNull() {
		updateReq.Updates.MemberName = data.MemberName.ValueString()
	}
	if !data.RefID.IsNull() {
		updateReq.Updates.RefID = data.RefID.ValueString()
	}
	if !data.Phone.IsNull() {
		updateReq.Updates.Phone = data.Phone.ValueString()
	}
	if !data.CountryCode.IsNull() {
		updateReq.Updates.CountryCode = data.CountryCode.ValueString()
	}
	if !data.Email.IsNull() {
		updateReq.Updates.Email = data.Email.ValueString()
	}

	err = r.client.UpdateMember(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Member", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MemberInviteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MemberInviteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	memberID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Member ID", err.Error())
		return
	}

	err = r.client.DeleteMember(ctx, &client.DeleteMemberRequest{
		MemberID: memberID,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Member", err.Error())
		return
	}

	tflog.Debug(ctx, "Deleted member", map[string]interface{}{
		"member_id": memberID,
	})
}

func (r *MemberInviteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
