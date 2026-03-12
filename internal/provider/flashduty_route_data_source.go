package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &RouteDataSource{}
	_ datasource.DataSourceWithConfigure = &RouteDataSource{}
)

func NewRouteDataSource() datasource.DataSource {
	return &RouteDataSource{}
}

type RouteDataSource struct {
	client *client.Client
}

type RouteDataSourceModel struct {
	ID            types.String               `tfsdk:"id"`
	IntegrationID types.Int64                `tfsdk:"integration_id"`
	Cases         []RouteCaseDataSourceModel `tfsdk:"cases"`
	Sections      []RouteSectionModel        `tfsdk:"sections"`
	Default       *RouteDefaultModel         `tfsdk:"default"`
	Version       types.Int64                `tfsdk:"version"`
	Status        types.String               `tfsdk:"status"`
	CreatedAt     types.Int64                `tfsdk:"created_at"`
	UpdatedAt     types.Int64                `tfsdk:"updated_at"`
}

type RouteCaseDataSourceModel struct {
	If               []RouteFilterDataSourceModel `tfsdk:"if"`
	ChannelIDs       types.List                   `tfsdk:"channel_ids"`
	Fallthrough      types.Bool                   `tfsdk:"fallthrough"`
	RoutingMode      types.String                 `tfsdk:"routing_mode"`
	NameMappingLabel types.String                 `tfsdk:"name_mapping_label"`
}

type RouteFilterDataSourceModel struct {
	Key  types.String `tfsdk:"key"`
	Oper types.String `tfsdk:"oper"`
	Vals types.List   `tfsdk:"vals"`
}

func (d *RouteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route"
}

func (d *RouteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves routing configuration for a shared alert integration from Flashduty.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the route.",
				Computed:            true,
			},
			"integration_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the shared alert integration.",
				Required:            true,
			},
			"cases": schema.ListNestedAttribute{
				MarkdownDescription: "Conditional routing rules.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"if": schema.ListNestedAttribute{
							MarkdownDescription: "Filter conditions.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										MarkdownDescription: "The attribute or label key.",
										Computed:            true,
									},
									"oper": schema.StringAttribute{
										MarkdownDescription: "The operator.",
										Computed:            true,
									},
									"vals": schema.ListAttribute{
										MarkdownDescription: "The values to match.",
										Computed:            true,
										ElementType:         types.StringType,
									},
								},
							},
						},
						"channel_ids": schema.ListAttribute{
							MarkdownDescription: "Target channel IDs.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"fallthrough": schema.BoolAttribute{
							MarkdownDescription: "Whether to continue matching after this rule.",
							Computed:            true,
						},
						"routing_mode": schema.StringAttribute{
							MarkdownDescription: "Routing mode.",
							Computed:            true,
						},
						"name_mapping_label": schema.StringAttribute{
							MarkdownDescription: "Label key for name mapping mode.",
							Computed:            true,
						},
					},
				},
			},
			"sections": schema.ListNestedAttribute{
				MarkdownDescription: "Section dividers in route rules.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Section name.",
							Computed:            true,
						},
						"position": schema.Int64Attribute{
							MarkdownDescription: "Position.",
							Computed:            true,
						},
					},
				},
			},
			"default": schema.SingleNestedAttribute{
				MarkdownDescription: "Default routing configuration.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"channel_ids": schema.ListAttribute{
						MarkdownDescription: "Default target channel IDs.",
						Computed:            true,
						ElementType:         types.Int64Type,
					},
				},
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "The version number of the route configuration.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the route.",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The timestamp when the route was created.",
				Computed:            true,
			},
			"updated_at": schema.Int64Attribute{
				MarkdownDescription: "The timestamp when the route was last updated.",
				Computed:            true,
			},
		},
	}
}

