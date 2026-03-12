package provider

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ScheduleResource{}
	_ resource.ResourceWithImportState = &ScheduleResource{}
)

func NewScheduleResource() resource.Resource {
	return &ScheduleResource{}
}

type ScheduleResource struct {
	client *client.Client
}

type ScheduleLayerGroupMemberModel struct {
	RoleID    types.Int64 `tfsdk:"role_id"`
	PersonIDs types.List  `tfsdk:"person_ids"`
}

type ScheduleLayerGroupModel struct {
	GroupName types.String                    `tfsdk:"group_name"`
	Members   []ScheduleLayerGroupMemberModel `tfsdk:"members"`
}

type ScheduleRestrictPeriodModel struct {
	RestrictStart types.Int64 `tfsdk:"restrict_start"`
	RestrictEnd   types.Int64 `tfsdk:"restrict_end"`
}

type ScheduleDayMaskModel struct {
	Repeat types.List `tfsdk:"repeat"`
}

type ScheduleLayerModel struct {
	LayerName             types.String                  `tfsdk:"layer_name"`
	Mode                  types.Int64                   `tfsdk:"mode"`
	Groups                []ScheduleLayerGroupModel     `tfsdk:"groups"`
	FairRotation          types.Bool                    `tfsdk:"fair_rotation"`
	HandoffTime           types.Int64                   `tfsdk:"handoff_time"`
	LayerStart            types.Int64                   `tfsdk:"layer_start"`
	LayerEnd              types.Int64                   `tfsdk:"layer_end"`
	RestrictMode          types.Int64                   `tfsdk:"restrict_mode"`
	RestrictPeriods       []ScheduleRestrictPeriodModel `tfsdk:"restrict_periods"`
	DayMask               *ScheduleDayMaskModel         `tfsdk:"day_mask"`
	MaskContinuousEnabled types.Bool                    `tfsdk:"mask_continuous_enabled"`
	RotationUnit          types.String                  `tfsdk:"rotation_unit"`
	RotationValue         types.Int64                   `tfsdk:"rotation_value"`
}

type ScheduleNotifyFixedTimeModel struct {
	Cycle types.String `tfsdk:"cycle"`
	Start types.String `tfsdk:"start"`
}

type ScheduleNotifyByModel struct {
	FollowPreference types.Bool `tfsdk:"follow_preference"`
	PersonalChannels types.List `tfsdk:"personal_channels"`
}

type ScheduleNotifyWebhookModel struct {
	Type     types.String `tfsdk:"type"`
	Settings types.String `tfsdk:"settings"`
}

type ScheduleNotifyModel struct {
	AdvanceInTime types.Int64                   `tfsdk:"advance_in_time"`
	FixedTime     *ScheduleNotifyFixedTimeModel `tfsdk:"fixed_time"`
	By            *ScheduleNotifyByModel        `tfsdk:"by"`
	Webhooks      []ScheduleNotifyWebhookModel  `tfsdk:"webhooks"`
}

type ScheduleResourceModel struct {
	ID           types.String         `tfsdk:"id"`
	ScheduleName types.String         `tfsdk:"schedule_name"`
	Description  types.String         `tfsdk:"description"`
	TeamID       types.Int64          `tfsdk:"team_id"`
	Status       types.Int64          `tfsdk:"status"`
	Layers       []ScheduleLayerModel `tfsdk:"layers"`
	Notify       *ScheduleNotifyModel `tfsdk:"notify"`
}

func (r *ScheduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schedule"
}

