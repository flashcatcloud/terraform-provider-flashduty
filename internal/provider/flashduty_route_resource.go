package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &RouteResource{}
	_ resource.ResourceWithConfigure = &RouteResource{}
)

func NewRouteResource() resource.Resource {
	return &RouteResource{}
}

type RouteResource struct {
	client *client.Client
}

type RouteResourceModel struct {
	ID            types.String        `tfsdk:"id"`
	IntegrationID types.Int64         `tfsdk:"integration_id"`
	Cases         []RouteCaseModel    `tfsdk:"cases"`
	Sections      []RouteSectionModel `tfsdk:"sections"`
	Default       *RouteDefaultModel  `tfsdk:"default"`
	Version       types.Int64         `tfsdk:"version"`
	Status        types.String        `tfsdk:"status"`
	CreatedAt     types.Int64         `tfsdk:"created_at"`
	UpdatedAt     types.Int64         `tfsdk:"updated_at"`
}

type RouteCaseModel struct {
	If               []RouteFilterModel `tfsdk:"if"`
	ChannelIDs       types.List         `tfsdk:"channel_ids"`
	Fallthrough      types.Bool         `tfsdk:"fallthrough"`
	RoutingMode      types.String       `tfsdk:"routing_mode"`
	NameMappingLabel types.String       `tfsdk:"name_mapping_label"`
}

type RouteFilterModel struct {
	Key  types.String `tfsdk:"key"`
	Oper types.String `tfsdk:"oper"`
	Vals types.List   `tfsdk:"vals"`
}

type RouteSectionModel struct {
	Name     types.String `tfsdk:"name"`
	Position types.Int64  `tfsdk:"position"`
}

type RouteDefaultModel struct {
	ChannelIDs types.List `tfsdk:"channel_ids"`
}

func (r *RouteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route"
}

