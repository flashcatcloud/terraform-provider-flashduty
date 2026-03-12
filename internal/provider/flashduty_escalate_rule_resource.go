package provider

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &EscalateRuleResource{}
	_ resource.ResourceWithImportState = &EscalateRuleResource{}
)

func NewEscalateRuleResource() resource.Resource {
	return &EscalateRuleResource{}
}

type EscalateRuleResource struct {
	client *client.Client
}

type EscalateTargetByModel struct {
	FollowPreference types.Bool `tfsdk:"follow_preference"`
	Critical         types.List `tfsdk:"critical"`
	Warning          types.List `tfsdk:"warning"`
	Info             types.List `tfsdk:"info"`
}

type EscalateWebhookModel struct {
	Type     types.String `tfsdk:"type"`
	Settings types.String `tfsdk:"settings"`
}

type EscalateTargetModel struct {
	PersonIDs   types.List             `tfsdk:"person_ids"`
	TeamIDs     types.List             `tfsdk:"team_ids"`
	ScheduleIDs types.List             `tfsdk:"schedule_ids"`
	Emails      types.List             `tfsdk:"emails"`
	By          *EscalateTargetByModel `tfsdk:"by"`
	Webhooks    []EscalateWebhookModel `tfsdk:"webhooks"`
}

type EscalateTimeFilterModel struct {
	Start  types.String `tfsdk:"start"`
	End    types.String `tfsdk:"end"`
	Repeat types.List   `tfsdk:"repeat"`
	CalID  types.String `tfsdk:"cal_id"`
	IsOff  types.Bool   `tfsdk:"is_off"`
}

type EscalateFilterModel struct {
	Key  types.String `tfsdk:"key"`
	Oper types.String `tfsdk:"oper"`
	Vals types.List   `tfsdk:"vals"`
}

type EscalateFilterGroupModel struct {
	Conditions []EscalateFilterModel `tfsdk:"conditions"`
}

type EscalateLayerModel struct {
	MaxTimes       types.Int64          `tfsdk:"max_times"`
	NotifyStep     types.Float64        `tfsdk:"notify_step"`
	EscalateWindow types.Int64          `tfsdk:"escalate_window"`
	ForceEscalate  types.Bool           `tfsdk:"force_escalate"`
	Target         *EscalateTargetModel `tfsdk:"target"`
}

type EscalateRuleResourceModel struct {
	ID          types.String               `tfsdk:"id"`
	ChannelID   types.Int64                `tfsdk:"channel_id"`
	RuleName    types.String               `tfsdk:"rule_name"`
	Description types.String               `tfsdk:"description"`
	TemplateID  types.String               `tfsdk:"template_id"`
	AggrWindow  types.Int64                `tfsdk:"aggr_window"`
	Priority    types.Int64                `tfsdk:"priority"`
	Layers      []EscalateLayerModel       `tfsdk:"layers"`
	TimeFilters []EscalateTimeFilterModel  `tfsdk:"time_filters"`
	Filters     []EscalateFilterGroupModel `tfsdk:"filters"`
}

func (r *EscalateRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_escalate_rule"
}

