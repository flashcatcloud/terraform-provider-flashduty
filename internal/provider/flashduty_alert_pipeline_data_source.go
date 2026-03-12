package provider

import (
	"context"
	"encoding/json"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &AlertPipelineDataSource{}

func NewAlertPipelineDataSource() datasource.DataSource {
	return &AlertPipelineDataSource{}
}

type AlertPipelineDataSource struct {
	client *client.Client
}

type AlertPipelineFilterDataSourceModel struct {
	Key  types.String `tfsdk:"key"`
	Oper types.String `tfsdk:"oper"`
	Vals types.List   `tfsdk:"vals"`
}

type AlertPipelineRuleDataSourceModel struct {
	Kind     types.String                         `tfsdk:"kind"`
	If       []AlertPipelineFilterDataSourceModel `tfsdk:"if"`
	Settings types.String                         `tfsdk:"settings"`
}

type AlertPipelineDataSourceModel struct {
	IntegrationID types.Int64                        `tfsdk:"integration_id"`
	Rules         []AlertPipelineRuleDataSourceModel `tfsdk:"rules"`
}

func (d *AlertPipelineDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_pipeline"
}

func (d *AlertPipelineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about alert processing pipeline rules for a Flashduty integration.",

		Attributes: map[string]schema.Attribute{
			"integration_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the integration.",
			},
			"rules": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Ordered list of pipeline rules.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"kind": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The rule kind: `title_reset`, `description_reset`, `severity_reset`, `alert_drop`, or `alert_inhibit`.",
						},
						"if": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Filter conditions (AND logic).",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The attribute or label key.",
									},
									"oper": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The operator: `IN` or `NOTIN`.",
									},
									"vals": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "The values to match.",
									},
								},
							},
						},
						"settings": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "JSON-encoded settings for the rule.",
						},
					},
				},
			},
		},
	}
}

func (d *AlertPipelineDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *AlertPipelineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AlertPipelineDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipeline, err := d.client.GetAlertPipeline(ctx, &client.GetAlertPipelineRequest{
		IntegrationID: data.IntegrationID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Alert Pipeline", err.Error())
		return
	}

	if pipeline == nil {
		data.Rules = []AlertPipelineRuleDataSourceModel{}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	data.Rules = make([]AlertPipelineRuleDataSourceModel, len(pipeline.Rules))
	for i, rule := range pipeline.Rules {
		settingsJSON, _ := json.Marshal(rule.Settings)

		var filters []AlertPipelineFilterDataSourceModel
		for _, f := range rule.If {
			vals, d := types.ListValueFrom(ctx, types.StringType, f.Vals)
			resp.Diagnostics.Append(d...)
			filters = append(filters, AlertPipelineFilterDataSourceModel{
				Key:  types.StringValue(f.Key),
				Oper: types.StringValue(f.Oper),
				Vals: vals,
			})
		}

		data.Rules[i] = AlertPipelineRuleDataSourceModel{
			Kind:     types.StringValue(rule.Kind),
			If:       filters,
			Settings: types.StringValue(string(settingsJSON)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
