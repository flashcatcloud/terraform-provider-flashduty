package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &FieldDataSource{}
	_ datasource.DataSourceWithConfigure = &FieldDataSource{}
)

func NewFieldDataSource() datasource.DataSource {
	return &FieldDataSource{}
}

type FieldDataSource struct {
	client *client.Client
}

type FieldDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	FieldName    types.String `tfsdk:"field_name"`
	DisplayName  types.String `tfsdk:"display_name"`
	Description  types.String `tfsdk:"description"`
	FieldType    types.String `tfsdk:"field_type"`
	ValueType    types.String `tfsdk:"value_type"`
	Options      types.String `tfsdk:"options"`
	DefaultValue types.String `tfsdk:"default_value"`
	Status       types.String `tfsdk:"status"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
	UpdatedAt    types.Int64  `tfsdk:"updated_at"`
}

func (d *FieldDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_field"
}

func (d *FieldDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a custom field from Flashduty.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the field.",
				Required:            true,
			},
			"field_name": schema.StringAttribute{
				MarkdownDescription: "The name of the field.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the field.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the field.",
				Computed:            true,
			},
			"field_type": schema.StringAttribute{
				MarkdownDescription: "The type of the field.",
				Computed:            true,
			},
			"value_type": schema.StringAttribute{
				MarkdownDescription: "The value type of the field.",
				Computed:            true,
			},
			"options": schema.StringAttribute{
				MarkdownDescription: "The available options as JSON array.",
				Computed:            true,
			},
			"default_value": schema.StringAttribute{
				MarkdownDescription: "The default value as JSON.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the field.",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The timestamp when the field was created.",
				Computed:            true,
			},
			"updated_at": schema.Int64Attribute{
				MarkdownDescription: "The timestamp when the field was last updated.",
				Computed:            true,
			},
		},
	}
}

func (d *FieldDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *FieldDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state FieldDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	field, err := d.client.GetField(ctx, &client.GetFieldRequest{FieldID: state.ID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Field", fmt.Sprintf("Could not read field: %s", err))
		return
	}

	state.ID = types.StringValue(field.FieldID)
	state.FieldName = types.StringValue(field.FieldName)
	state.DisplayName = types.StringValue(field.DisplayName)
	state.Description = types.StringValue(field.Description)
	state.FieldType = types.StringValue(field.FieldType)
	state.ValueType = types.StringValue(field.ValueType)
	state.Status = types.StringValue(field.Status)
	state.CreatedAt = types.Int64Value(field.CreatedAt)
	state.UpdatedAt = types.Int64Value(field.UpdatedAt)

	if len(field.Options) > 0 {
		optionsJSON, _ := json.Marshal(field.Options)
		state.Options = types.StringValue(string(optionsJSON))
	} else {
		state.Options = types.StringNull()
	}

	if field.DefaultValue != nil {
		defaultValueJSON, _ := json.Marshal(field.DefaultValue)
		state.DefaultValue = types.StringValue(string(defaultValueJSON))
	} else {
		state.DefaultValue = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// FieldsDataSource - list all fields.
var (
	_ datasource.DataSource              = &FieldsDataSource{}
	_ datasource.DataSourceWithConfigure = &FieldsDataSource{}
)

func NewFieldsDataSource() datasource.DataSource {
	return &FieldsDataSource{}
}

type FieldsDataSource struct {
	client *client.Client
}

type FieldsDataSourceModel struct {
	Fields []FieldDataSourceModel `tfsdk:"fields"`
}

func (d *FieldsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fields"
}

func (d *FieldsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves all custom fields from Flashduty.",
		Attributes: map[string]schema.Attribute{
			"fields": schema.ListNestedAttribute{
				MarkdownDescription: "List of custom fields.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the field.",
							Computed:            true,
						},
						"field_name": schema.StringAttribute{
							MarkdownDescription: "The name of the field.",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the field.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the field.",
							Computed:            true,
						},
						"field_type": schema.StringAttribute{
							MarkdownDescription: "The type of the field.",
							Computed:            true,
						},
						"value_type": schema.StringAttribute{
							MarkdownDescription: "The value type of the field.",
							Computed:            true,
						},
						"options": schema.StringAttribute{
							MarkdownDescription: "The available options as JSON array.",
							Computed:            true,
						},
						"default_value": schema.StringAttribute{
							MarkdownDescription: "The default value as JSON.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "The status of the field.",
							Computed:            true,
						},
						"created_at": schema.Int64Attribute{
							MarkdownDescription: "The timestamp when the field was created.",
							Computed:            true,
						},
						"updated_at": schema.Int64Attribute{
							MarkdownDescription: "The timestamp when the field was last updated.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *FieldsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (d *FieldsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	result, err := d.client.ListFields(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Fields", fmt.Sprintf("Could not read fields: %s", err))
		return
	}

	var state FieldsDataSourceModel
	for _, field := range result.Items {
		fieldModel := FieldDataSourceModel{
			ID:          types.StringValue(field.FieldID),
			FieldName:   types.StringValue(field.FieldName),
			DisplayName: types.StringValue(field.DisplayName),
			Description: types.StringValue(field.Description),
			FieldType:   types.StringValue(field.FieldType),
			ValueType:   types.StringValue(field.ValueType),
			Status:      types.StringValue(field.Status),
			CreatedAt:   types.Int64Value(field.CreatedAt),
			UpdatedAt:   types.Int64Value(field.UpdatedAt),
		}

		if len(field.Options) > 0 {
			optionsJSON, _ := json.Marshal(field.Options)
			fieldModel.Options = types.StringValue(string(optionsJSON))
		} else {
			fieldModel.Options = types.StringNull()
		}

		if field.DefaultValue != nil {
			defaultValueJSON, _ := json.Marshal(field.DefaultValue)
			fieldModel.DefaultValue = types.StringValue(string(defaultValueJSON))
		} else {
			fieldModel.DefaultValue = types.StringNull()
		}

		state.Fields = append(state.Fields, fieldModel)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
