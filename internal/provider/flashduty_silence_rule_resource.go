package provider

import (
	"context"
	"strconv"
	"strings"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &SilenceRuleResource{}
	_ resource.ResourceWithImportState = &SilenceRuleResource{}
)

func NewSilenceRuleResource() resource.Resource {
	return &SilenceRuleResource{}
}

type SilenceRuleResource struct {
	client *client.Client
}

type SingleTimeModel struct {
	StartTime types.Int64 `tfsdk:"start_time"`
	EndTime   types.Int64 `tfsdk:"end_time"`
}

type FilterModel struct {
	Key  types.String `tfsdk:"key"`
	Oper types.String `tfsdk:"oper"`
	Vals types.List   `tfsdk:"vals"`
}

type FilterGroupModel struct {
	Conditions []FilterModel `tfsdk:"conditions"`
}

type TimeFilterModel struct {
	Start  types.String `tfsdk:"start"`
	End    types.String `tfsdk:"end"`
	Repeat types.List   `tfsdk:"repeat"`
}

type SilenceRuleResourceModel struct {
	ID                types.String       `tfsdk:"id"`
	ChannelID         types.Int64        `tfsdk:"channel_id"`
	RuleName          types.String       `tfsdk:"rule_name"`
	Description       types.String       `tfsdk:"description"`
	Priority          types.Int64        `tfsdk:"priority"`
	Filters           []FilterGroupModel `tfsdk:"filters"`
	TimeFilter        *SingleTimeModel   `tfsdk:"time_filter"`
	TimeFilters       []TimeFilterModel  `tfsdk:"time_filters"`
	IsDirectlyDiscard types.Bool         `tfsdk:"is_directly_discard"`
}

func (r *SilenceRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_silence_rule"
}

func (r *SilenceRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Flashduty silence rule for a channel.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the silence rule.",
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
				MarkdownDescription: "The name of the silence rule.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The description of the rule.",
			},
			"priority": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "The priority of the silence rule for ordering.",
			},
			"filters": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "The filter conditions (OR groups of AND conditions).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"conditions": schema.ListNestedAttribute{
							Required:            true,
							MarkdownDescription: "AND conditions within this group.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The attribute or label key (e.g., `title`, `severity`, `labels.xxx`).",
									},
									"oper": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The operator: `IN` or `NOTIN`.",
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
			"time_filter": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Single time range (mutually exclusive with `time_filters`).",
				Attributes: map[string]schema.Attribute{
					"start_time": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Start timestamp.",
					},
					"end_time": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "End timestamp.",
					},
				},
			},
			"time_filters": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Recurring time filters (mutually exclusive with `time_filter`).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Start time (e.g., `10:00`).",
						},
						"end": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "End time (e.g., `23:59`).",
						},
						"repeat": schema.ListAttribute{
							Required:            true,
							ElementType:         types.Int64Type,
							MarkdownDescription: "Days of week (0=Sunday, 1-6=Monday-Saturday).",
						},
					},
				},
			},
			"is_directly_discard": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to directly discard matching alerts.",
			},
		},
	}
}

