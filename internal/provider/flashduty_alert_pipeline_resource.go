package provider

import (
	"context"
	"encoding/json"
	"strconv"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &AlertPipelineResource{}
	_ resource.ResourceWithImportState = &AlertPipelineResource{}
)

func NewAlertPipelineResource() resource.Resource {
	return &AlertPipelineResource{}
}

type AlertPipelineResource struct {
	client *client.Client
}

type AlertPipelineFilterModel struct {
	Key  types.String `tfsdk:"key"`
	Oper types.String `tfsdk:"oper"`
	Vals types.List   `tfsdk:"vals"`
}

type AlertPipelineRuleModel struct {
	Kind     types.String               `tfsdk:"kind"`
	If       []AlertPipelineFilterModel `tfsdk:"if"`
	Settings types.String               `tfsdk:"settings"`
}

type AlertPipelineResourceModel struct {
	ID            types.String             `tfsdk:"id"`
	IntegrationID types.Int64              `tfsdk:"integration_id"`
	Rules         []AlertPipelineRuleModel `tfsdk:"rules"`
}

func (r *AlertPipelineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_pipeline"
}

func (r *AlertPipelineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages alert processing pipeline rules for a Flashduty integration. " +
			"Pipelines allow you to transform, drop, or inhibit alerts before they create incidents.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The resource identifier (same as integration_id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"integration_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the integration this pipeline belongs to.",
			},
			"rules": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "Ordered list of pipeline rules. Rules are evaluated top-down.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"kind": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The rule kind. Valid values: `title_reset`, `description_reset`, " +
								"`severity_reset`, `alert_drop`, `alert_inhibit`.",
						},
						"if": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Filter conditions that must all match for the rule to apply (AND logic).",
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
						"settings": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "JSON-encoded settings for the rule. Use `jsonencode()` in HCL. " +
								"Contents vary by kind:\n" +
								"  - `title_reset`: `{\"title\": \"...\"}`\n" +
								"  - `description_reset`: `{\"description\": \"...\"}`\n" +
								"  - `severity_reset`: `{\"severity\": \"Critical\"|\"Warning\"|\"Info\"}`\n" +
								"  - `alert_drop`: `{}`\n" +
								"  - `alert_inhibit`: `{\"source_filters\": [...], \"equals\": [...]}`",
						},
					},
				},
			},
		},
	}
}

func (r *AlertPipelineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (r *AlertPipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AlertPipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rules, err := r.buildRules(ctx, data.Rules)
	if err != nil {
		resp.Diagnostics.AddError("Error Building Pipeline Rules", err.Error())
		return
	}

	upsertReq := &client.UpsertAlertPipelineRequest{
		IntegrationID: data.IntegrationID.ValueInt64(),
		Rules:         rules,
	}

	if err := r.client.UpsertAlertPipeline(ctx, upsertReq); err != nil {
		resp.Diagnostics.AddError("Error Creating Alert Pipeline", err.Error())
		return
	}

	data.ID = types.StringValue(strconv.FormatInt(data.IntegrationID.ValueInt64(), 10))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertPipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AlertPipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrationID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error Parsing ID", err.Error())
		return
	}

	pipeline, err := r.client.GetAlertPipeline(ctx, &client.GetAlertPipelineRequest{
		IntegrationID: integrationID,
	})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Alert Pipeline", err.Error())
		return
	}

	if pipeline == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.IntegrationID = types.Int64Value(pipeline.IntegrationID)
	data.Rules = r.readRules(ctx, pipeline.Rules, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertPipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AlertPipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rules, err := r.buildRules(ctx, data.Rules)
	if err != nil {
		resp.Diagnostics.AddError("Error Building Pipeline Rules", err.Error())
		return
	}

	upsertReq := &client.UpsertAlertPipelineRequest{
		IntegrationID: data.IntegrationID.ValueInt64(),
		Rules:         rules,
	}

	if err := r.client.UpsertAlertPipeline(ctx, upsertReq); err != nil {
		resp.Diagnostics.AddError("Error Updating Alert Pipeline", err.Error())
		return
	}

	data.ID = types.StringValue(strconv.FormatInt(data.IntegrationID.ValueInt64(), 10))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete clears the pipeline by upserting an empty rules list (no dedicated delete endpoint).
func (r *AlertPipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AlertPipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upsertReq := &client.UpsertAlertPipelineRequest{
		IntegrationID: data.IntegrationID.ValueInt64(),
		Rules:         []client.AlertPipelineRule{},
	}

	if err := r.client.UpsertAlertPipeline(ctx, upsertReq); err != nil {
		resp.Diagnostics.AddError("Error Deleting Alert Pipeline", err.Error())
		return
	}
}

func (r *AlertPipelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AlertPipelineResource) buildRules(ctx context.Context, ruleModels []AlertPipelineRuleModel) ([]client.AlertPipelineRule, error) {
	rules := make([]client.AlertPipelineRule, 0, len(ruleModels))

	for _, rm := range ruleModels {
		var settings map[string]interface{}
		if err := json.Unmarshal([]byte(rm.Settings.ValueString()), &settings); err != nil {
			return nil, err
		}

		var filters []client.Filter
		for _, f := range rm.If {
			var vals []string
			if !f.Vals.IsNull() && !f.Vals.IsUnknown() {
				f.Vals.ElementsAs(ctx, &vals, false)
			}
			filters = append(filters, client.Filter{
				Key:  f.Key.ValueString(),
				Oper: f.Oper.ValueString(),
				Vals: vals,
			})
		}

		rules = append(rules, client.AlertPipelineRule{
			Kind:     rm.Kind.ValueString(),
			If:       filters,
			Settings: settings,
		})
	}

	return rules, nil
}

func (r *AlertPipelineResource) readRules(ctx context.Context, rules []client.AlertPipelineRule, diags *diag.Diagnostics) []AlertPipelineRuleModel {
	result := make([]AlertPipelineRuleModel, 0, len(rules))

	for _, rule := range rules {
		settingsJSON, err := json.Marshal(rule.Settings)
		if err != nil {
			diags.AddError("Error Serializing Settings", err.Error())
			return nil
		}

		var filters []AlertPipelineFilterModel
		for _, f := range rule.If {
			vals, d := types.ListValueFrom(ctx, types.StringType, f.Vals)
			diags.Append(d...)
			filters = append(filters, AlertPipelineFilterModel{
				Key:  types.StringValue(f.Key),
				Oper: types.StringValue(f.Oper),
				Vals: vals,
			})
		}

		result = append(result, AlertPipelineRuleModel{
			Kind:     types.StringValue(rule.Kind),
			If:       filters,
			Settings: types.StringValue(string(settingsJSON)),
		})
	}

	return result
}
