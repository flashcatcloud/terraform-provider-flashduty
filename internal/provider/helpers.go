package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// configureClient extracts *client.Client from provider data.
// Shared by all resource and data source Configure methods.
func configureClient(providerData any, diags *diag.Diagnostics) *client.Client {
	if providerData == nil {
		return nil
	}

	c, ok := providerData.(*client.Client)
	if !ok {
		diags.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", providerData),
		)
		return nil
	}

	return c
}

// Int64ToString converts an int64 to a Terraform string value.
func Int64ToString(v int64) types.String {
	return types.StringValue(strconv.FormatInt(v, 10))
}

// StringToInt64 converts a Terraform string value to int64.
func StringToInt64(v types.String) (int64, error) {
	return strconv.ParseInt(v.ValueString(), 10, 64)
}

// Int64ListToSlice converts a Terraform List of Int64 to a Go []int64 slice.
func Int64ListToSlice(ctx context.Context, list types.List, diags *diag.Diagnostics) []int64 {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var result []int64
	diags.Append(list.ElementsAs(ctx, &result, false)...)
	return result
}

// StringListToSlice converts a Terraform List of String to a Go []string slice.
func StringListToSlice(ctx context.Context, list types.List, diags *diag.Diagnostics) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var result []string
	diags.Append(list.ElementsAs(ctx, &result, false)...)
	return result
}

// SliceToInt64List converts a Go []int64 slice to a Terraform List of Int64.
func SliceToInt64List(ctx context.Context, slice []int64, diags *diag.Diagnostics) types.List {
	if slice == nil {
		return types.ListNull(types.Int64Type)
	}

	list, d := types.ListValueFrom(ctx, types.Int64Type, slice)
	diags.Append(d...)
	return list
}

// SliceToStringList converts a Go []string slice to a Terraform List of String.
func SliceToStringList(ctx context.Context, slice []string, diags *diag.Diagnostics) types.List {
	if slice == nil {
		return types.ListNull(types.StringType)
	}

	list, d := types.ListValueFrom(ctx, types.StringType, slice)
	diags.Append(d...)
	return list
}
