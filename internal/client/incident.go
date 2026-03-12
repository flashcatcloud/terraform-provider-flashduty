package client

import (
	"context"
	"net/http"
)

// Incident represents a Flashduty incident.
type Incident struct {
	IncidentID       string              `json:"incident_id"`
	Title            string              `json:"title,omitempty"`
	Description      string              `json:"description,omitempty"`
	IncidentSeverity string              `json:"incident_severity,omitempty"`
	IncidentStatus   string              `json:"incident_status,omitempty"`
	Progress         string              `json:"progress,omitempty"`
	ChannelID        int64               `json:"channel_id,omitempty"`
	ChannelName      string              `json:"channel_name,omitempty"`
	StartTime        int64               `json:"start_time,omitempty"`
	EndTime          int64               `json:"end_time,omitempty"`
	AckTime          int64               `json:"ack_time,omitempty"`
	CloseTime        int64               `json:"close_time,omitempty"`
	SnoozedBefore    int64               `json:"snoozed_before,omitempty"`
	CreatorID        int64               `json:"creator_id,omitempty"`
	CloserID         int64               `json:"closer_id,omitempty"`
	AssignedTo       *IncidentAssignment `json:"assigned_to,omitempty"`
	Impact           string              `json:"impact,omitempty"`
	RootCause        string              `json:"root_cause,omitempty"`
	Resolution       string              `json:"resolution,omitempty"`
	AlertCnt         int                 `json:"alert_cnt,omitempty"`
	CreatedAt        int64               `json:"created_at,omitempty"`
	UpdatedAt        int64               `json:"updated_at,omitempty"`
}

// IncidentAssignment represents assignment information for an incident.
type IncidentAssignment struct {
	Type           string  `json:"type,omitempty"`
	PersonIDs      []int64 `json:"person_ids,omitempty"`
	EscalateRuleID string  `json:"escalate_rule_id,omitempty"`
	LayerIdx       int     `json:"layer_idx,omitempty"`
}

// CreateIncidentRequest represents the request body for creating an incident.
type CreateIncidentRequest struct {
	Title            string              `json:"title"`
	Description      string              `json:"description,omitempty"`
	IncidentSeverity string              `json:"incident_severity"`
	ChannelID        int64               `json:"channel_id,omitempty"`
	AssignedTo       *IncidentAssignment `json:"assigned_to,omitempty"`
}

// CreateIncidentResult represents the response data for incident creation.
type CreateIncidentResult struct {
	IncidentID string `json:"incident_id"`
	Title      string `json:"title"`
}

// GetIncidentRequest represents the request body for getting incident info.
type GetIncidentRequest struct {
	IncidentID string `json:"incident_id"`
}

// DeleteIncidentRequest represents the request body for deleting incidents.
type DeleteIncidentRequest struct {
	IncidentIDs []string `json:"incident_ids"`
}

// UpdateIncidentRequest represents the request body for updating an incident.
type UpdateIncidentRequest struct {
	IncidentID       string `json:"incident_id"`
	Title            string `json:"title,omitempty"`
	Description      string `json:"description,omitempty"`
	IncidentSeverity string `json:"incident_severity,omitempty"`
	Impact           string `json:"impact,omitempty"`
	RootCause        string `json:"root_cause,omitempty"`
	Resolution       string `json:"resolution,omitempty"`
}

// CreateIncident creates a new incident.
func (c *Client) CreateIncident(ctx context.Context, req *CreateIncidentRequest) (*CreateIncidentResult, error) {
	result, _, err := doRequestWithResponse[CreateIncidentResult](c, ctx, http.MethodPost, "/incident/create", req)
	return result, err
}

// GetIncident retrieves incident information.
func (c *Client) GetIncident(ctx context.Context, req *GetIncidentRequest) (*Incident, error) {
	result, _, err := doRequestWithResponse[Incident](c, ctx, http.MethodPost, "/incident/info", req)
	return result, err
}

// UpdateIncident updates an incident's details.
func (c *Client) UpdateIncident(ctx context.Context, req *UpdateIncidentRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/incident/reset", req)
	return err
}

// DeleteIncident deletes incidents by IDs.
func (c *Client) DeleteIncident(ctx context.Context, req *DeleteIncidentRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/incident/remove", req)
	return err
}
