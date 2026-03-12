package provider

import (
	"context"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &TeamDataSource{}

func NewTeamDataSource() datasource.DataSource {
	return &TeamDataSource{}
}

type TeamDataSource struct {
	client *client.Client
}

type TeamDataSourceModel struct {
	TeamID      types.Int64  `tfsdk:"team_id"`
	TeamName    types.String `tfsdk:"team_name"`
	Description types.String `tfsdk:"description"`
	RefID       types.String `tfsdk:"ref_id"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	UpdatedAt   types.Int64  `tfsdk:"updated_at"`
}

func (d *TeamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *TeamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a Flashduty team.",

		Attributes: map[string]schema.Attribute{
			"team_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The ID of the team. One of `team_id`, `team_name`, or `ref_id` must be specified.",
			},
			"team_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the team.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the team.",
			},
			"ref_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The external reference ID of the team.",
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

func (d *TeamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *TeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.TeamID.IsNull() && data.TeamName.IsNull() && data.RefID.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"One of `team_id`, `team_name`, or `ref_id` must be specified.",
		)
		return
	}

	getReq := &client.GetTeamRequest{}
	if !data.TeamID.IsNull() {
		getReq.TeamID = data.TeamID.ValueInt64()
	}
	if !data.TeamName.IsNull() {
		getReq.TeamName = data.TeamName.ValueString()
	}
	if !data.RefID.IsNull() {
		getReq.RefID = data.RefID.ValueString()
	}

	team, err := d.client.GetTeam(ctx, getReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Team", err.Error())
		return
	}

	data.TeamID = types.Int64Value(team.TeamID)
	data.TeamName = types.StringValue(team.TeamName)
	data.Description = types.StringValue(team.Description)
	if team.RefID != "" {
		data.RefID = types.StringValue(team.RefID)
	}
	data.CreatedAt = types.Int64Value(team.CreatedAt)
	data.UpdatedAt = types.Int64Value(team.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