func (r *RouteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages routing configuration for a shared alert integration in Flashduty.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the route (same as integration_id).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"integration_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the shared alert integration.",
				Required:            true,
			},
			"cases": schema.ListNestedAttribute{
				MarkdownDescription: "Conditional routing rules.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"if": schema.ListNestedAttribute{
							MarkdownDescription: "Filter conditions (AND relationship).",
							Required:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										MarkdownDescription: "The attribute or label key (e.g., title, severity, labels.xxx).",
										Required:            true,
									},
									"oper": schema.StringAttribute{
										MarkdownDescription: "The operator (IN, NOTIN).",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("IN", "NOTIN"),
										},
									},
									"vals": schema.ListAttribute{
										MarkdownDescription: "The values to match.",
										Required:            true,
										ElementType:         types.StringType,
									},
								},
							},
						},
						"channel_ids": schema.ListAttribute{
							MarkdownDescription: "Target channel IDs (required for standard routing mode).",
							Optional:            true,
							ElementType:         types.Int64Type,
						},
						"fallthrough": schema.BoolAttribute{
							MarkdownDescription: "Whether to continue matching after this rule.",
							Optional:            true,
						},
						"routing_mode": schema.StringAttribute{
							MarkdownDescription: "Routing mode (standard, name_mapping).",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("standard", "name_mapping"),
							},
						},
						"name_mapping_label": schema.StringAttribute{
							MarkdownDescription: "Label key for name mapping mode (e.g., labels.ChannelName).",
							Optional:            true,
						},
					},
				},
			},
			"sections": schema.ListNestedAttribute{
				MarkdownDescription: "Section dividers in route rules.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Section name (must be unique).",
							Required:            true,
						},
						"position": schema.Int64Attribute{
							MarkdownDescription: "Position (0 means before case 0).",
							Required:            true,
						},
					},
				},
			},
			"default": schema.SingleNestedAttribute{
				MarkdownDescription: "Default routing configuration.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"channel_ids": schema.ListAttribute{
						MarkdownDescription: "Default target channel IDs.",
						Required:            true,
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

func (r *RouteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (r *RouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RouteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upsertReq := r.buildUpsertRequest(ctx, &plan, 0, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpsertRoute(ctx, upsertReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Route", fmt.Sprintf("Could not create route: %s", err))
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(plan.IntegrationID.ValueInt64(), 10))

	// Only read back computed fields, preserve plan values for user-configured fields
	route, err := r.client.GetRoute(ctx, &client.GetRouteRequest{IntegrationID: plan.IntegrationID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Route", fmt.Sprintf("Could not read route after creation: %s", err))
		return
	}

	plan.Version = types.Int64Value(int64(route.Version))
	plan.Status = types.StringValue(route.Status)
	plan.CreatedAt = types.Int64Value(route.CreatedAt)
	plan.UpdatedAt = types.Int64Value(route.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RouteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	route, err := r.client.GetRoute(ctx, &client.GetRouteRequest{IntegrationID: state.IntegrationID.ValueInt64()})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Route", fmt.Sprintf("Could not read route: %s", err))
		return
	}

	r.mapRouteToModel(ctx, route, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RouteResourceModel
	var state RouteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	currentVersion := int(state.Version.ValueInt64())
	upsertReq := r.buildUpsertRequest(ctx, &plan, currentVersion, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpsertRoute(ctx, upsertReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Route", fmt.Sprintf("Could not update route: %s", err))
		return
	}

	// Only read back computed fields, preserve plan values for user-configured fields
	route, err := r.client.GetRoute(ctx, &client.GetRouteRequest{IntegrationID: plan.IntegrationID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Route", fmt.Sprintf("Could not read route after update: %s", err))
		return
	}

	plan.Version = types.Int64Value(int64(route.Version))
	plan.Status = types.StringValue(route.Status)
	plan.CreatedAt = types.Int64Value(route.CreatedAt)
	plan.UpdatedAt = types.Int64Value(route.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Route cannot be deleted via API - it can only be updated.
	// Removing from Terraform state only. The route configuration
	// will remain on the integration until manually changed.
	resp.Diagnostics.AddWarning(
		"Route Not Deleted",
		"Route resources cannot be deleted via API. The route has been removed from Terraform state but still exists on the integration.",
	)
}

func (r *RouteResource) buildUpsertRequest(ctx context.Context, plan *RouteResourceModel, version int, diags *diag.Diagnostics) *client.UpsertRouteRequest {
	upsertReq := &client.UpsertRouteRequest{
		IntegrationID: plan.IntegrationID.ValueInt64(),
		Version:       version,
	}

	// Build cases
	if len(plan.Cases) > 0 {
		for _, caseModel := range plan.Cases {
			routeCase := client.RouteCase{
				Fallthrough:      caseModel.Fallthrough.ValueBool(),
				RoutingMode:      caseModel.RoutingMode.ValueString(),
				NameMappingLabel: caseModel.NameMappingLabel.ValueString(),
			}

			// Build filters
			for _, filterModel := range caseModel.If {
				filter := client.RouteFilter{
					Key:  filterModel.Key.ValueString(),
					Oper: filterModel.Oper.ValueString(),
				}

				var vals []string
				diags.Append(filterModel.Vals.ElementsAs(ctx, &vals, false)...)
				filter.Vals = vals

				routeCase.If = append(routeCase.If, filter)
			}

			// Build channel IDs
			if !caseModel.ChannelIDs.IsNull() && !caseModel.ChannelIDs.IsUnknown() {
				var channelIDs []int64
				diags.Append(caseModel.ChannelIDs.ElementsAs(ctx, &channelIDs, false)...)
				routeCase.ChannelIDs = channelIDs
			}

			upsertReq.Cases = append(upsertReq.Cases, routeCase)
		}
	}

	// Build sections
	if len(plan.Sections) > 0 {
		for _, sectionModel := range plan.Sections {
			section := client.RouteSection{
				Name:     sectionModel.Name.ValueString(),
				Position: int(sectionModel.Position.ValueInt64()),
			}
			upsertReq.Sections = append(upsertReq.Sections, section)
		}
	}

	// Build default
	if plan.Default != nil {
		var channelIDs []int64
		diags.Append(plan.Default.ChannelIDs.ElementsAs(ctx, &channelIDs, false)...)
		upsertReq.Default = &client.RouteDefault{
			ChannelIDs: channelIDs,
		}
	}

	return upsertReq
}

func (r *RouteResource) mapRouteToModel(ctx context.Context, route *client.Route, model *RouteResourceModel, diags *diag.Diagnostics) {
	model.Version = types.Int64Value(int64(route.Version))
	model.Status = types.StringValue(route.Status)
	model.CreatedAt = types.Int64Value(route.CreatedAt)
	model.UpdatedAt = types.Int64Value(route.UpdatedAt)

	// Map cases
	if len(route.Cases) > 0 {
		model.Cases = make([]RouteCaseModel, len(route.Cases))
		for i, routeCase := range route.Cases {
			caseModel := RouteCaseModel{}

			if routeCase.Fallthrough {
				caseModel.Fallthrough = types.BoolValue(true)
			}
			if routeCase.RoutingMode != "" {
				caseModel.RoutingMode = types.StringValue(routeCase.RoutingMode)
			}
			if routeCase.NameMappingLabel != "" {
				caseModel.NameMappingLabel = types.StringValue(routeCase.NameMappingLabel)
			}

			// Map filters
			caseModel.If = make([]RouteFilterModel, len(routeCase.If))
			for j, filter := range routeCase.If {
				vals, d := types.ListValueFrom(ctx, types.StringType, filter.Vals)
				diags.Append(d...)

				caseModel.If[j] = RouteFilterModel{
					Key:  types.StringValue(filter.Key),
					Oper: types.StringValue(filter.Oper),
					Vals: vals,
				}
			}

			// Map channel IDs
			if len(routeCase.ChannelIDs) > 0 {
				channelIDs, d := types.ListValueFrom(ctx, types.Int64Type, routeCase.ChannelIDs)
				diags.Append(d...)
				caseModel.ChannelIDs = channelIDs
			} else {
				caseModel.ChannelIDs = types.ListNull(types.Int64Type)
			}

			model.Cases[i] = caseModel
		}
	} else {
		model.Cases = nil
	}

	// Map sections
	if len(route.Sections) > 0 {
		model.Sections = make([]RouteSectionModel, len(route.Sections))
		for i, section := range route.Sections {
			model.Sections[i] = RouteSectionModel{
				Name:     types.StringValue(section.Name),
				Position: types.Int64Value(int64(section.Position)),
			}
		}
	} else {
		model.Sections = nil
	}

	// Map default
	if route.Default != nil && len(route.Default.ChannelIDs) > 0 {
		channelIDs, d := types.ListValueFrom(ctx, types.Int64Type, route.Default.ChannelIDs)
		diags.Append(d...)
		model.Default = &RouteDefaultModel{
			ChannelIDs: channelIDs,
		}
	} else {
		model.Default = nil
	}
}
