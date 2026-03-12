package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &FieldResource{}
	_ resource.ResourceWithConfigure   = &FieldResource{}
	_ resource.ResourceWithImportState = &FieldResource{}
)

func NewFieldResource() resource.Resource {
	return &FieldResource{}
}

type FieldResource struct {
	client *client.Client
}

type FieldResourceModel struct {
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

func (r *FieldResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_field"
}

func (r *FieldResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a custom field in Flashduty.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the field.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"field_name": schema.StringAttribute{
				MarkdownDescription: "The name of the field (must match ^[a-zA-Z_][a-zA-Z0-9_]{0,39}$).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the field.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the field.",
				Optional:            true,
			},
			"field_type": schema.StringAttribute{
				MarkdownDescription: "The type of the field (text, single_select, multi_select, checkbox).",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("text", "single_select", "multi_select", "checkbox"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value_type": schema.StringAttribute{
				MarkdownDescription: "The value type of the field (string, bool).",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("string", "bool"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"options": schema.StringAttribute{
				MarkdownDescription: "The available options as JSON array (for single_select, multi_select types).",
				Optional:            true,
			},
			"default_value": schema.StringAttribute{
				MarkdownDescription: "The default value as JSON.",
				Optional:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the field (enabled, disabled, deleted).",
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

func (r *FieldResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureClient(req.ProviderData, &resp.Diagnostics)
}

func (r *FieldResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FieldResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateFieldRequest{
		FieldName:   plan.FieldName.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
		FieldType:   plan.FieldType.ValueString(),
		ValueType:   plan.ValueType.ValueString(),
	}

	if !plan.Options.IsNull() && !plan.Options.IsUnknown() {
		var options []any
		if err := json.Unmarshal([]byte(plan.Options.ValueString()), &options); err != nil {
			resp.Diagnostics.AddError("Invalid Options", fmt.Sprintf("Failed to parse options JSON: %s", err))
			return
		}
		createReq.Options = options
	}

	if !plan.DefaultValue.IsNull() && !plan.DefaultValue.IsUnknown() {
		var defaultValue any
		if err := json.Unmarshal([]byte(plan.DefaultValue.ValueString()), &defaultValue); err != nil {
			resp.Diagnostics.AddError("Invalid Default Value", fmt.Sprintf("Failed to parse default_value JSON: %s", err))
			return
		}
		createReq.DefaultValue = defaultValue
	}

	result, err := r.client.CreateField(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Field", fmt.Sprintf("Could not create field: %s", err))
		return
	}

	plan.ID = types.StringValue(result.FieldID)

	// Read back the field to get computed values
	field, err := r.client.GetField(ctx, &client.GetFieldRequest{FieldID: result.FieldID})
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Field", fmt.Sprintf("Could not read field after creation: %s", err))
		return
	}

	r.mapFieldToModel(field, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FieldResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FieldResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	field, err := r.client.GetField(ctx, &client.GetFieldRequest{FieldID: state.ID.ValueString()})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Field", fmt.Sprintf("Could not read field: %s", err))
		return
	}

	r.mapFieldToModel(field, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FieldResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FieldResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &client.UpdateFieldRequest{
		FieldID:     plan.ID.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	if !plan.Options.IsNull() && !plan.Options.IsUnknown() {
		var options []any
		if err := json.Unmarshal([]byte(plan.Options.ValueString()), &options); err != nil {
			resp.Diagnostics.AddError("Invalid Options", fmt.Sprintf("Failed to parse options JSON: %s", err))
			return
		}
		updateReq.Options = options
	}

	if !plan.DefaultValue.IsNull() && !plan.DefaultValue.IsUnknown() {
		var defaultValue any
		if err := json.Unmarshal([]byte(plan.DefaultValue.ValueString()), &defaultValue); err != nil {
			resp.Diagnostics.AddError("Invalid Default Value", fmt.Sprintf("Failed to parse default_value JSON: %s", err))
			return
		}
		updateReq.DefaultValue = defaultValue
	}

	err := r.client.UpdateField(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Field", fmt.Sprintf("Could not update field: %s", err))
		return
	}

	// Read back the field to get updated values
	field, err := r.client.GetField(ctx, &client.GetFieldRequest{FieldID: plan.ID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Field", fmt.Sprintf("Could not read field after update: %s", err))
		return
	}

	r.mapFieldToModel(field, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FieldResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FieldResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteField(ctx, &client.DeleteFieldRequest{FieldID: state.ID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Field", fmt.Sprintf("Could not delete field: %s", err))
		return
	}
}

func (r *FieldResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *FieldResource) mapFieldToModel(field *client.Field, model *FieldResourceModel, diags *diag.Diagnostics) {
	model.ID = types.StringValue(field.FieldID)
	model.FieldName = types.StringValue(field.FieldName)
	model.DisplayName = types.StringValue(field.DisplayName)
	model.Description = types.StringValue(field.Description)
	model.FieldType = types.StringValue(field.FieldType)
	model.ValueType = types.StringValue(field.ValueType)
	model.Status = types.StringValue(field.Status)
	model.CreatedAt = types.Int64Value(field.CreatedAt)
	model.UpdatedAt = types.Int64Value(field.UpdatedAt)

	if len(field.Options) > 0 {
		optionsJSON, err := json.Marshal(field.Options)
		if err != nil {
			diags.AddError("Error Serializing Options", err.Error())
			return
		}
		model.Options = types.StringValue(string(optionsJSON))
	}

	if field.DefaultValue != nil {
		defaultValueJSON, err := json.Marshal(field.DefaultValue)
		if err != nil {
			diags.AddError("Error Serializing Default Value", err.Error())
			return
		}
		model.DefaultValue = types.StringValue(string(defaultValueJSON))
	}
}
