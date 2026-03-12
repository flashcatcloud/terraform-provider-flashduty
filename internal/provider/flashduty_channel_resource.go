package provider

import (
	"context"
	"strconv"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ChannelResource{}
	_ resource.ResourceWithImportState = &ChannelResource{}
)

func NewChannelResource() resource.Resource {
	return &ChannelResource{}
}

type ChannelResource struct {
	client *client.Client
}

type ChannelGroupCaseModel struct {
	If     []ChannelGroupFilterModel `tfsdk:"if"`
	Equals types.List                `tfsdk:"equals"`
}

type ChannelGroupFilterModel struct {
	Key  types.String `tfsdk:"key"`
	Oper types.String `tfsdk:"oper"`
	Vals types.List   `tfsdk:"vals"`
}

type ChannelGroupModel struct {
	Method            types.String            `tfsdk:"method"`
	Cases             []ChannelGroupCaseModel `tfsdk:"cases"`
	Equals            types.List              `tfsdk:"equals"`
	AllEqualsRequired types.Bool              `tfsdk:"all_equals_required"`
	TimeWindow        types.Int64             `tfsdk:"time_window"`
	IKeys             types.List              `tfsdk:"i_keys"`
	IScoreThreshold   types.Float64           `tfsdk:"i_score_threshold"`
	StormThresholds   types.List              `tfsdk:"storm_thresholds"`
}

type ChannelFlappingModel struct {
	IsDisabled types.Bool  `tfsdk:"is_disabled"`
	MaxChanges types.Int64 `tfsdk:"max_changes"`
	InMins     types.Int64 `tfsdk:"in_mins"`
	MuteMins   types.Int64 `tfsdk:"mute_mins"`
}

type ChannelResourceModel struct {
	ID                      types.String          `tfsdk:"id"`
	ChannelName             types.String          `tfsdk:"channel_name"`
	Description             types.String          `tfsdk:"description"`
	TeamID                  types.Int64           `tfsdk:"team_id"`
	ManagingTeamIDs         types.List            `tfsdk:"managing_team_ids"`
	IsPrivate               types.Bool            `tfsdk:"is_private"`
	AutoResolveTimeout      types.Int64           `tfsdk:"auto_resolve_timeout"`
	AutoResolveMode         types.String          `tfsdk:"auto_resolve_mode"`
	DisableOutlierDetection types.Bool            `tfsdk:"disable_outlier_detection"`
	DisableAutoClose        types.Bool            `tfsdk:"disable_auto_close"`
	Group                   *ChannelGroupModel    `tfsdk:"group"`
	Flapping                *ChannelFlappingModel `tfsdk:"flapping"`
}

