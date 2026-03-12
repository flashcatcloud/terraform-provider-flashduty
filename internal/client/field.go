package client

import (
	"context"
	"net/http"
)

// Field represents a custom field in Flashduty.
type Field struct {
	FieldID      string `json:"field_id,omitempty"`
	AccountID    int64  `json:"account_id,omitempty"`
	FieldName    string `json:"field_name,omitempty"`
	DisplayName  string `json:"display_name,omitempty"`
	Description  string `json:"description,omitempty"`
	FieldType    string `json:"field_type,omitempty"`
	ValueType    string `json:"value_type,omitempty"`
	Options      []any  `json:"options,omitempty"`
	DefaultValue any    `json:"default_value,omitempty"`
	Status       string `json:"status,omitempty"`
	CreatorID    int64  `json:"creator_id,omitempty"`
	UpdatedBy    int64  `json:"updated_by,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
	UpdatedAt    int64  `json:"updated_at,omitempty"`
}

// CreateFieldRequest represents a request to create a custom field.
type CreateFieldRequest struct {
	FieldName    string `json:"field_name"`
	DisplayName  string `json:"display_name"`
	Description  string `json:"description,omitempty"`
	FieldType    string `json:"field_type"`
	ValueType    string `json:"value_type"`
	Options      []any  `json:"options,omitempty"`
	DefaultValue any    `json:"default_value,omitempty"`
}

// CreateFieldResult represents the result of creating a custom field.
type CreateFieldResult struct {
	FieldID   string `json:"field_id"`
	FieldName string `json:"field_name"`
}

// UpdateFieldRequest represents a request to update a custom field.
type UpdateFieldRequest struct {
	FieldID      string `json:"field_id"`
	DisplayName  string `json:"display_name,omitempty"`
	Description  string `json:"description,omitempty"`
	Options      []any  `json:"options,omitempty"`
	DefaultValue any    `json:"default_value,omitempty"`
}

// GetFieldRequest represents a request to get a custom field.
type GetFieldRequest struct {
	FieldID string `json:"field_id"`
}

// DeleteFieldRequest represents a request to delete a custom field.
type DeleteFieldRequest struct {
	FieldID string `json:"field_id"`
}

// ListFieldsResult represents the result of listing custom fields.
type ListFieldsResult struct {
	Items []Field `json:"items"`
}

// CreateField creates a new custom field.
func (c *Client) CreateField(ctx context.Context, req *CreateFieldRequest) (*CreateFieldResult, error) {
	result, _, err := doRequestWithResponse[CreateFieldResult](c, ctx, http.MethodPost, "/field/create", req)
	return result, err
}

// GetField retrieves a custom field by ID.
func (c *Client) GetField(ctx context.Context, req *GetFieldRequest) (*Field, error) {
	result, _, err := doRequestWithResponse[Field](c, ctx, http.MethodPost, "/field/info", req)
	return result, err
}

// UpdateField updates a custom field.
func (c *Client) UpdateField(ctx context.Context, req *UpdateFieldRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/field/update", req)
	return err
}

// DeleteField deletes a custom field.
func (c *Client) DeleteField(ctx context.Context, req *DeleteFieldRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/field/delete", req)
	return err
}

// ListFields retrieves all custom fields.
func (c *Client) ListFields(ctx context.Context) (*ListFieldsResult, error) {
	result, _, err := doRequestWithResponse[ListFieldsResult](c, ctx, http.MethodPost, "/field/list", map[string]any{})
	return result, err
}
