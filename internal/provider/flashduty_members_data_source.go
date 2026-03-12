package provider

import (
	"context"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &MembersDataSource{}

func NewMembersDataSource() datasource.DataSource {
	return &MembersDataSource{}
}

type MembersDataSource struct {
	client *client.Client
}

type MemberItemModel struct {
	MemberID   types.Int64  `tfsdk:"member_id"`
	MemberName types.String `tfsdk:"member_name"`
	Email      types.String `tfsdk:"email"`
	Status     types.String `tfsdk:"status"`
}

type MembersDataSourceModel struct {
	Query   types.String      `tfsdk:"query"`
	Members []MemberItemModel `tfsdk:"members"`
}

func (d *MembersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_members"
}

func (d *MembersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get a list of Flashduty members.",

		Attributes: map[string]schema.Attribute{
			"query": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Search query to filter members by name, email, or phone.",
			},
			"members": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of members.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"member_id": schema.Int64Attribute{
							Computed:            true,
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
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status of the member.",
						},
					},
				},
			},
		},
	}
}

func (d *MembersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *MembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MembersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := data.Query.ValueString()
	var allMembers []client.Member
	page := 1
	limit := 100
	for {
		result, err := d.client.ListMembers(ctx, page, limit, query)
		if err != nil {
			resp.Diagnostics.AddError("Error Listing Members", err.Error())
			return
		}
		allMembers = append(allMembers, result.Items...)
		if len(result.Items) < limit {
			break
		}
		page++
	}

	data.Members = make([]MemberItemModel, len(allMembers))
	for i, member := range allMembers {
		data.Members[i] = MemberItemModel{
			MemberID:   types.Int64Value(member.MemberID),
			MemberName: types.StringValue(member.MemberName),
			Email:      types.StringValue(member.Email),
			Status:     types.StringValue(member.Status),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