func (r *ChannelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

func (r *ChannelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Flashduty channel.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the channel.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"channel_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the channel.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the channel.",
			},
			"team_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the team that owns this channel.",
			},
			"managing_team_ids": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "IDs of managing teams (max 3). Managing teams take over edit permissions from the owning team.",
			},
			"is_private": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the channel is private. Defaults to `false`.",
			},
			"auto_resolve_timeout": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Auto-resolve timeout in seconds (0-86400). If not set, auto-resolve is disabled.",
			},
			"auto_resolve_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Auto-resolve mode. Valid values: `trigger`, `update`.",
			},
			"disable_outlier_detection": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to disable outlier detection. Defaults to `false`.",
			},
			"disable_auto_close": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to disable automatic incident closure. Defaults to `false`.",
			},
			"group": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Alert grouping configuration.",
				Attributes: map[string]schema.Attribute{
					"method": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Grouping method: `i` (intelligent), `p` (pattern/rule-based), `n` (none).",
					},
					"time_window": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Time window in minutes (0-60). 0 means merge until incident closes.",
					},
					"cases": schema.ListNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Branch grouping rules (max 10). Matched top-down; first match determines grouping dimensions.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"if": schema.ListNestedAttribute{
									Required:            true,
									MarkdownDescription: "Filter conditions (AND logic).",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "Filter key (e.g. `title`, `severity`, `labels.xxx`).",
											},
											"oper": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "Operator: `IN` or `NOTIN`.",
											},
											"vals": schema.ListAttribute{
												Required:            true,
												ElementType:         types.StringType,
												MarkdownDescription: "The values to match.",
											},
										},
									},
								},
								"equals": schema.ListAttribute{
									Required:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "Grouping dimensions (e.g. `title`, `severity`, `labels.xxx`).",
								},
							},
						},
					},
					"equals": schema.ListAttribute{
						Optional:            true,
						ElementType:         types.ListType{ElemType: types.StringType},
						MarkdownDescription: "Default grouping dimensions (OR of AND groups). Used when no case matches.",
					},
					"all_equals_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether all grouping dimensions must be present. Default `false`.",
					},
					"i_keys": schema.ListAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "Fields for intelligent similarity scoring (e.g. `title`, `description`, `labels.service`).",
					},
					"i_score_threshold": schema.Float64Attribute{
						Optional:            true,
						MarkdownDescription: "Similarity threshold for intelligent grouping (0.5-1.0, default 0.9).",
					},
					"storm_thresholds": schema.ListAttribute{
						Optional:            true,
						ElementType:         types.Int64Type,
						MarkdownDescription: "Alert storm warning thresholds (max 5 values, each 2-10000).",
					},
				},
			},
			"flapping": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Incident flap detection configuration.",
				Attributes: map[string]schema.Attribute{
					"is_disabled": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "Whether flap detection is disabled.",
					},
					"max_changes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Max state changes within the time window (2-100, default 4).",
					},
					"in_mins": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Statistics time window in minutes (1-1440, default 60).",
					},
					"mute_mins": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Mute window in minutes (0-1440, default 120). 0 means no muting.",
					},
				},
			},
		},
	}
}

