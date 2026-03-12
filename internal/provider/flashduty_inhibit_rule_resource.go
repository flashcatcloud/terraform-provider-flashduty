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
	_ resource.Resource                = &InhibitRuleResource{}
	_ resource.ResourceWithImportState = &InhibitRuleResource{}
)

func NewInhibitRuleResource() resource.Resource {
	return &InhibitRuleResource{}
}

type InhibitRuleResource struct {
	client *client.Client
}

type InhibitFilterGroupModel struct {
	Conditions []FilterModel `tfsdk:"conditions"`
}

type InhibitRuleResourceModel struct {
	ID                types.String              `tfsdk:"id"`
	ChannelID         types.Int64               `tfsdk:"channel_id"`
	RuleName          types.String              `tfsdk:"rule_name"`
	Description       types.String              `tfsdk:"description"`
	Priority          types.Int64               `tfsdk:"priority"`
	SourceFilters     []InhibitFilterGroupModel `tfsdk:"source_filters"`
	TargetFilters     []InhibitFilterGroupModel `tfsdk:"target_filters"`
	Equals            types.List                `tfsdk:"equals"`
	IsDirectlyDiscard types.Bool                `tfsdk:"is_directly_discard"`
}

func (r *InhibitRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inhibit_rule"
}

func (r *InhibitRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Flashduty inhibit rule for a channel.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the inhibit rule.",
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
				MarkdownDescription: "The name of the inhibit rule.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The description of the rule.",
			},
			"priority": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "The priority of the rule.",
			},
			"source_filters": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "Source event filter conditions (OR groups of AND conditions).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"conditions": schema.ListNestedAttribute{
							Required:            true,
							MarkdownDescription: "AND conditions within this group.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The attribute or label key.",
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
			"target_filters": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "Target event filter conditions (events to be inhibited).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"conditions": schema.ListNestedAttribute{
							Required:            true,
							MarkdownDescription: "AND conditions within this group.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The attribute or label key.",
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
			"equals": schema.ListAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Fields that must be equal between source and target (e.g., `labels.host`).",
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

func (r *InhibitRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (r *InhibitRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InhibitRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceFilters := r.buildFilters(ctx, data.SourceFilters, &resp.Diagnostics)
	targetFilters := r.buildFilters(ctx, data.TargetFilters, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var equals []string
	resp.Diagnostics.Append(data.Equals.ElementsAs(ctx, &equals, false)...)

	createReq := &client.CreateInhibitRuleRequest{
		ChannelID:         data.ChannelID.ValueInt64(),
		RuleName:          data.RuleName.ValueString(),
		Description:       data.Description.ValueString(),
		Priority:          int(data.Priority.ValueInt64()),
		SourceFilters:     sourceFilters,
		TargetFilters:     targetFilters,
		Equals:            equals,
		IsDirectlyDiscard: data.IsDirectlyDiscard.ValueBool(),
	}

	result, err := r.client.CreateInhibitRule(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Inhibit Rule", err.Error())
		return
	}

	data.ID = types.StringValue(result.RuleID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InhibitRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InhibitRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetInhibitRule(ctx, &client.GetInhibitRuleRequest{
		ChannelID: data.ChannelID.ValueInt64(),
		RuleID:    data.ID.ValueString(),
	})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Inhibit Rule", err.Error())
		return
	}

	data.RuleName = types.StringValue(rule.RuleName)
	data.Description = types.StringValue(rule.Description)
	data.Priority = types.Int64Value(int64(rule.Priority))
	data.IsDirectlyDiscard = types.BoolValue(rule.IsDirectlyDiscard)

	data.SourceFilters = r.readFilters(ctx, rule.SourceFilters)
	data.TargetFilters = r.readFilters(ctx, rule.TargetFilters)

	if len(rule.Equals) > 0 {
		equals, diags := types.ListValueFrom(ctx, types.StringType, rule.Equals)
		resp.Diagnostics.Append(diags...)
		data.Equals = equals
	} else {
		data.Equals = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InhibitRuleResource) readFilters(ctx context.Context, filters [][]client.Filter) []InhibitFilterGroupModel {
	var result []InhibitFilterGroupModel
	for _, group := range filters {
		var conditions []FilterModel
		for _, f := range group {
			vals, _ := types.ListValueFrom(ctx, types.StringType, f.Vals)
			conditions = append(conditions, FilterModel{
				Key:  types.StringValue(f.Key),
				Oper: types.StringValue(f.Oper),
				Vals: vals,
			})
		}
		result = append(result, InhibitFilterGroupModel{Conditions: conditions})
	}
	return result
}

func (r *InhibitRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InhibitRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceFilters := r.buildFilters(ctx, data.SourceFilters, &resp.Diagnostics)
	targetFilters := r.buildFilters(ctx, data.TargetFilters, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var equals []string
	resp.Diagnostics.Append(data.Equals.ElementsAs(ctx, &equals, false)...)

	priority := int(data.Priority.ValueInt64())
	isDirectlyDiscard := data.IsDirectlyDiscard.ValueBool()

	updateReq := &client.UpdateInhibitRuleRequest{
		ChannelID:         data.ChannelID.ValueInt64(),
		RuleID:            data.ID.ValueString(),
		RuleName:          data.RuleName.ValueString(),
		Description:       data.Description.ValueString(),
		Priority:          &priority,
		SourceFilters:     sourceFilters,
		TargetFilters:     targetFilters,
		Equals:            equals,
		IsDirectlyDiscard: &isDirectlyDiscard,
	}

	err := r.client.UpdateInhibitRule(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Inhibit Rule", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InhibitRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InhibitRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Disable the rule first (required before deletion)
	err := r.client.DisableInhibitRule(ctx, &client.EnableInhibitRuleRequest{
		ChannelID: data.ChannelID.ValueInt64(),
		RuleID:    data.ID.ValueString(),
	})
	if err != nil && !client.IsNotFoundError(err) {
		resp.Diagnostics.AddError("Error Disabling Inhibit Rule", err.Error())
		return
	}

	// Now delete the rule
	err = r.client.DeleteInhibitRule(ctx, &client.DeleteInhibitRuleRequest{
		ChannelID: data.ChannelID.ValueInt64(),
		RuleID:    data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Inhibit Rule", err.Error())
		return
	}
}

func (r *InhibitRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *InhibitRuleResource) buildFilters(ctx context.Context, filters []InhibitFilterGroupModel, diags *diag.Diagnostics) [][]client.Filter {
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
