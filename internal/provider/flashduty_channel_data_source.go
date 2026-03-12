package provider

import (
	"context"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &ChannelDataSource{}

func NewChannelDataSource() datasource.DataSource {
	return &ChannelDataSource{}
}

type ChannelDataSource struct {
	client *client.Client
}

type ChannelDataSourceModel struct {
	ChannelID               types.Int64  `tfsdk:"channel_id"`
	ChannelName             types.String `tfsdk:"channel_name"`
	Description             types.String `tfsdk:"description"`
	TeamID                  types.Int64  `tfsdk:"team_id"`
	ManagingTeamIDs         types.List   `tfsdk:"managing_team_ids"`
	AutoResolveTimeout      types.Int64  `tfsdk:"auto_resolve_timeout"`
	AutoResolveMode         types.String `tfsdk:"auto_resolve_mode"`
	IsPrivate               types.Bool   `tfsdk:"is_private"`
	DisableOutlierDetection types.Bool   `tfsdk:"disable_outlier_detection"`
	DisableAutoClose        types.Bool   `tfsdk:"disable_auto_close"`
	Status                  types.String `tfsdk:"status"`
	CreatedAt               types.Int64  `tfsdk:"created_at"`
	UpdatedAt               types.Int64  `tfsdk:"updated_at"`
}

func (d *ChannelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

func (d *ChannelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a Flashduty channel (collaboration space).",

		Attributes: map[string]schema.Attribute{
			"channel_id": schema.Int64Attribute{
				Required:            true,
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
				MarkdownDescription: "The ID of the team that owns this channel.",
			},
			"managing_team_ids": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "IDs of managing teams.",
			},
			"auto_resolve_timeout": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Auto-resolve timeout in seconds.",
			},
			"auto_resolve_mode": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Auto-resolve mode: `trigger` or `update`.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the channel.",
			},
			"is_private": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the channel is private.",
			},
			"disable_outlier_detection": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether outlier detection is disabled.",
			},
			"disable_auto_close": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether automatic incident closure is disabled.",
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

func (d *ChannelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *ChannelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ChannelDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channel, err := d.client.GetChannel(ctx, &client.GetChannelRequest{
		ChannelID: data.ChannelID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Channel", err.Error())
		return
	}

	data.ChannelName = types.StringValue(channel.ChannelName)
	data.Description = types.StringValue(channel.Description)
	data.TeamID = types.Int64Value(channel.TeamID)
	data.Status = types.StringValue(channel.Status)
	data.IsPrivate = types.BoolValue(channel.IsPrivate)
	data.DisableOutlierDetection = types.BoolValue(channel.DisableOutlierDetection)
	data.DisableAutoClose = types.BoolValue(channel.DisableAutoClose)
	data.CreatedAt = types.Int64Value(channel.CreatedAt)
	data.UpdatedAt = types.Int64Value(channel.UpdatedAt)

	if len(channel.ManagingTeamIDs) > 0 {
		ids, d := types.ListValueFrom(ctx, types.Int64Type, channel.ManagingTeamIDs)
		resp.Diagnostics.Append(d...)
		data.ManagingTeamIDs = ids
	} else {
		data.ManagingTeamIDs = types.ListNull(types.Int64Type)
	}

	data.AutoResolveTimeout = types.Int64Value(int64(channel.AutoResolveTimeout))
	data.AutoResolveMode = types.StringValue(channel.AutoResolveMode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