func (r *ChannelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (r *ChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateChannelRequest{
		ChannelName:             data.ChannelName.ValueString(),
		Description:             data.Description.ValueString(),
		TeamID:                  data.TeamID.ValueInt64(),
		IsPrivate:               data.IsPrivate.ValueBoolPointer(),
		DisableOutlierDetection: data.DisableOutlierDetection.ValueBoolPointer(),
	}
	createReq.DisableAutoClose = data.DisableAutoClose.ValueBoolPointer()

	if !data.ManagingTeamIDs.IsNull() {
		var ids []int64
		resp.Diagnostics.Append(data.ManagingTeamIDs.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.ManagingTeamIDs = ids
	}
	if !data.AutoResolveTimeout.IsNull() {
		createReq.AutoResolveTimeout = int(data.AutoResolveTimeout.ValueInt64())
	}
	if !data.AutoResolveMode.IsNull() {
		createReq.AutoResolveMode = data.AutoResolveMode.ValueString()
	}

	createReq.Group = r.buildGroup(ctx, data.Group, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	createReq.Flapping = r.buildFlapping(data.Flapping)

	result, err := r.client.CreateChannel(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Channel", err.Error())
		return
	}

	data.ID = types.StringValue(strconv.FormatInt(result.ChannelID, 10))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channelID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Channel ID", err.Error())
		return
	}

	channel, err := r.client.GetChannel(ctx, &client.GetChannelRequest{
		ChannelID: channelID,
	})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Channel", err.Error())
		return
	}

	data.ChannelName = types.StringValue(channel.ChannelName)
	data.Description = types.StringValue(channel.Description)
	data.TeamID = types.Int64Value(channel.TeamID)
	data.IsPrivate = types.BoolValue(channel.IsPrivate)
	data.DisableOutlierDetection = types.BoolValue(channel.DisableOutlierDetection)
	data.DisableAutoClose = types.BoolValue(channel.DisableAutoClose)

	if len(channel.ManagingTeamIDs) > 0 {
		ids, d := types.ListValueFrom(ctx, types.Int64Type, channel.ManagingTeamIDs)
		resp.Diagnostics.Append(d...)
		data.ManagingTeamIDs = ids
	} else {
		data.ManagingTeamIDs = types.ListNull(types.Int64Type)
	}

	if !data.AutoResolveTimeout.IsNull() {
		data.AutoResolveTimeout = types.Int64Value(int64(channel.AutoResolveTimeout))
	}
	if !data.AutoResolveMode.IsNull() {
		data.AutoResolveMode = types.StringValue(channel.AutoResolveMode)
	}

	if data.Group != nil {
		data.Group = r.readGroup(ctx, channel.Group, &resp.Diagnostics)
	}
	if data.Flapping != nil {
		data.Flapping = r.readFlapping(channel.Flapping)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channelID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Channel ID", err.Error())
		return
	}

	isPrivate := data.IsPrivate.ValueBool()
	disableOutlier := data.DisableOutlierDetection.ValueBool()
	disableAutoClose := data.DisableAutoClose.ValueBool()

	updateReq := &client.UpdateChannelRequest{
		ChannelID:               channelID,
		ChannelName:             data.ChannelName.ValueString(),
		Description:             data.Description.ValueString(),
		TeamID:                  data.TeamID.ValueInt64(),
		IsPrivate:               &isPrivate,
		DisableOutlierDetection: &disableOutlier,
		DisableAutoClose:        &disableAutoClose,
	}

	if !data.ManagingTeamIDs.IsNull() {
		var ids []int64
		resp.Diagnostics.Append(data.ManagingTeamIDs.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.ManagingTeamIDs = ids
	}
	if !data.AutoResolveTimeout.IsNull() {
		timeout := int(data.AutoResolveTimeout.ValueInt64())
		updateReq.AutoResolveTimeout = &timeout
	}
	if !data.AutoResolveMode.IsNull() {
		updateReq.AutoResolveMode = data.AutoResolveMode.ValueString()
	}

	updateReq.Group = r.buildGroup(ctx, data.Group, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	updateReq.Flapping = r.buildFlapping(data.Flapping)

	err = r.client.UpdateChannel(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Channel", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channelID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Channel ID", err.Error())
		return
	}

	// Disable the channel first (required before deletion)
	err = r.client.DisableChannel(ctx, &client.ChannelStatusRequest{
		ChannelID: channelID,
	})
	if err != nil && !client.IsNotFoundError(err) {
		resp.Diagnostics.AddError("Error Disabling Channel", err.Error())
		return
	}

	// Now delete the channel
	err = r.client.DeleteChannel(ctx, &client.DeleteChannelRequest{
		ChannelID: channelID,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Channel", err.Error())
		return
	}
}

func (r *ChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ChannelResource) buildGroup(ctx context.Context, g *ChannelGroupModel, diags *diag.Diagnostics) *client.ChannelGroup {
	if g == nil {
		return nil
	}

	result := &client.ChannelGroup{
		Method: g.Method.ValueString(),
	}

	if !g.TimeWindow.IsNull() {
		result.TimeWindow = int(g.TimeWindow.ValueInt64())
	}
	if !g.AllEqualsRequired.IsNull() {
		result.AllEqualsRequired = g.AllEqualsRequired.ValueBool()
	}
	if !g.IScoreThreshold.IsNull() {
		result.IScoreThreshold = g.IScoreThreshold.ValueFloat64()
	}

	if !g.IKeys.IsNull() {
		var keys []string
		diags.Append(g.IKeys.ElementsAs(ctx, &keys, false)...)
		result.IKeys = keys
	}

	if !g.StormThresholds.IsNull() {
		var thresholds []int
		var int64s []int64
		diags.Append(g.StormThresholds.ElementsAs(ctx, &int64s, false)...)
		for _, v := range int64s {
			thresholds = append(thresholds, int(v))
		}
		result.StormThresholds = thresholds
	}

	if !g.Equals.IsNull() {
		var outerList []types.List
		diags.Append(g.Equals.ElementsAs(ctx, &outerList, false)...)
		for _, inner := range outerList {
			var strs []string
			diags.Append(inner.ElementsAs(ctx, &strs, false)...)
			result.Equals = append(result.Equals, strs)
		}
	}

	for _, c := range g.Cases {
		rule := client.GroupCaseRule{}
		if !c.Equals.IsNull() {
			var eqs []string
			diags.Append(c.Equals.ElementsAs(ctx, &eqs, false)...)
			rule.Equals = eqs
		}
		for _, f := range c.If {
			filter := client.Filter{
				Key:  f.Key.ValueString(),
				Oper: f.Oper.ValueString(),
			}
			if !f.Vals.IsNull() {
				var vals []string
				diags.Append(f.Vals.ElementsAs(ctx, &vals, false)...)
				filter.Vals = vals
			}
			rule.If = append(rule.If, filter)
		}
		result.Cases = append(result.Cases, rule)
	}

	return result
}

func (r *ChannelResource) readGroup(ctx context.Context, g *client.ChannelGroup, diags *diag.Diagnostics) *ChannelGroupModel {
	if g == nil {
		return nil
	}

	result := &ChannelGroupModel{
		Method: types.StringValue(g.Method),
	}

	if g.TimeWindow != 0 {
		result.TimeWindow = types.Int64Value(int64(g.TimeWindow))
	}
	if g.AllEqualsRequired {
		result.AllEqualsRequired = types.BoolValue(true)
	}
	if g.IScoreThreshold != 0 {
		result.IScoreThreshold = types.Float64Value(g.IScoreThreshold)
	}

	if len(g.IKeys) > 0 {
		v, d := types.ListValueFrom(ctx, types.StringType, g.IKeys)
		diags.Append(d...)
		result.IKeys = v
	} else {
		result.IKeys = types.ListNull(types.StringType)
	}

	if len(g.StormThresholds) > 0 {
		int64s := make([]int64, len(g.StormThresholds))
		for i, v := range g.StormThresholds {
			int64s[i] = int64(v)
		}
		v, d := types.ListValueFrom(ctx, types.Int64Type, int64s)
		diags.Append(d...)
		result.StormThresholds = v
	} else {
		result.StormThresholds = types.ListNull(types.Int64Type)
	}

	if len(g.Equals) > 0 {
		innerListType := types.ListType{ElemType: types.StringType}
		outerVals := make([]attr.Value, len(g.Equals))
		for i, inner := range g.Equals {
			v, d := types.ListValueFrom(ctx, types.StringType, inner)
			diags.Append(d...)
			outerVals[i] = v
		}
		outerList, d := types.ListValue(innerListType, outerVals)
		diags.Append(d...)
		result.Equals = outerList
	} else {
		result.Equals = types.ListNull(types.ListType{ElemType: types.StringType})
	}

	for _, c := range g.Cases {
		caseModel := ChannelGroupCaseModel{}
		if len(c.Equals) > 0 {
			v, d := types.ListValueFrom(ctx, types.StringType, c.Equals)
			diags.Append(d...)
			caseModel.Equals = v
		} else {
			caseModel.Equals = types.ListNull(types.StringType)
		}
		for _, f := range c.If {
			fm := ChannelGroupFilterModel{
				Key:  types.StringValue(f.Key),
				Oper: types.StringValue(f.Oper),
			}
			if len(f.Vals) > 0 {
				v, d := types.ListValueFrom(ctx, types.StringType, f.Vals)
				diags.Append(d...)
				fm.Vals = v
			} else {
				fm.Vals = types.ListNull(types.StringType)
			}
			caseModel.If = append(caseModel.If, fm)
		}
		result.Cases = append(result.Cases, caseModel)
	}

	return result
}

func (r *ChannelResource) buildFlapping(f *ChannelFlappingModel) *client.ChannelFlapping {
	if f == nil {
		return nil
	}
	result := &client.ChannelFlapping{
		IsDisabled: f.IsDisabled.ValueBool(),
	}
	if !f.MaxChanges.IsNull() {
		result.MaxChanges = int(f.MaxChanges.ValueInt64())
	}
	if !f.InMins.IsNull() {
		result.InMins = int(f.InMins.ValueInt64())
	}
	if !f.MuteMins.IsNull() {
		result.MuteMins = int(f.MuteMins.ValueInt64())
	}
	return result
}

func (r *ChannelResource) readFlapping(f *client.ChannelFlapping) *ChannelFlappingModel {
	if f == nil {
		return nil
	}
	return &ChannelFlappingModel{
		IsDisabled: types.BoolValue(f.IsDisabled),
		MaxChanges: types.Int64Value(int64(f.MaxChanges)),
		InMins:     types.Int64Value(int64(f.InMins)),
		MuteMins:   types.Int64Value(int64(f.MuteMins)),
	}
}
