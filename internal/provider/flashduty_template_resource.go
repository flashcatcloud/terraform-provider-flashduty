package provider

import (
	"context"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &TemplateResource{}
	_ resource.ResourceWithImportState = &TemplateResource{}
)

func NewTemplateResource() resource.Resource {
	return &TemplateResource{}
}

type TemplateResource struct {
	client *client.Client
}

type TemplateResourceModel struct {
	ID           types.String `tfsdk:"id"`
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
}

func (r *TemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (r *TemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Flashduty notification template.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"template_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the template.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the template.",
			},
			"team_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The ID of the team this template belongs to.",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Email notification template content.",
			},
			"sms": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "SMS notification template content.",
			},
			"dingtalk": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "DingTalk bot notification template content.",
			},
			"wecom": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "WeCom bot notification template content.",
			},
			"feishu": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Feishu bot notification template content.",
			},
			"feishu_app": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Feishu app notification template content.",
			},
			"dingtalk_app": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "DingTalk app notification template content.",
			},
			"wecom_app": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "WeCom app notification template content.",
			},
			"teams_app": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Microsoft Teams app notification template content.",
			},
			"slack_app": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Slack app notification template content.",
			},
			"slack": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Slack bot notification template content.",
			},
			"zoom": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Zoom bot notification template content.",
			},
			"telegram": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Telegram bot notification template content.",
			},
		},
	}
}

func (r *TemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (r *TemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateTemplateRequest{
		TemplateName: data.TemplateName.ValueString(),
		Description:  data.Description.ValueString(),
		Email:        data.Email.ValueString(),
		SMS:          data.SMS.ValueString(),
		Dingtalk:     data.Dingtalk.ValueString(),
		Wecom:        data.Wecom.ValueString(),
		Feishu:       data.Feishu.ValueString(),
		FeishuApp:    data.FeishuApp.ValueString(),
		DingtalkApp:  data.DingtalkApp.ValueString(),
		WecomApp:     data.WecomApp.ValueString(),
		TeamsApp:     data.TeamsApp.ValueString(),
		SlackApp:     data.SlackApp.ValueString(),
		Slack:        data.Slack.ValueString(),
		Zoom:         data.Zoom.ValueString(),
		Telegram:     data.Telegram.ValueString(),
	}

	if !data.TeamID.IsNull() {
		createReq.TeamID = data.TeamID.ValueInt64()
	}

	result, err := r.client.CreateTemplate(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Template", err.Error())
		return
	}

	data.ID = types.StringValue(result.TemplateID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template, err := r.client.GetTemplate(ctx, &client.GetTemplateRequest{
		TemplateID: data.ID.ValueString(),
	})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Template", err.Error())
		return
	}

	data.TemplateName = types.StringValue(template.TemplateName)
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

	if template.Description != "" {
		data.Description = types.StringValue(template.Description)
	}

	if template.TeamID != 0 {
		data.TeamID = types.Int64Value(template.TeamID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &client.UpdateTemplateRequest{
		TemplateID:   data.ID.ValueString(),
		TemplateName: data.TemplateName.ValueString(),
		Description:  data.Description.ValueString(),
		Email:        data.Email.ValueString(),
		SMS:          data.SMS.ValueString(),
		Dingtalk:     data.Dingtalk.ValueString(),
		Wecom:        data.Wecom.ValueString(),
		Feishu:       data.Feishu.ValueString(),
		FeishuApp:    data.FeishuApp.ValueString(),
		DingtalkApp:  data.DingtalkApp.ValueString(),
		WecomApp:     data.WecomApp.ValueString(),
		TeamsApp:     data.TeamsApp.ValueString(),
		SlackApp:     data.SlackApp.ValueString(),
		Slack:        data.Slack.ValueString(),
		Zoom:         data.Zoom.ValueString(),
		Telegram:     data.Telegram.ValueString(),
	}

	if !data.TeamID.IsNull() {
		updateReq.TeamID = data.TeamID.ValueInt64()
	}

	err := r.client.UpdateTemplate(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Template", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTemplate(ctx, &client.DeleteTemplateRequest{
		TemplateID: data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Template", err.Error())
		return
	}
}

func (r *TemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
