package provider

import (
	"context"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &ChannelsDataSource{}

func NewChannelsDataSource() datasource.DataSource {
	return &ChannelsDataSource{}
}

type ChannelsDataSource struct {
	client *client.Client
}

type ChannelItemModel struct {
	ChannelID   types.Int64  `tfsdk:"channel_id"`
	ChannelName types.String `tfsdk:"channel_name"`
	Description types.String `tfsdk:"description"`
	TeamID      types.Int64  `tfsdk:"team_id"`
	Status      types.String `tfsdk:"status"`
}

type ChannelsDataSourceModel struct {
	Query    types.String       `tfsdk:"query"`
	TeamIDs  types.List         `tfsdk:"team_ids"`
	Channels []ChannelItemModel `tfsdk:"channels"`
}

func (d *ChannelsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channels"
}

func (d *ChannelsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get a list of Flashduty channels (collaboration spaces).",

		Attributes: map[string]schema.Attribute{
			"query": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Search query to filter channels by name or description.",
			},
			"team_ids": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "Filter channels by team IDs.",
			},
			"channels": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of channels.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"channel_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the channel.",
						},
						"channel_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the channel.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the channel.",
						},
						"team_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the team.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status of the channel.",
						},
					},
				},
			},
		},
	}
}

func (d *ChannelsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *ChannelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ChannelsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := ""
	if !data.Query.IsNull() {
		query = data.Query.ValueString()
	}
	var teamIDs []int64
	if !data.TeamIDs.IsNull() {
		resp.Diagnostics.Append(data.TeamIDs.ElementsAs(ctx, &teamIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var allChannels []client.Channel
	page := 1
	limit := 100
	for {
		result, err := d.client.ListChannels(ctx, &client.ListChannelsRequest{
			Page:    page,
			Limit:   limit,
			Query:   query,
			TeamIDs: teamIDs,
		})
		if err != nil {
			resp.Diagnostics.AddError("Error Listing Channels", err.Error())
			return
		}
		allChannels = append(allChannels, result.Items...)
		if len(result.Items) < limit {
			break
		}
		page++
	}

	data.Channels = make([]ChannelItemModel, len(allChannels))
	for i, channel := range allChannels {
		data.Channels[i] = ChannelItemModel{
			ChannelID:   types.Int64Value(channel.ChannelID),
			ChannelName: types.StringValue(channel.ChannelName),
			Description: types.StringValue(channel.Description),
			TeamID:      types.Int64Value(channel.TeamID),
			Status:      types.StringValue(channel.Status),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