func (r *SilenceRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (r *SilenceRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SilenceRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filters := r.buildFilters(ctx, data.Filters, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateSilenceRuleRequest{
		ChannelID:         data.ChannelID.ValueInt64(),
		RuleName:          data.RuleName.ValueString(),
		Description:       data.Description.ValueString(),
		Priority:          int(data.Priority.ValueInt64()),
		Filters:           filters,
		IsDirectlyDiscard: data.IsDirectlyDiscard.ValueBool(),
	}

	if data.TimeFilter != nil {
		createReq.TimeFilter = &client.SingleTime{
			StartTime: data.TimeFilter.StartTime.ValueInt64(),
			EndTime:   data.TimeFilter.EndTime.ValueInt64(),
		}
	}

	if len(data.TimeFilters) > 0 {
		createReq.TimeFilters = r.buildTimeFilters(ctx, data.TimeFilters, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	result, err := r.client.CreateSilenceRule(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Silence Rule", err.Error())
		return
	}

	data.ID = types.StringValue(result.RuleID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SilenceRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SilenceRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetSilenceRule(ctx, &client.GetSilenceRuleRequest{
		ChannelID: data.ChannelID.ValueInt64(),
		RuleID:    data.ID.ValueString(),
	})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Silence Rule", err.Error())
		return
	}

	data.RuleName = types.StringValue(rule.RuleName)
	data.Description = types.StringValue(rule.Description)
	data.Priority = types.Int64Value(int64(rule.Priority))
	data.IsDirectlyDiscard = types.BoolValue(rule.IsDirectlyDiscard)

	data.Filters = r.readFilters(ctx, rule.Filters, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if rule.TimeFilter != nil {
		data.TimeFilter = &SingleTimeModel{
			StartTime: types.Int64Value(rule.TimeFilter.StartTime),
			EndTime:   types.Int64Value(rule.TimeFilter.EndTime),
		}
	}

	data.TimeFilters = r.readTimeFilters(ctx, rule.TimeFilters, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SilenceRuleResource) readFilters(ctx context.Context, filters [][]client.Filter, diags *diag.Diagnostics) []FilterGroupModel {
	var result []FilterGroupModel
	for _, group := range filters {
		var conditions []FilterModel
		for _, f := range group {
			vals, d := types.ListValueFrom(ctx, types.StringType, f.Vals)
			diags.Append(d...)
			conditions = append(conditions, FilterModel{
				Key:  types.StringValue(f.Key),
				Oper: types.StringValue(f.Oper),
				Vals: vals,
			})
		}
		result = append(result, FilterGroupModel{Conditions: conditions})
	}
	return result
}

func (r *SilenceRuleResource) readTimeFilters(ctx context.Context, timeFilters []client.TimeFilter, diags *diag.Diagnostics) []TimeFilterModel {
	var result []TimeFilterModel
	for _, tf := range timeFilters {
		repeat, d := types.ListValueFrom(ctx, types.Int64Type, tf.Repeat)
		diags.Append(d...)
		result = append(result, TimeFilterModel{
			Start:  types.StringValue(tf.Start),
			End:    types.StringValue(tf.End),
			Repeat: repeat,
		})
	}
	return result
}

func (r *SilenceRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SilenceRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filters := r.buildFilters(ctx, data.Filters, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	isDirectlyDiscard := data.IsDirectlyDiscard.ValueBool()
	priority := int(data.Priority.ValueInt64())
	updateReq := &client.UpdateSilenceRuleRequest{
		ChannelID:         data.ChannelID.ValueInt64(),
		RuleID:            data.ID.ValueString(),
		RuleName:          data.RuleName.ValueString(),
		Description:       data.Description.ValueString(),
		Priority:          &priority,
		Filters:           filters,
		IsDirectlyDiscard: &isDirectlyDiscard,
	}

	if data.TimeFilter != nil {
		updateReq.TimeFilter = &client.SingleTime{
			StartTime: data.TimeFilter.StartTime.ValueInt64(),
			EndTime:   data.TimeFilter.EndTime.ValueInt64(),
		}
	}

	if len(data.TimeFilters) > 0 {
		updateReq.TimeFilters = r.buildTimeFilters(ctx, data.TimeFilters, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	err := r.client.UpdateSilenceRule(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Silence Rule", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SilenceRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SilenceRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Disable the rule first (required before deletion)
	err := r.client.DisableSilenceRule(ctx, &client.EnableSilenceRuleRequest{
		ChannelID: data.ChannelID.ValueInt64(),
		RuleID:    data.ID.ValueString(),
	})
	if err != nil && !client.IsNotFoundError(err) {
		resp.Diagnostics.AddError("Error Disabling Silence Rule", err.Error())
		return
	}

	// Now delete the rule
	err = r.client.DeleteSilenceRule(ctx, &client.DeleteSilenceRuleRequest{
		ChannelID: data.ChannelID.ValueInt64(),
		RuleID:    data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Silence Rule", err.Error())
		return
	}
}

func (r *SilenceRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *SilenceRuleResource) buildFilters(ctx context.Context, filters []FilterGroupModel, diags *diag.Diagnostics) [][]client.Filter {
	var result [][]client.Filter

	for _, group := range filters {
		var andGroup []client.Filter
		for _, f := range group.Conditions {
			var values []string
			if !f.Vals.IsNull() && !f.Vals.IsUnknown() {
				diags.Append(f.Vals.ElementsAs(ctx, &values, false)...)
			}
			filter := client.Filter{
				Key:  f.Key.ValueString(),
				Oper: f.Oper.ValueString(),
				Vals: values,
			}
			andGroup = append(andGroup, filter)
		}
		result = append(result, andGroup)
	}

	return result
}

func (r *SilenceRuleResource) buildTimeFilters(ctx context.Context, timeFilters []TimeFilterModel, diags *diag.Diagnostics) []client.TimeFilter {
	var result []client.TimeFilter

	for _, tf := range timeFilters {
		filter := client.TimeFilter{
			Start: tf.Start.ValueString(),
			End:   tf.End.ValueString(),
		}
		var repeat []int
		diags.Append(tf.Repeat.ElementsAs(ctx, &repeat, false)...)
		filter.Repeat = repeat
		result = append(result, filter)
	}

	return result
}
