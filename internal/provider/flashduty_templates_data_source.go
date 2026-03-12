package provider

import (
	"context"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &TemplatesDataSource{}

func NewTemplatesDataSource() datasource.DataSource {
	return &TemplatesDataSource{}
}

type TemplatesDataSource struct {
	client *client.Client
}

type TemplateItemModel struct {
	TemplateID   types.String `tfsdk:"template_id"`
	TemplateName types.String `tfsdk:"template_name"`
	Description  types.String `tfsdk:"description"`
	TeamID       types.Int64  `tfsdk:"team_id"`
	Status       types.String `tfsdk:"status"`
}

type TemplatesDataSourceModel struct {
	Query     types.String        `tfsdk:"query"`
	TeamIDs   types.List          `tfsdk:"team_ids"`
	Templates []TemplateItemModel `tfsdk:"templates"`
}

func (d *TemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_templates"
}

func (d *TemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get a list of Flashduty notification templates.",

		Attributes: map[string]schema.Attribute{
			"query": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Search query to filter templates by name.",
			},
			"team_ids": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "Filter templates by team IDs.",
			},
			"templates": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of templates.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"template_id": schema.StringAttribute{
							Computed:            true,
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
							MarkdownDescription: "The ID of the team.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status of the template.",
						},
					},
				},
			},
		},
	}
}

func (d *TemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *TemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TemplatesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := ""
	if !data.Query.IsNull() {
		query = data.Query.ValueString()
	}
	var teamIDs []int64
	if !data.TeamIDs.IsNull() {
		resp.Diagnostics.Append(data.TeamIDs.ElementsAs(ctx, &teamIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var allTemplates []client.Template
	page := 1
	limit := 100
	for {
		result, err := d.client.ListTemplates(ctx, &client.ListTemplatesRequest{
			Page:    page,
			Limit:   limit,
			Query:   query,
			TeamIDs: teamIDs,
		})
		if err != nil {
			resp.Diagnostics.AddError("Error Listing Templates", err.Error())
			return
		}
		allTemplates = append(allTemplates, result.Items...)
		if len(result.Items) < limit {
			break
		}
		page++
	}

	data.Templates = make([]TemplateItemModel, len(allTemplates))
	for i, t := range allTemplates {
		data.Templates[i] = TemplateItemModel{
			TemplateID:   types.StringValue(t.TemplateID),
			TemplateName: types.StringValue(t.TemplateName),
			Description:  types.StringValue(t.Description),
			TeamID:       types.Int64Value(t.TeamID),
			Status:       types.StringValue(t.Status),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
