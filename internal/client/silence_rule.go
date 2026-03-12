package client

import (
	"context"
	"net/http"
)

// SilenceRule represents a Flashduty silence rule.
type SilenceRule struct {
	RuleID            string       `json:"rule_id"`
	RuleName          string       `json:"rule_name"`
	ChannelID         int64        `json:"channel_id,omitempty"`
	Description       string       `json:"description,omitempty"`
	Priority          int          `json:"priority,omitempty"`
	Filters           [][]Filter   `json:"filters,omitempty"`
	TimeFilter        *SingleTime  `json:"time_filter,omitempty"`
	TimeFilters       []TimeFilter `json:"time_filters,omitempty"`
	IsDirectlyDiscard bool         `json:"is_directly_discard,omitempty"`
	FromIncidentID    string       `json:"from_incident_id,omitempty"`
	Status            string       `json:"status,omitempty"`
	CreatedAt         int64        `json:"created_at,omitempty"`
	UpdatedAt         int64        `json:"updated_at,omitempty"`
}

// SingleTime represents a single time range.
type SingleTime struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// CreateSilenceRuleRequest represents the request body for creating a silence rule.
type CreateSilenceRuleRequest struct {
	ChannelID         int64        `json:"channel_id"`
	RuleName          string       `json:"rule_name"`
	Description       string       `json:"description"`
	Priority          int          `json:"priority,omitempty"`
	Filters           [][]Filter   `json:"filters"`
	TimeFilter        *SingleTime  `json:"time_filter,omitempty"`
	TimeFilters       []TimeFilter `json:"time_filters,omitempty"`
	IsDirectlyDiscard bool         `json:"is_directly_discard,omitempty"`
	FromIncidentID    string       `json:"from_incident_id,omitempty"`
}

// CreateSilenceRuleResult represents the response data for silence rule creation.
type CreateSilenceRuleResult struct {
	RuleID   string `json:"rule_id"`
	RuleName string `json:"rule_name"`
}

// UpdateSilenceRuleRequest represents the request body for updating a silence rule.
type UpdateSilenceRuleRequest struct {
	ChannelID         int64        `json:"channel_id"`
	RuleID            string       `json:"rule_id"`
	RuleName          string       `json:"rule_name,omitempty"`
	Description       string       `json:"description,omitempty"`
	Priority          *int         `json:"priority,omitempty"`
	Filters           [][]Filter   `json:"filters,omitempty"`
	TimeFilter        *SingleTime  `json:"time_filter,omitempty"`
	TimeFilters       []TimeFilter `json:"time_filters,omitempty"`
	IsDirectlyDiscard *bool        `json:"is_directly_discard,omitempty"`
}

// GetSilenceRuleRequest represents the request body for getting silence rule info.
type GetSilenceRuleRequest struct {
	ChannelID int64  `json:"channel_id"`
	RuleID    string `json:"rule_id"`
}

// DeleteSilenceRuleRequest represents the request body for deleting a silence rule.
type DeleteSilenceRuleRequest struct {
	ChannelID int64  `json:"channel_id"`
	RuleID    string `json:"rule_id"`
}

// EnableSilenceRuleRequest represents the request body for enabling a silence rule.
type EnableSilenceRuleRequest struct {
	ChannelID int64  `json:"channel_id"`
	RuleID    string `json:"rule_id"`
}

// ListSilenceRulesRequest represents the request body for listing silence rules.
type ListSilenceRulesRequest struct {
	ChannelID int64 `json:"channel_id"`
}

// ListSilenceRulesResult represents the response data for listing silence rules.
type ListSilenceRulesResult struct {
	Items []SilenceRule `json:"items"`
}

// CreateSilenceRule creates a new silence rule.
func (c *Client) CreateSilenceRule(ctx context.Context, req *CreateSilenceRuleRequest) (*CreateSilenceRuleResult, error) {
	result, _, err := doRequestWithResponse[CreateSilenceRuleResult](c, ctx, http.MethodPost, "/channel/silence/rule/create", req)
	return result, err
}

// ListSilenceRules lists all silence rules for a channel.
func (c *Client) ListSilenceRules(ctx context.Context, req *ListSilenceRulesRequest) (*ListSilenceRulesResult, error) {
	result, _, err := doRequestWithResponse[ListSilenceRulesResult](c, ctx, http.MethodPost, "/channel/silence/rule/list", req)
	return result, err
}

// GetSilenceRule retrieves silence rule information by listing all rules and finding the target.
func (c *Client) GetSilenceRule(ctx context.Context, req *GetSilenceRuleRequest) (*SilenceRule, error) {
	result, err := c.ListSilenceRules(ctx, &ListSilenceRulesRequest{ChannelID: req.ChannelID})
	if err != nil {
		return nil, err
	}
	if result != nil {
		for i := range result.Items {
			if result.Items[i].RuleID == req.RuleID {
				return &result.Items[i], nil
			}
		}
	}
	return nil, ErrNotFound
}

// UpdateSilenceRule updates a silence rule.
func (c *Client) UpdateSilenceRule(ctx context.Context, req *UpdateSilenceRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/silence/rule/update", req)
	return err
}

// DeleteSilenceRule deletes a silence rule.
func (c *Client) DeleteSilenceRule(ctx context.Context, req *DeleteSilenceRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/silence/rule/delete", req)
	return err
}

// EnableSilenceRule enables a silence rule.
func (c *Client) EnableSilenceRule(ctx context.Context, req *EnableSilenceRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/silence/rule/enable", req)
	return err
}

// DisableSilenceRule disables a silence rule.
func (c *Client) DisableSilenceRule(ctx context.Context, req *EnableSilenceRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/silence/rule/disable", req)
	return err
}