func (r *EscalateRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Flashduty escalation rule for a channel.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the escalate rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"channel_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the channel this rule belongs to.",
			},
			"rule_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the escalation rule.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the rule.",
			},
			"template_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The notification template ID.",
			},
			"aggr_window": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "Aggregation window in seconds (0-3600).",
			},
			"priority": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "The priority of the escalation rule for ordering. Defaults to `1`.",
			},
			"layers": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "The escalation layers.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"max_times": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							Default:             int64default.StaticInt64(2),
							MarkdownDescription: "Maximum notification times. Defaults to 2.",
						},
						"notify_step": schema.Float64Attribute{
							Optional:            true,
							Computed:            true,
							Default:             float64default.StaticFloat64(10),
							MarkdownDescription: "Notification interval in minutes. Defaults to 10.",
						},
						"escalate_window": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							Default:             int64default.StaticInt64(30),
							MarkdownDescription: "Minutes before escalating to next layer. Defaults to 30.",
						},
						"force_escalate": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether to escalate even after acknowledgment.",
						},
						"target": schema.SingleNestedAttribute{
							Required:            true,
							MarkdownDescription: "The escalation target.",
							Attributes: map[string]schema.Attribute{
								"person_ids": schema.ListAttribute{
									Optional:            true,
									ElementType:         types.Int64Type,
									MarkdownDescription: "List of person IDs.",
								},
								"team_ids": schema.ListAttribute{
									Optional:            true,
									ElementType:         types.Int64Type,
									MarkdownDescription: "List of team IDs.",
								},
								"schedule_ids": schema.ListAttribute{
									Optional:            true,
									ElementType:         types.Int64Type,
									MarkdownDescription: "List of schedule IDs. Maps to `schedule_to_role_ids` in API with default role `[0]`.",
								},
								"emails": schema.ListAttribute{
									Optional:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "List of email addresses for dispatch.",
								},
								"by": schema.SingleNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Notification preference settings.",
									Attributes: map[string]schema.Attribute{
										"follow_preference": schema.BoolAttribute{
											Optional:            true,
											MarkdownDescription: "Whether to follow user's notification preferences.",
										},
										"critical": schema.ListAttribute{
											Optional:            true,
											ElementType:         types.StringType,
											MarkdownDescription: "Notification channels for Critical severity (e.g., `email`, `sms`, `phone`, `im`).",
										},
										"warning": schema.ListAttribute{
											Optional:            true,
											ElementType:         types.StringType,
											MarkdownDescription: "Notification channels for Warning severity.",
										},
										"info": schema.ListAttribute{
											Optional:            true,
											ElementType:         types.StringType,
											MarkdownDescription: "Notification channels for Info severity.",
										},
									},
								},
								"webhooks": schema.ListNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Webhook notification targets.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"type": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "Webhook type (e.g., `feishu_app`, `wecom`, `dingtalk`, `slack`).",
											},
											"settings": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "Webhook settings as JSON string.",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"time_filters": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Time-based filter conditions for when this rule applies.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Start time in HH:MM format.",
						},
						"end": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "End time in HH:MM format.",
						},
						"repeat": schema.ListAttribute{
							Optional:            true,
							ElementType:         types.Int64Type,
							MarkdownDescription: "Days of the week (0=Sunday, 1=Monday, ..., 6=Saturday).",
						},
						"cal_id": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Calendar ID for calendar-based filtering.",
						},
						"is_off": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether to match off-duty time from the calendar.",
						},
					},
				},
			},
			"filters": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Alert matching filter conditions (OR between groups, AND within conditions).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"conditions": schema.ListNestedAttribute{
							Required:            true,
							MarkdownDescription: "Filter conditions (AND logic).",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The attribute key to filter on (e.g., `title`, `severity`, `labels.xxx`).",
									},
									"oper": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The operator (`IN`, `NOTIN`).",
									},
									"vals": schema.ListAttribute{
										Required:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "The values to match.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *EscalateRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (r *EscalateRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EscalateRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	layers := r.buildLayers(ctx, data.Layers, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateEscalateRuleRequest{
		ChannelID:   data.ChannelID.ValueInt64(),
		RuleName:    data.RuleName.ValueString(),
		TemplateID:  data.TemplateID.ValueString(),
		AggrWindow:  int(data.AggrWindow.ValueInt64()),
		Priority:    int(data.Priority.ValueInt64()),
		Layers:      layers,
		TimeFilters: r.buildTimeFilters(ctx, data.TimeFilters),
		Filters:     r.buildFilters(ctx, data.Filters, &resp.Diagnostics),
	}
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Description.IsNull() {
		createReq.Description = data.Description.ValueString()
	}

	result, err := r.client.CreateEscalateRule(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Escalate Rule", err.Error())
		return
	}

	data.ID = types.StringValue(result.RuleID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EscalateRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EscalateRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetEscalateRule(ctx, &client.GetEscalateRuleRequest{
		ChannelID: data.ChannelID.ValueInt64(),
		RuleID:    data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Escalate Rule", err.Error())
		return
	}

	if rule == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.RuleName = types.StringValue(rule.RuleName)
	data.Description = types.StringValue(rule.Description)
	data.TemplateID = types.StringValue(rule.TemplateID)
	data.AggrWindow = types.Int64Value(int64(rule.AggrWindow))
	data.Priority = types.Int64Value(int64(rule.Priority))
	data.Layers = r.readLayers(ctx, rule.Layers, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.TimeFilters = r.readTimeFilters(ctx, rule.TimeFilters, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Filters = r.readFilters(ctx, rule.Filters, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EscalateRuleResource) readLayers(ctx context.Context, layers []client.EscalateLayer, diags *diag.Diagnostics) []EscalateLayerModel {
	var result []EscalateLayerModel

	for _, layer := range layers {
		l := EscalateLayerModel{
			MaxTimes:       types.Int64Value(int64(layer.MaxTimes)),
			NotifyStep:     types.Float64Value(layer.NotifyStep),
			EscalateWindow: types.Int64Value(int64(layer.EscalateWindow)),
			ForceEscalate:  types.BoolValue(layer.ForceEscalate),
		}

		if layer.Target != nil {
			target := &EscalateTargetModel{}

			if len(layer.Target.PersonIDs) > 0 {
				vals, d := types.ListValueFrom(ctx, types.Int64Type, layer.Target.PersonIDs)
				diags.Append(d...)
				target.PersonIDs = vals
			} else {
				target.PersonIDs = types.ListNull(types.Int64Type)
			}

			if len(layer.Target.TeamIDs) > 0 {
				vals, d := types.ListValueFrom(ctx, types.Int64Type, layer.Target.TeamIDs)
				diags.Append(d...)
				target.TeamIDs = vals
			} else {
				target.TeamIDs = types.ListNull(types.Int64Type)
			}

			if len(layer.Target.ScheduleToRoleIDs) > 0 {
				var scheduleIDs []int64
				for sid := range layer.Target.ScheduleToRoleIDs {
					if id, err := strconv.ParseInt(sid, 10, 64); err == nil {
						scheduleIDs = append(scheduleIDs, id)
					}
				}
				if len(scheduleIDs) > 0 {
					vals, d := types.ListValueFrom(ctx, types.Int64Type, scheduleIDs)
					diags.Append(d...)
					target.ScheduleIDs = vals
				} else {
					target.ScheduleIDs = types.ListNull(types.Int64Type)
				}
			} else {
				target.ScheduleIDs = types.ListNull(types.Int64Type)
			}

			if len(layer.Target.Emails) > 0 {
				vals, d := types.ListValueFrom(ctx, types.StringType, layer.Target.Emails)
				diags.Append(d...)
				target.Emails = vals
			} else {
				target.Emails = types.ListNull(types.StringType)
			}

			if layer.Target.By != nil {
				by := &EscalateTargetByModel{
					FollowPreference: types.BoolValue(layer.Target.By.FollowPreference),
				}
				if len(layer.Target.By.Critical) > 0 {
					vals, d := types.ListValueFrom(ctx, types.StringType, layer.Target.By.Critical)
					diags.Append(d...)
					by.Critical = vals
				} else {
					by.Critical = types.ListNull(types.StringType)
				}
				if len(layer.Target.By.Warning) > 0 {
					vals, d := types.ListValueFrom(ctx, types.StringType, layer.Target.By.Warning)
					diags.Append(d...)
					by.Warning = vals
				} else {
					by.Warning = types.ListNull(types.StringType)
				}
				if len(layer.Target.By.Info) > 0 {
					vals, d := types.ListValueFrom(ctx, types.StringType, layer.Target.By.Info)
					diags.Append(d...)
					by.Info = vals
				} else {
					by.Info = types.ListNull(types.StringType)
				}
				target.By = by
			}

			if len(layer.Target.Webhooks) > 0 {
				for _, wh := range layer.Target.Webhooks {
					settingsJSON, err := json.Marshal(wh.Settings)
					if err != nil {
						diags.AddError("Error Serializing Webhook Settings", err.Error())
						return nil
					}
					target.Webhooks = append(target.Webhooks, EscalateWebhookModel{
						Type:     types.StringValue(wh.Type),
						Settings: types.StringValue(string(settingsJSON)),
					})
				}
			}

			l.Target = target
		}

		result = append(result, l)
	}

	return result
}

func (r *EscalateRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EscalateRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	layers := r.buildLayers(ctx, data.Layers, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	aggrWindow := int(data.AggrWindow.ValueInt64())
	priority := int(data.Priority.ValueInt64())
	updateReq := &client.UpdateEscalateRuleRequest{
		ChannelID:   data.ChannelID.ValueInt64(),
		RuleID:      data.ID.ValueString(),
		RuleName:    data.RuleName.ValueString(),
		TemplateID:  data.TemplateID.ValueString(),
		AggrWindow:  &aggrWindow,
		Priority:    &priority,
		Layers:      layers,
		TimeFilters: r.buildTimeFilters(ctx, data.TimeFilters),
		Filters:     r.buildFilters(ctx, data.Filters, &resp.Diagnostics),
	}
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Description.IsNull() {
		updateReq.Description = data.Description.ValueString()
	}

	err := r.client.UpdateEscalateRule(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Escalate Rule", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EscalateRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EscalateRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Disable the rule first (required before deletion)
	err := r.client.DisableEscalateRule(ctx, &client.EnableEscalateRuleRequest{
		ChannelID: data.ChannelID.ValueInt64(),
		RuleID:    data.ID.ValueString(),
	})
	if err != nil && !client.IsNotFoundError(err) {
		resp.Diagnostics.AddError("Error Disabling Escalate Rule", err.Error())
		return
	}

	// Now delete the rule
	err = r.client.DeleteEscalateRule(ctx, &client.DeleteEscalateRuleRequest{
		ChannelID: data.ChannelID.ValueInt64(),
		RuleID:    data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Escalate Rule", err.Error())
		return
	}
}

func (r *EscalateRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Expected import ID format: channel_id/rule_id. Got: "+req.ID,
		)
		return
	}

	channelID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", "channel_id must be a valid integer: "+parts[0])
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("channel_id"), types.Int64Value(channelID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(parts[1]))...)
}

func (r *EscalateRuleResource) buildLayers(ctx context.Context, layers []EscalateLayerModel, diags *diag.Diagnostics) []client.EscalateLayer {
	var result []client.EscalateLayer

	for _, layer := range layers {
		l := client.EscalateLayer{
			MaxTimes:       int(layer.MaxTimes.ValueInt64()),
			NotifyStep:     layer.NotifyStep.ValueFloat64(),
			EscalateWindow: int(layer.EscalateWindow.ValueInt64()),
			ForceEscalate:  layer.ForceEscalate.ValueBool(),
		}

		if layer.Target != nil {
			target := &client.EscalateTarget{}

			if !layer.Target.PersonIDs.IsNull() {
				var ids []int64
				diags.Append(layer.Target.PersonIDs.ElementsAs(ctx, &ids, false)...)
				target.PersonIDs = ids
			}
			if !layer.Target.TeamIDs.IsNull() {
				var ids []int64
				diags.Append(layer.Target.TeamIDs.ElementsAs(ctx, &ids, false)...)
				target.TeamIDs = ids
			}
			if !layer.Target.ScheduleIDs.IsNull() {
				var ids []int64
				diags.Append(layer.Target.ScheduleIDs.ElementsAs(ctx, &ids, false)...)
				if len(ids) > 0 {
					m := make(map[string][]int64, len(ids))
					for _, id := range ids {
						m[strconv.FormatInt(id, 10)] = []int64{0}
					}
					target.ScheduleToRoleIDs = m
				}
			}
			if !layer.Target.Emails.IsNull() {
				var emails []string
				diags.Append(layer.Target.Emails.ElementsAs(ctx, &emails, false)...)
				target.Emails = emails
			}
			if layer.Target.By != nil {
				by := &client.EscalateTargetBy{
					FollowPreference: layer.Target.By.FollowPreference.ValueBool(),
				}
				if !layer.Target.By.Critical.IsNull() {
					var channels []string
					diags.Append(layer.Target.By.Critical.ElementsAs(ctx, &channels, false)...)
					by.Critical = channels
				}
				if !layer.Target.By.Warning.IsNull() {
					var channels []string
					diags.Append(layer.Target.By.Warning.ElementsAs(ctx, &channels, false)...)
					by.Warning = channels
				}
				if !layer.Target.By.Info.IsNull() {
					var channels []string
					diags.Append(layer.Target.By.Info.ElementsAs(ctx, &channels, false)...)
					by.Info = channels
				}
				target.By = by
			}

			for _, wh := range layer.Target.Webhooks {
				var settings map[string]interface{}
				if !wh.Settings.IsNull() {
					if err := json.Unmarshal([]byte(wh.Settings.ValueString()), &settings); err != nil {
						diags.AddWarning("Invalid Webhook Settings JSON",
							"Could not parse webhook settings: "+err.Error())
					}
				}
				target.Webhooks = append(target.Webhooks, client.EscalateWebhook{
					Type:     wh.Type.ValueString(),
					Settings: settings,
				})
			}

			l.Target = target
		}

		result = append(result, l)
	}

	return result
}

func (r *EscalateRuleResource) buildTimeFilters(ctx context.Context, tfs []EscalateTimeFilterModel) []client.TimeFilter {
	if len(tfs) == 0 {
		return nil
	}
	var result []client.TimeFilter
	for _, tf := range tfs {
		f := client.TimeFilter{
			Start: tf.Start.ValueString(),
			End:   tf.End.ValueString(),
			CalID: tf.CalID.ValueString(),
			IsOff: tf.IsOff.ValueBool(),
		}
		if !tf.Repeat.IsNull() {
			var int64s []int64
			tf.Repeat.ElementsAs(ctx, &int64s, false)
			repeat := make([]int, len(int64s))
			for i, v := range int64s {
				repeat[i] = int(v)
			}
			f.Repeat = repeat
		}
		result = append(result, f)
	}
	return result
}

func (r *EscalateRuleResource) buildFilters(ctx context.Context, groups []EscalateFilterGroupModel, diags *diag.Diagnostics) [][]client.Filter {
	if len(groups) == 0 {
		return nil
	}
	var result [][]client.Filter
	for _, group := range groups {
		var conditions []client.Filter
		for _, c := range group.Conditions {
			f := client.Filter{
				Key:  c.Key.ValueString(),
				Oper: c.Oper.ValueString(),
			}
			if !c.Vals.IsNull() {
				var vals []string
				diags.Append(c.Vals.ElementsAs(ctx, &vals, false)...)
				f.Vals = vals
			}
			conditions = append(conditions, f)
		}
		result = append(result, conditions)
	}
	return result
}

func (r *EscalateRuleResource) readTimeFilters(ctx context.Context, tfs []client.TimeFilter, diags *diag.Diagnostics) []EscalateTimeFilterModel {
	if len(tfs) == 0 {
		return nil
	}
	var result []EscalateTimeFilterModel
	for _, tf := range tfs {
		m := EscalateTimeFilterModel{
			Start: types.StringValue(tf.Start),
			End:   types.StringValue(tf.End),
			IsOff: types.BoolValue(tf.IsOff),
		}
		if tf.CalID != "" {
			m.CalID = types.StringValue(tf.CalID)
		}
		if len(tf.Repeat) > 0 {
			int64s := make([]int64, len(tf.Repeat))
			for i, v := range tf.Repeat {
				int64s[i] = int64(v)
			}
			vals, d := types.ListValueFrom(ctx, types.Int64Type, int64s)
			diags.Append(d...)
			m.Repeat = vals
		} else {
			m.Repeat = types.ListNull(types.Int64Type)
		}
		result = append(result, m)
	}
	return result
}

func (r *EscalateRuleResource) readFilters(ctx context.Context, groups [][]client.Filter, diags *diag.Diagnostics) []EscalateFilterGroupModel {
	if len(groups) == 0 {
		return nil
	}
	var result []EscalateFilterGroupModel
	for _, group := range groups {
		var conditions []EscalateFilterModel
		for _, f := range group {
			m := EscalateFilterModel{
				Key:  types.StringValue(f.Key),
				Oper: types.StringValue(f.Oper),
			}
			if len(f.Vals) > 0 {
				vals, d := types.ListValueFrom(ctx, types.StringType, f.Vals)
				diags.Append(d...)
				m.Vals = vals
			} else {
				m.Vals = types.ListNull(types.StringType)
			}
			conditions = append(conditions, m)
		}
		result = append(result, EscalateFilterGroupModel{Conditions: conditions})
	}
	return result
}
