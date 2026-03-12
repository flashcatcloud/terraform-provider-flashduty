package client

import (
	"context"
	"net/http"
)

// InhibitRule represents a Flashduty inhibit rule.
type InhibitRule struct {
	RuleID            string     `json:"rule_id"`
	RuleName          string     `json:"rule_name"`
	ChannelID         int64      `json:"channel_id,omitempty"`
	Description       string     `json:"description,omitempty"`
	Priority          int        `json:"priority,omitempty"`
	SourceFilters     [][]Filter `json:"source_filters,omitempty"`
	TargetFilters     [][]Filter `json:"target_filters,omitempty"`
	Equals            []string   `json:"equals,omitempty"`
	IsDirectlyDiscard bool       `json:"is_directly_discard,omitempty"`
	Status            string     `json:"status,omitempty"`
	CreatedAt         int64      `json:"created_at,omitempty"`
	UpdatedAt         int64      `json:"updated_at,omitempty"`
}

// CreateInhibitRuleRequest represents the request body for creating an inhibit rule.
type CreateInhibitRuleRequest struct {
	ChannelID         int64      `json:"channel_id"`
	RuleName          string     `json:"rule_name"`
	Description       string     `json:"description"`
	Priority          int        `json:"priority,omitempty"`
	SourceFilters     [][]Filter `json:"source_filters"`
	TargetFilters     [][]Filter `json:"target_filters"`
	Equals            []string   `json:"equals"`
	IsDirectlyDiscard bool       `json:"is_directly_discard,omitempty"`
}

// CreateInhibitRuleResult represents the response data for inhibit rule creation.
type CreateInhibitRuleResult struct {
	RuleID   string `json:"rule_id"`
	RuleName string `json:"rule_name"`
}

// UpdateInhibitRuleRequest represents the request body for updating an inhibit rule.
type UpdateInhibitRuleRequest struct {
	ChannelID         int64      `json:"channel_id"`
	RuleID            string     `json:"rule_id"`
	RuleName          string     `json:"rule_name,omitempty"`
	Description       string     `json:"description,omitempty"`
	Priority          *int       `json:"priority,omitempty"`
	SourceFilters     [][]Filter `json:"source_filters,omitempty"`
	TargetFilters     [][]Filter `json:"target_filters,omitempty"`
	Equals            []string   `json:"equals,omitempty"`
	IsDirectlyDiscard *bool      `json:"is_directly_discard,omitempty"`
}

// GetInhibitRuleRequest represents the request body for getting inhibit rule info.
type GetInhibitRuleRequest struct {
	ChannelID int64  `json:"channel_id"`
	RuleID    string `json:"rule_id"`
}

// DeleteInhibitRuleRequest represents the request body for deleting an inhibit rule.
type DeleteInhibitRuleRequest struct {
	ChannelID int64  `json:"channel_id"`
	RuleID    string `json:"rule_id"`
}

// EnableInhibitRuleRequest represents the request body for enabling an inhibit rule.
type EnableInhibitRuleRequest struct {
	ChannelID int64  `json:"channel_id"`
	RuleID    string `json:"rule_id"`
}

// ListInhibitRulesRequest represents the request body for listing inhibit rules.
type ListInhibitRulesRequest struct {
	ChannelID int64 `json:"channel_id"`
}

// ListInhibitRulesResult represents the response data for listing inhibit rules.
type ListInhibitRulesResult struct {
	Items []InhibitRule `json:"items"`
}

// CreateInhibitRule creates a new inhibit rule.
func (c *Client) CreateInhibitRule(ctx context.Context, req *CreateInhibitRuleRequest) (*CreateInhibitRuleResult, error) {
	result, _, err := doRequestWithResponse[CreateInhibitRuleResult](c, ctx, http.MethodPost, "/channel/inhibit/rule/create", req)
	return result, err
}

// ListInhibitRules lists all inhibit rules for a channel.
func (c *Client) ListInhibitRules(ctx context.Context, req *ListInhibitRulesRequest) (*ListInhibitRulesResult, error) {
	result, _, err := doRequestWithResponse[ListInhibitRulesResult](c, ctx, http.MethodPost, "/channel/inhibit/rule/list", req)
	return result, err
}

// GetInhibitRule retrieves inhibit rule information by listing all rules and finding the target.
func (c *Client) GetInhibitRule(ctx context.Context, req *GetInhibitRuleRequest) (*InhibitRule, error) {
	result, err := c.ListInhibitRules(ctx, &ListInhibitRulesRequest{ChannelID: req.ChannelID})
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

// UpdateInhibitRule updates an inhibit rule.
func (c *Client) UpdateInhibitRule(ctx context.Context, req *UpdateInhibitRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/inhibit/rule/update", req)
	return err
}

// DeleteInhibitRule deletes an inhibit rule.
func (c *Client) DeleteInhibitRule(ctx context.Context, req *DeleteInhibitRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/inhibit/rule/delete", req)
	return err
}

// EnableInhibitRule enables an inhibit rule.
func (c *Client) EnableInhibitRule(ctx context.Context, req *EnableInhibitRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/inhibit/rule/enable", req)
	return err
}

// DisableInhibitRule disables an inhibit rule.
func (c *Client) DisableInhibitRule(ctx context.Context, req *EnableInhibitRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/inhibit/rule/disable", req)
	return err
}