func (r *ScheduleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Flashduty on-call schedule.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the schedule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the schedule.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The description of the schedule.",
			},
			"team_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The ID of the team that owns this schedule.",
			},
			"status": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "The status of the schedule. 0 = enabled.",
			},
			"layers": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "The schedule layers (rotation rules).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"layer_name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The name of the layer.",
						},
						"mode": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "The mode of the layer. 0 = normal, 1 = temporary.",
						},
						"layer_start": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "The start time of the layer (Unix timestamp).",
						},
						"layer_end": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "The end time of the layer (Unix timestamp).",
						},
						"rotation_unit": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The rotation unit. Valid values: `hour`, `day`, `week`, `month`.",
						},
						"rotation_value": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "The rotation value.",
						},
						"fair_rotation": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether to enable fair rotation.",
						},
						"handoff_time": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "The handoff time in seconds.",
						},
						"restrict_mode": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "The restriction mode. 0 = none.",
						},
						"restrict_periods": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Time restriction periods.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"restrict_start": schema.Int64Attribute{
										Required:            true,
										MarkdownDescription: "Restriction start time (seconds from midnight).",
									},
									"restrict_end": schema.Int64Attribute{
										Required:            true,
										MarkdownDescription: "Restriction end time (seconds from midnight).",
									},
								},
							},
						},
						"day_mask": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Day of week mask.",
							Attributes: map[string]schema.Attribute{
								"repeat": schema.ListAttribute{
									Required:            true,
									ElementType:         types.Int64Type,
									MarkdownDescription: "Days of week (0=Sunday, 1=Monday, ..., 6=Saturday).",
								},
							},
						},
						"mask_continuous_enabled": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether continuous mask is enabled.",
						},
						"groups": schema.ListNestedAttribute{
							Required:            true,
							MarkdownDescription: "The on-call groups for this layer.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"group_name": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "The name of the group.",
									},
									"members": schema.ListNestedAttribute{
										Required:            true,
										MarkdownDescription: "The members of the group.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"role_id": schema.Int64Attribute{
													Optional:            true,
													Computed:            true,
													Default:             int64default.StaticInt64(0),
													MarkdownDescription: "The role ID.",
												},
												"person_ids": schema.ListAttribute{
													Required:            true,
													ElementType:         types.Int64Type,
													MarkdownDescription: "The list of person IDs.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"notify": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Notification settings for the schedule.",
				Attributes: map[string]schema.Attribute{
					"advance_in_time": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Advance notification time in seconds.",
					},
					"fixed_time": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Fixed time notification settings.",
						Attributes: map[string]schema.Attribute{
							"cycle": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Notification cycle (e.g., `day`).",
							},
							"start": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Start time string.",
							},
						},
					},
					"by": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Notification channel settings.",
						Attributes: map[string]schema.Attribute{
							"follow_preference": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether to follow user preference.",
							},
							"personal_channels": schema.ListAttribute{
								Optional:            true,
								ElementType:         types.StringType,
								MarkdownDescription: "Personal notification channels (e.g., `email`, `sms`, `phone`).",
							},
						},
					},
					"webhooks": schema.ListNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Webhook notification settings.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Webhook type (e.g., `feishu`, `dingtalk`, `slack`).",
								},
								"settings": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Webhook settings as JSON string. Use `jsonencode()` to construct.",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *ScheduleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (r *ScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	layers := r.buildLayers(ctx, data.Layers, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateScheduleRequest{
		ScheduleName: data.ScheduleName.ValueString(),
		Description:  data.Description.ValueString(),
		Status:       int(data.Status.ValueInt64()),
		Layers:       layers,
	}

	if !data.TeamID.IsNull() {
		createReq.TeamID = data.TeamID.ValueInt64()
	}

	createReq.Notify = r.buildNotify(ctx, data.Notify, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateSchedule(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Schedule", err.Error())
		return
	}

	data.ID = types.StringValue(strconv.FormatInt(result.ScheduleID, 10))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scheduleID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Schedule ID", err.Error())
		return
	}

	now := time.Now().Unix()
	schedule, err := r.client.GetSchedule(ctx, &client.GetScheduleRequest{
		ScheduleID: scheduleID,
		Start:      now,
		End:        now + 7*24*3600,
	})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Schedule", err.Error())
		return
	}

	data.ScheduleName = types.StringValue(schedule.ScheduleName)
	data.Description = types.StringValue(schedule.Description)
	data.Status = types.Int64Value(int64(schedule.Status))

	if schedule.TeamID != 0 {
		data.TeamID = types.Int64Value(schedule.TeamID)
	}

	data.Layers = r.readLayers(ctx, schedule.Layers, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Notify = r.readNotify(ctx, schedule.Notify, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scheduleID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Schedule ID", err.Error())
		return
	}

	layers := r.buildLayers(ctx, data.Layers, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	status := int(data.Status.ValueInt64())
	updateReq := &client.UpdateScheduleRequest{
		ScheduleID:   scheduleID,
		ScheduleName: data.ScheduleName.ValueString(),
		Description:  data.Description.ValueString(),
		Status:       &status,
		Layers:       layers,
	}

	if !data.TeamID.IsNull() {
		updateReq.TeamID = data.TeamID.ValueInt64()
	}

	updateReq.Notify = r.buildNotify(ctx, data.Notify, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err = r.client.UpdateSchedule(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Schedule", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scheduleID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing Schedule ID", err.Error())
		return
	}

	err = r.client.DeleteSchedule(ctx, &client.DeleteScheduleRequest{
		ScheduleIDs: []int64{scheduleID},
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Schedule", err.Error())
		return
	}
}

func (r *ScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ScheduleResource) buildLayers(ctx context.Context, layers []ScheduleLayerModel, diags interface{ Append(...diag.Diagnostic) }) []client.ScheduleLayer {
	var result []client.ScheduleLayer

	for _, layer := range layers {
		l := client.ScheduleLayer{
			LayerName:     layer.LayerName.ValueString(),
			Mode:          int(layer.Mode.ValueInt64()),
			LayerStart:    layer.LayerStart.ValueInt64(),
			RotationUnit:  layer.RotationUnit.ValueString(),
			RotationValue: int(layer.RotationValue.ValueInt64()),
		}

		if !layer.LayerEnd.IsNull() {
			l.LayerEnd = layer.LayerEnd.ValueInt64()
		}
		if !layer.FairRotation.IsNull() {
			l.FairRotation = layer.FairRotation.ValueBool()
		}
		if !layer.HandoffTime.IsNull() {
			l.HandoffTime = int(layer.HandoffTime.ValueInt64())
		}
		if !layer.RestrictMode.IsNull() {
			l.RestrictMode = int(layer.RestrictMode.ValueInt64())
		}
		if !layer.MaskContinuousEnabled.IsNull() {
			l.MaskContinuousEnabled = layer.MaskContinuousEnabled.ValueBool()
		}

		for _, rp := range layer.RestrictPeriods {
			l.RestrictPeriods = append(l.RestrictPeriods, client.RestrictPeriod{
				RestrictStart: int(rp.RestrictStart.ValueInt64()),
				RestrictEnd:   int(rp.RestrictEnd.ValueInt64()),
			})
		}

		if layer.DayMask != nil {
			var repeat []int
			diags.Append(layer.DayMask.Repeat.ElementsAs(ctx, &repeat, false)...)
			l.DayMask = &client.DayMask{Repeat: repeat}
		}

		for _, group := range layer.Groups {
			g := client.ScheduleGroup{
				GroupName: group.GroupName.ValueString(),
			}

			for _, member := range group.Members {
				m := client.ScheduleGroupMember{
					RoleID: member.RoleID.ValueInt64(),
				}

				var personIDs []int64
				diags.Append(member.PersonIDs.ElementsAs(ctx, &personIDs, false)...)
				m.PersonIDs = personIDs

				g.Members = append(g.Members, m)
			}

			l.Groups = append(l.Groups, g)
		}

		result = append(result, l)
	}

	return result
}

func (r *ScheduleResource) readLayers(_ context.Context, layers []client.ScheduleLayer, diags *diag.Diagnostics) []ScheduleLayerModel {
	var result []ScheduleLayerModel

	for _, layer := range layers {
		l := ScheduleLayerModel{
			LayerName:     types.StringValue(layer.LayerName),
			Mode:          types.Int64Value(int64(layer.Mode)),
			LayerStart:    types.Int64Value(layer.LayerStart),
			RotationUnit:  types.StringValue(layer.RotationUnit),
			RotationValue: types.Int64Value(int64(layer.RotationValue)),
		}

		if layer.LayerEnd != 0 {
			l.LayerEnd = types.Int64Value(layer.LayerEnd)
		}
		if layer.FairRotation {
			l.FairRotation = types.BoolValue(true)
		}
		if layer.HandoffTime != 0 {
			l.HandoffTime = types.Int64Value(int64(layer.HandoffTime))
		}
		if layer.RestrictMode != 0 {
			l.RestrictMode = types.Int64Value(int64(layer.RestrictMode))
		}
		if layer.MaskContinuousEnabled {
			l.MaskContinuousEnabled = types.BoolValue(true)
		}

		for _, rp := range layer.RestrictPeriods {
			l.RestrictPeriods = append(l.RestrictPeriods, ScheduleRestrictPeriodModel{
				RestrictStart: types.Int64Value(int64(rp.RestrictStart)),
				RestrictEnd:   types.Int64Value(int64(rp.RestrictEnd)),
			})
		}

		if layer.DayMask != nil && len(layer.DayMask.Repeat) > 0 {
			repeatVals := make([]attr.Value, len(layer.DayMask.Repeat))
			for i, v := range layer.DayMask.Repeat {
				repeatVals[i] = types.Int64Value(int64(v))
			}
			repeatList, d := types.ListValue(types.Int64Type, repeatVals)
			diags.Append(d...)
			l.DayMask = &ScheduleDayMaskModel{Repeat: repeatList}
		}

		for _, group := range layer.Groups {
			g := ScheduleLayerGroupModel{}
			if group.GroupName != "" {
				g.GroupName = types.StringValue(group.GroupName)
			}

			for _, member := range group.Members {
				m := ScheduleLayerGroupMemberModel{
					RoleID: types.Int64Value(member.RoleID),
				}

				personIDsValues := make([]attr.Value, len(member.PersonIDs))
				for i, id := range member.PersonIDs {
					personIDsValues[i] = types.Int64Value(id)
				}
				var d diag.Diagnostics
				m.PersonIDs, d = types.ListValue(types.Int64Type, personIDsValues)
				diags.Append(d...)

				g.Members = append(g.Members, m)
			}

			l.Groups = append(l.Groups, g)
		}

		result = append(result, l)
	}

	return result
}

func (r *ScheduleResource) buildNotify(ctx context.Context, notify *ScheduleNotifyModel, diags *diag.Diagnostics) *client.ScheduleNotify {
	if notify == nil {
		return nil
	}

	result := &client.ScheduleNotify{}

	if !notify.AdvanceInTime.IsNull() {
		result.AdvanceInTime = int(notify.AdvanceInTime.ValueInt64())
	}

	if notify.FixedTime != nil {
		result.FixedTime = &client.NotifyFixedTime{
			Cycle: notify.FixedTime.Cycle.ValueString(),
			Start: notify.FixedTime.Start.ValueString(),
		}
	}

	if notify.By != nil {
		result.By = &client.NotifyBy{
			FollowPreference: notify.By.FollowPreference.ValueBool(),
		}
		if !notify.By.PersonalChannels.IsNull() {
			var channels []string
			diags.Append(notify.By.PersonalChannels.ElementsAs(ctx, &channels, false)...)
			result.By.PersonalChannels = channels
		}
	}

	for _, wh := range notify.Webhooks {
		webhook := client.NotifyWebhook{
			Type: wh.Type.ValueString(),
		}
		if !wh.Settings.IsNull() && wh.Settings.ValueString() != "" {
			var settings map[string]interface{}
			if err := json.Unmarshal([]byte(wh.Settings.ValueString()), &settings); err != nil {
				diags.AddError("Invalid Webhook Settings JSON", err.Error())
				return nil
			}
			webhook.Settings = settings
		}
		result.Webhooks = append(result.Webhooks, webhook)
	}

	return result
}

func (r *ScheduleResource) readNotify(_ context.Context, notify *client.ScheduleNotify, diags *diag.Diagnostics) *ScheduleNotifyModel {
	if notify == nil {
		return nil
	}

	result := &ScheduleNotifyModel{}

	if notify.AdvanceInTime != 0 {
		result.AdvanceInTime = types.Int64Value(int64(notify.AdvanceInTime))
	}

	if notify.FixedTime != nil {
		result.FixedTime = &ScheduleNotifyFixedTimeModel{
			Cycle: types.StringValue(notify.FixedTime.Cycle),
			Start: types.StringValue(notify.FixedTime.Start),
		}
	}

	if notify.By != nil {
		result.By = &ScheduleNotifyByModel{
			FollowPreference: types.BoolValue(notify.By.FollowPreference),
		}
		if len(notify.By.PersonalChannels) > 0 {
			vals := make([]attr.Value, len(notify.By.PersonalChannels))
			for i, ch := range notify.By.PersonalChannels {
				vals[i] = types.StringValue(ch)
			}
			var d diag.Diagnostics
			result.By.PersonalChannels, d = types.ListValue(types.StringType, vals)
			diags.Append(d...)
		}
	}

	for _, wh := range notify.Webhooks {
		whModel := ScheduleNotifyWebhookModel{
			Type: types.StringValue(wh.Type),
		}
		if wh.Settings != nil {
			settingsJSON, err := json.Marshal(wh.Settings)
			if err != nil {
				diags.AddError("Error Serializing Webhook Settings", err.Error())
				return nil
			}
			whModel.Settings = types.StringValue(string(settingsJSON))
		}
		result.Webhooks = append(result.Webhooks, whModel)
	}

	if result.AdvanceInTime.IsNull() && result.FixedTime == nil && result.By == nil && len(result.Webhooks) == 0 {
		return nil
	}

	return result
}
