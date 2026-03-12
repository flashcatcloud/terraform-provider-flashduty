package provider

import (
	"context"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &TemplateDataSource{}

func NewTemplateDataSource() datasource.DataSource {
	return &TemplateDataSource{}
}

type TemplateDataSource struct {
	client *client.Client
}

type TemplateDataSourceModel struct {
	TemplateID   types.String `tfsdk:"template_id"`
	TemplateName types.String `tfsdk:"template_name"`
	Description  types.String `tfsdk:"description"`
	TeamID       types.Int64  `tfsdk:"team_id"`
	Email        types.String `tfsdk:"email"`
	SMS          types.String `tfsdk:"sms"`
	Dingtalk     types.String `tfsdk:"dingtalk"`
	Wecom        types.String `tfsdk:"wecom"`
	Feishu       types.String `tfsdk:"feishu"`
	FeishuApp    types.String `tfsdk:"feishu_app"`
	DingtalkApp  types.String `tfsdk:"dingtalk_app"`
	WecomApp     types.String `tfsdk:"wecom_app"`
	TeamsApp     types.String `tfsdk:"teams_app"`
	SlackApp     types.String `tfsdk:"slack_app"`
	Slack        types.String `tfsdk:"slack"`
	Zoom         types.String `tfsdk:"zoom"`
	Telegram     types.String `tfsdk:"telegram"`
	Status       types.String `tfsdk:"status"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
	UpdatedAt    types.Int64  `tfsdk:"updated_at"`
}

func (d *TemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (d *TemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a Flashduty notification template.",

		Attributes: map[string]schema.Attribute{
			"template_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the template.",
			},
			"template_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the template.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the template.",
			},
			"team_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The ID of the team this template belongs to.",
			},
			"email": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Email notification template content.",
			},
			"sms": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SMS notification template content.",
			},
			"dingtalk": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DingTalk bot notification template content.",
			},
			"wecom": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "WeCom bot notification template content.",
			},
			"feishu": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Feishu bot notification template content.",
			},
			"feishu_app": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Feishu app notification template content.",
			},
			"dingtalk_app": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DingTalk app notification template content.",
			},
			"wecom_app": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "WeCom app notification template content.",
			},
			"teams_app": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Microsoft Teams app notification template content.",
			},
			"slack_app": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Slack app notification template content.",
			},
			"slack": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Slack bot notification template content.",
			},
			"zoom": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Zoom bot notification template content.",
			},
			"telegram": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Telegram bot notification template content.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the template: `enabled`, `disabled`, or `deleted`.",
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

func (d *TemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *TemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TemplateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template, err := d.client.GetTemplate(ctx, &client.GetTemplateRequest{
		TemplateID: data.TemplateID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Template", err.Error())
		return
	}

	data.TemplateName = types.StringValue(template.TemplateName)
	data.Description = types.StringValue(template.Description)
	data.TeamID = types.Int64Value(template.TeamID)
	data.Email = types.StringValue(template.Email)
	data.SMS = types.StringValue(template.SMS)
	data.Dingtalk = types.StringValue(template.Dingtalk)
	data.Wecom = types.StringValue(template.Wecom)
	data.Feishu = types.StringValue(template.Feishu)
	data.FeishuApp = types.StringValue(template.FeishuApp)
	data.DingtalkApp = types.StringValue(template.DingtalkApp)
	data.WecomApp = types.StringValue(template.WecomApp)
	data.TeamsApp = types.StringValue(template.TeamsApp)
	data.SlackApp = types.StringValue(template.SlackApp)
	data.Slack = types.StringValue(template.Slack)
	data.Zoom = types.StringValue(template.Zoom)
	data.Telegram = types.StringValue(template.Telegram)
	data.Status = types.StringValue(template.Status)
	data.CreatedAt = types.Int64Value(template.CreatedAt)
	data.UpdatedAt = types.Int64Value(template.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
