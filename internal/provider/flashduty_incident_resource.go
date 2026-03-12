package provider

import (
	"context"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &IncidentResource{}
	_ resource.ResourceWithImportState = &IncidentResource{}
)

func NewIncidentResource() resource.Resource {
	return &IncidentResource{}
}

type IncidentResource struct {
	client *client.Client
}

type IncidentResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Title            types.String `tfsdk:"title"`
	Description      types.String `tfsdk:"description"`
	IncidentSeverity types.String `tfsdk:"incident_severity"`
	ChannelID        types.Int64  `tfsdk:"channel_id"`
	IncidentStatus   types.String `tfsdk:"incident_status"`
	Progress         types.String `tfsdk:"progress"`
	Impact           types.String `tfsdk:"impact"`
	RootCause        types.String `tfsdk:"root_cause"`
	Resolution       types.String `tfsdk:"resolution"`
}

func (r *IncidentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_incident"
}

func (r *IncidentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Flashduty incident.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the incident (ObjectID).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The title of the incident.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the incident. Can be plain text or markdown.",
			},
			"incident_severity": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The severity of the incident. Valid values: `Critical`, `Warning`, `Info`.",
				Validators: []validator.String{
					stringvalidator.OneOf("Critical", "Warning", "Info"),
				},
			},
			"channel_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The ID of the channel (collaboration space) for this incident.",
			},
			"incident_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the incident. Values: `Info`, `Warning`, `Critical`, `Ok`.",
			},
			"progress": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The processing progress. Values: `Triggered`, `Processing`, `Closed`.",
			},
			"impact": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The impact description of the incident.",
			},
			"root_cause": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The root cause analysis of the incident.",
			},
			"resolution": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The resolution/solution for the incident.",
			},
		},
	}
}

func (r *IncidentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (r *IncidentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IncidentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateIncidentRequest{
		Title:            data.Title.ValueString(),
		IncidentSeverity: data.IncidentSeverity.ValueString(),
	}

	if !data.Description.IsNull() {
		createReq.Description = data.Description.ValueString()
	}
	if !data.ChannelID.IsNull() {
		createReq.ChannelID = data.ChannelID.ValueInt64()
	}

	result, err := r.client.CreateIncident(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Incident", err.Error())
		return
	}

	data.ID = types.StringValue(result.IncidentID)

	incident, err := r.client.GetIncident(ctx, &client.GetIncidentRequest{
		IncidentID: result.IncidentID,
	})
	if err == nil {
		data.IncidentStatus = types.StringValue(incident.IncidentStatus)
		data.Progress = types.StringValue(incident.Progress)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncidentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IncidentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	incident, err := r.client.GetIncident(ctx, &client.GetIncidentRequest{
		IncidentID: data.ID.ValueString(),
	})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Incident", err.Error())
		return
	}

	data.Title = types.StringValue(incident.Title)
	data.IncidentSeverity = types.StringValue(incident.IncidentSeverity)
	data.IncidentStatus = types.StringValue(incident.IncidentStatus)
	data.Progress = types.StringValue(incident.Progress)

	if incident.Description != "" {
		data.Description = types.StringValue(incident.Description)
	}
	if incident.ChannelID != 0 {
		data.ChannelID = types.Int64Value(incident.ChannelID)
	}
	if incident.Impact != "" {
		data.Impact = types.StringValue(incident.Impact)
	}
	if incident.RootCause != "" {
		data.RootCause = types.StringValue(incident.RootCause)
	}
	if incident.Resolution != "" {
		data.Resolution = types.StringValue(incident.Resolution)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncidentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IncidentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &client.UpdateIncidentRequest{
		IncidentID:       data.ID.ValueString(),
		Title:            data.Title.ValueString(),
		IncidentSeverity: data.IncidentSeverity.ValueString(),
	}

	if !data.Description.IsNull() {
		updateReq.Description = data.Description.ValueString()
	}
	if !data.Impact.IsNull() {
		updateReq.Impact = data.Impact.ValueString()
	}
	if !data.RootCause.IsNull() {
		updateReq.RootCause = data.RootCause.ValueString()
	}
	if !data.Resolution.IsNull() {
		updateReq.Resolution = data.Resolution.ValueString()
	}

	err := r.client.UpdateIncident(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Incident", err.Error())
		return
	}

	incident, err := r.client.GetIncident(ctx, &client.GetIncidentRequest{
		IncidentID: data.ID.ValueString(),
	})
	if err == nil {
		data.IncidentStatus = types.StringValue(incident.IncidentStatus)
		data.Progress = types.StringValue(incident.Progress)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncidentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IncidentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIncident(ctx, &client.DeleteIncidentRequest{
		IncidentIDs: []string{data.ID.ValueString()},
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Incident", err.Error())
		return
	}
}

func (r *IncidentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