func (d *RouteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *RouteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RouteDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	route, err := d.client.GetRoute(ctx, &client.GetRouteRequest{IntegrationID: state.IntegrationID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Route", fmt.Sprintf("Could not read route: %s", err))
		return
	}

	d.mapRouteToModel(ctx, route, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = types.StringValue(strconv.FormatInt(state.IntegrationID.ValueInt64(), 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *RouteDataSource) mapRouteToModel(ctx context.Context, route *client.Route, model *RouteDataSourceModel, diags *diag.Diagnostics) {
	model.Version = types.Int64Value(int64(route.Version))
	model.Status = types.StringValue(route.Status)
	model.CreatedAt = types.Int64Value(route.CreatedAt)
	model.UpdatedAt = types.Int64Value(route.UpdatedAt)
	model.Cases = mapRouteCasesToDataSourceModels(ctx, route.Cases, diags)
	model.Sections = mapRouteSectionsToModels(route.Sections)
	model.Default = mapRouteDefaultToModel(ctx, route.Default, diags)
}

func mapRouteCasesToDataSourceModels(ctx context.Context, cases []client.RouteCase, diags *diag.Diagnostics) []RouteCaseDataSourceModel {
	if len(cases) == 0 {
		return nil
	}
	result := make([]RouteCaseDataSourceModel, len(cases))
	for i, c := range cases {
		result[i] = mapRouteCaseToDataSourceModel(ctx, &c, diags)
	}
	return result
}

func mapRouteCaseToDataSourceModel(ctx context.Context, c *client.RouteCase, diags *diag.Diagnostics) RouteCaseDataSourceModel {
	m := RouteCaseDataSourceModel{
		Fallthrough:      types.BoolValue(c.Fallthrough),
		RoutingMode:      types.StringValue(c.RoutingMode),
		NameMappingLabel: types.StringValue(c.NameMappingLabel),
	}
	m.If = make([]RouteFilterDataSourceModel, len(c.If))
	for j, f := range c.If {
		vals, d := types.ListValueFrom(ctx, types.StringType, f.Vals)
		diags.Append(d...)
		m.If[j] = RouteFilterDataSourceModel{
			Key:  types.StringValue(f.Key),
			Oper: types.StringValue(f.Oper),
			Vals: vals,
		}
	}
	if len(c.ChannelIDs) > 0 {
		channelIDs, d := types.ListValueFrom(ctx, types.Int64Type, c.ChannelIDs)
		diags.Append(d...)
		m.ChannelIDs = channelIDs
	} else {
		m.ChannelIDs = types.ListNull(types.Int64Type)
	}
	return m
}

func mapRouteSectionsToModels(sections []client.RouteSection) []RouteSectionModel {
	if len(sections) == 0 {
		return nil
	}
	result := make([]RouteSectionModel, len(sections))
	for i, s := range sections {
		result[i] = RouteSectionModel{
			Name:     types.StringValue(s.Name),
			Position: types.Int64Value(int64(s.Position)),
		}
	}
	return result
}

func mapRouteDefaultToModel(ctx context.Context, def *client.RouteDefault, diags *diag.Diagnostics) *RouteDefaultModel {
	if def == nil || len(def.ChannelIDs) == 0 {
		return nil
	}
	channelIDs, d := types.ListValueFrom(ctx, types.Int64Type, def.ChannelIDs)
	diags.Append(d...)
	return &RouteDefaultModel{ChannelIDs: channelIDs}
}

// RouteHistoryDataSource - list route history.
var (
	_ datasource.DataSource              = &RouteHistoryDataSource{}
	_ datasource.DataSourceWithConfigure = &RouteHistoryDataSource{}
)

func NewRouteHistoryDataSource() datasource.DataSource {
	return &RouteHistoryDataSource{}
}

type RouteHistoryDataSource struct {
	client *client.Client
}

type RouteHistoryDataSourceModel struct {
	IntegrationID types.Int64             `tfsdk:"integration_id"`
	Items         []RouteHistoryItemModel `tfsdk:"items"`
}

type RouteHistoryItemModel struct {
	Cases     []RouteCaseDataSourceModel `tfsdk:"cases"`
	Sections  []RouteSectionModel        `tfsdk:"sections"`
	Default   *RouteDefaultModel         `tfsdk:"default"`
	Version   types.Int64                `tfsdk:"version"`
	UpdatedBy types.Int64                `tfsdk:"updated_by"`
	UpdatedAt types.Int64                `tfsdk:"updated_at"`
}

func (d *RouteHistoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route_history"
}

func (d *RouteHistoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves route history for a shared alert integration from Flashduty.",
		Attributes: map[string]schema.Attribute{
			"integration_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the shared alert integration.",
				Required:            true,
			},
			"items": schema.ListNestedAttribute{
				MarkdownDescription: "List of route history items.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cases": schema.ListNestedAttribute{
							MarkdownDescription: "Conditional routing rules.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"if": schema.ListNestedAttribute{
										MarkdownDescription: "Filter conditions.",
										Computed:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													MarkdownDescription: "The attribute or label key.",
													Computed:            true,
												},
												"oper": schema.StringAttribute{
													MarkdownDescription: "The operator.",
													Computed:            true,
												},
												"vals": schema.ListAttribute{
													MarkdownDescription: "The values to match.",
													Computed:            true,
													ElementType:         types.StringType,
												},
											},
										},
									},
									"channel_ids": schema.ListAttribute{
										MarkdownDescription: "Target channel IDs.",
										Computed:            true,
										ElementType:         types.Int64Type,
									},
									"fallthrough": schema.BoolAttribute{
										MarkdownDescription: "Whether to continue matching.",
										Computed:            true,
									},
									"routing_mode": schema.StringAttribute{
										MarkdownDescription: "Routing mode.",
										Computed:            true,
									},
									"name_mapping_label": schema.StringAttribute{
										MarkdownDescription: "Label key for name mapping mode.",
										Computed:            true,
									},
								},
							},
						},
						"sections": schema.ListNestedAttribute{
							MarkdownDescription: "Section dividers.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										MarkdownDescription: "Section name.",
										Computed:            true,
									},
									"position": schema.Int64Attribute{
										MarkdownDescription: "Position.",
										Computed:            true,
									},
								},
							},
						},
						"default": schema.SingleNestedAttribute{
							MarkdownDescription: "Default routing configuration.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"channel_ids": schema.ListAttribute{
									MarkdownDescription: "Default target channel IDs.",
									Computed:            true,
									ElementType:         types.Int64Type,
								},
							},
						},
						"version": schema.Int64Attribute{
							MarkdownDescription: "The version number.",
							Computed:            true,
						},
						"updated_by": schema.Int64Attribute{
							MarkdownDescription: "The user ID who updated.",
							Computed:            true,
						},
						"updated_at": schema.Int64Attribute{
							MarkdownDescription: "The timestamp when updated.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *RouteHistoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *RouteHistoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RouteHistoryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListRouteHistory(ctx, &client.ListRouteHistoryRequest{
		IntegrationID: state.IntegrationID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Route History", fmt.Sprintf("Could not read route history: %s", err))
		return
	}

	for _, history := range result.Items {
		state.Items = append(state.Items, RouteHistoryItemModel{
			Cases:     mapRouteCasesToDataSourceModels(ctx, history.Cases, &resp.Diagnostics),
			Sections:  mapRouteSectionsToModels(history.Sections),
			Default:   mapRouteDefaultToModel(ctx, history.Default, &resp.Diagnostics),
			Version:   types.Int64Value(int64(history.Version)),
			UpdatedBy: types.Int64Value(history.UpdatedBy),
			UpdatedAt: types.Int64Value(history.UpdatedAt),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
