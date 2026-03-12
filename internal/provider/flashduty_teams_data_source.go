package provider

import (
	"context"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &TeamsDataSource{}

func NewTeamsDataSource() datasource.DataSource {
	return &TeamsDataSource{}
}

type TeamsDataSource struct {
	client *client.Client
}

type TeamItemModel struct {
	TeamID      types.Int64  `tfsdk:"team_id"`
	TeamName    types.String `tfsdk:"team_name"`
	Description types.String `tfsdk:"description"`
}

type TeamsDataSourceModel struct {
	Query types.String    `tfsdk:"query"`
	Teams []TeamItemModel `tfsdk:"teams"`
}

func (d *TeamsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

func (d *TeamsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get a list of Flashduty teams.",

		Attributes: map[string]schema.Attribute{
			"query": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Search query to filter teams by name or description.",
			},
			"teams": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of teams.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"team_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the team.",
						},
						"team_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the team.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the team.",
						},
					},
				},
			},
		},
	}
}

func (d *TeamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *TeamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := ""
	if !data.Query.IsNull() {
		query = data.Query.ValueString()
	}

	var allTeams []client.Team
	page := 1
	limit := 100
	for {
		result, err := d.client.ListTeams(ctx, &client.ListTeamsRequest{
			Page:  page,
			Limit: limit,
			Query: query,
		})
		if err != nil {
			resp.Diagnostics.AddError("Error Listing Teams", err.Error())
			return
		}
		allTeams = append(allTeams, result.Items...)
		if len(result.Items) < limit {
			break
		}
		page++
	}

	data.Teams = make([]TeamItemModel, len(allTeams))
	for i, team := range allTeams {
		data.Teams[i] = TeamItemModel{
			TeamID:      types.Int64Value(team.TeamID),
			TeamName:    types.StringValue(team.TeamName),
			Description: types.StringValue(team.Description),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
