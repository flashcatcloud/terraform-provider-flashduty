package client

import (
	"context"
	"net/http"
)

// AlertPipelineRule represents a single alert processing rule.
type AlertPipelineRule struct {
	Kind     string         `json:"kind"`
	If       []Filter       `json:"if,omitempty"`
	Settings map[string]any `json:"settings"`
}

// AlertPipeline represents the full alert processing pipeline for an integration.
type AlertPipeline struct {
	IntegrationID int64               `json:"integration_id"`
	Rules         []AlertPipelineRule `json:"rules"`
	CreatedAt     int64               `json:"created_at,omitempty"`
	UpdatedAt     int64               `json:"updated_at,omitempty"`
}

type GetAlertPipelineRequest struct {
	IntegrationID int64 `json:"integration_id"`
}

type UpsertAlertPipelineRequest struct {
	IntegrationID int64               `json:"integration_id"`
	Rules         []AlertPipelineRule `json:"rules"`
}

func (c *Client) GetAlertPipeline(ctx context.Context, req *GetAlertPipelineRequest) (*AlertPipeline, error) {
	result, _, err := doRequestWithResponse[AlertPipeline](c, ctx, http.MethodPost, "/alert/pipeline/info", req)
	return result, err
}

func (c *Client) UpsertAlertPipeline(ctx context.Context, req *UpsertAlertPipelineRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/alert/pipeline/upsert", req)
	return err
}
