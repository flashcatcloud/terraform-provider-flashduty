package provider

import (
	"context"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &MemberDataSource{}

func NewMemberDataSource() datasource.DataSource {
	return &MemberDataSource{}
}

type MemberDataSource struct {
	client *client.Client
}

type MemberDataSourceModel struct {
	MemberID      types.Int64  `tfsdk:"member_id"`
	MemberName    types.String `tfsdk:"member_name"`
	Email         types.String `tfsdk:"email"`
	Phone         types.String `tfsdk:"phone"`
	Status        types.String `tfsdk:"status"`
	EmailVerified types.Bool   `tfsdk:"email_verified"`
	PhoneVerified types.Bool   `tfsdk:"phone_verified"`
	CreatedAt     types.Int64  `tfsdk:"created_at"`
	UpdatedAt     types.Int64  `tfsdk:"updated_at"`
}

func (d *MemberDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_member"
}

func (d *MemberDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a Flashduty member.",

		Attributes: map[string]schema.Attribute{
			"member_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the member.",
			},
			"member_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the member.",
			},
			"email": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The email address of the member.",
			},
			"phone": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The phone number of the member (encrypted).",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the member.",
			},
			"email_verified": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the email is verified.",
			},
			"phone_verified": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the phone is verified.",
			},
			"created_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The creation timestamp.",
			},
			"updated_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The last update timestamp.",
			},
		},
	}
}

func (d *MemberDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *MemberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MemberDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	member, err := d.client.GetMemberByID(ctx, data.MemberID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Member", err.Error())
		return
	}

	data.MemberName = types.StringValue(member.MemberName)
	data.Email = types.StringValue(member.Email)
	data.Phone = types.StringValue(member.Phone)
	data.Status = types.StringValue(member.Status)
	data.EmailVerified = types.BoolValue(member.EmailVerified)
	data.PhoneVerified = types.BoolValue(member.PhoneVerified)
	data.CreatedAt = types.Int64Value(member.CreatedAt)
	data.UpdatedAt = types.Int64Value(member.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
