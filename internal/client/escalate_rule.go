package client

import (
	"context"
	"net/http"
)

// EscalateRule represents a Flashduty escalation rule.
type EscalateRule struct {
	RuleID      string          `json:"rule_id"`
	RuleName    string          `json:"rule_name"`
	ChannelID   int64           `json:"channel_id,omitempty"`
	Description string          `json:"description,omitempty"`
	TemplateID  string          `json:"template_id,omitempty"`
	AggrWindow  int             `json:"aggr_window,omitempty"`
	Priority    int             `json:"priority,omitempty"`
	Layers      []EscalateLayer `json:"layers,omitempty"`
	TimeFilters []TimeFilter    `json:"time_filters,omitempty"`
	Filters     [][]Filter      `json:"filters,omitempty"`
	Status      string          `json:"status,omitempty"`
	CreatedAt   int64           `json:"created_at,omitempty"`
	UpdatedAt   int64           `json:"updated_at,omitempty"`
}

// EscalateLayer represents a layer in an escalation rule.
type EscalateLayer struct {
	MaxTimes       int             `json:"max_times,omitempty"`
	NotifyStep     float64         `json:"notify_step,omitempty"`
	EscalateWindow int             `json:"escalate_window,omitempty"`
	ForceEscalate  bool            `json:"force_escalate,omitempty"`
	Target         *EscalateTarget `json:"target,omitempty"`
}

// EscalateTargetBy represents notification preference settings.
type EscalateTargetBy struct {
	FollowPreference bool     `json:"follow_preference,omitempty"`
	Critical         []string `json:"critical,omitempty"`
	Warning          []string `json:"warning,omitempty"`
	Info             []string `json:"info,omitempty"`
}

// EscalateWebhook represents a webhook notification target.
type EscalateWebhook struct {
	Type     string         `json:"type,omitempty"`
	Settings map[string]any `json:"settings,omitempty"`
}

// EscalateTarget represents the target of an escalation.
type EscalateTarget struct {
	PersonIDs         []int64            `json:"person_ids,omitempty"`
	TeamIDs           []int64            `json:"team_ids,omitempty"`
	Emails            []string           `json:"emails,omitempty"`
	ScheduleToRoleIDs map[string][]int64 `json:"schedule_to_role_ids,omitempty"`
	By                *EscalateTargetBy  `json:"by,omitempty"`
	Webhooks          []EscalateWebhook  `json:"webhooks,omitempty"`
}

// TimeFilter represents a time filter condition.
type TimeFilter struct {
	Start  string `json:"start,omitempty"`
	End    string `json:"end,omitempty"`
	Repeat []int  `json:"repeat,omitempty"`
	CalID  string `json:"cal_id,omitempty"`
	IsOff  bool   `json:"is_off,omitempty"`
}

// Filter represents a filter condition.
type Filter struct {
	Key  string   `json:"key,omitempty"`
	Oper string   `json:"oper,omitempty"`
	Vals []string `json:"vals,omitempty"`
}

// CreateEscalateRuleRequest represents the request body for creating an escalate rule.
type CreateEscalateRuleRequest struct {
	ChannelID   int64           `json:"channel_id"`
	RuleName    string          `json:"rule_name"`
	TemplateID  string          `json:"template_id"`
	Description string          `json:"description,omitempty"`
	AggrWindow  int             `json:"aggr_window"`
	Priority    int             `json:"priority,omitempty"`
	Layers      []EscalateLayer `json:"layers"`
	TimeFilters []TimeFilter    `json:"time_filters,omitempty"`
	Filters     [][]Filter      `json:"filters,omitempty"`
}

// CreateEscalateRuleResult represents the response data for escalate rule creation.
type CreateEscalateRuleResult struct {
	RuleID   string `json:"rule_id"`
	RuleName string `json:"rule_name"`
}

// UpdateEscalateRuleRequest represents the request body for updating an escalate rule.
type UpdateEscalateRuleRequest struct {
	ChannelID   int64           `json:"channel_id"`
	RuleID      string          `json:"rule_id"`
	RuleName    string          `json:"rule_name,omitempty"`
	TemplateID  string          `json:"template_id,omitempty"`
	Description string          `json:"description,omitempty"`
	AggrWindow  *int            `json:"aggr_window,omitempty"`
	Priority    *int            `json:"priority,omitempty"`
	Layers      []EscalateLayer `json:"layers,omitempty"`
	TimeFilters []TimeFilter    `json:"time_filters,omitempty"`
	Filters     [][]Filter      `json:"filters,omitempty"`
}

// GetEscalateRuleRequest represents the request body for getting escalate rule info.
type GetEscalateRuleRequest struct {
	ChannelID int64  `json:"channel_id"`
	RuleID    string `json:"rule_id"`
}

// DeleteEscalateRuleRequest represents the request body for deleting an escalate rule.
type DeleteEscalateRuleRequest struct {
	ChannelID int64  `json:"channel_id"`
	RuleID    string `json:"rule_id"`
}

// EnableEscalateRuleRequest represents the request body for enabling an escalate rule.
type EnableEscalateRuleRequest struct {
	ChannelID int64  `json:"channel_id"`
	RuleID    string `json:"rule_id"`
}

// ListEscalateRulesRequest represents the request body for listing escalate rules.
type ListEscalateRulesRequest struct {
	ChannelID int64 `json:"channel_id"`
}

// ListEscalateRulesResult represents the response data for listing escalate rules.
type ListEscalateRulesResult struct {
	Items []EscalateRule `json:"items"`
}

// CreateEscalateRule creates a new escalate rule.
func (c *Client) CreateEscalateRule(ctx context.Context, req *CreateEscalateRuleRequest) (*CreateEscalateRuleResult, error) {
	result, _, err := doRequestWithResponse[CreateEscalateRuleResult](c, ctx, http.MethodPost, "/channel/escalate/rule/create", req)
	return result, err
}

// ListEscalateRules lists all escalate rules for a channel.
func (c *Client) ListEscalateRules(ctx context.Context, req *ListEscalateRulesRequest) (*ListEscalateRulesResult, error) {
	result, _, err := doRequestWithResponse[ListEscalateRulesResult](c, ctx, http.MethodPost, "/channel/escalate/rule/list", req)
	return result, err
}

// GetEscalateRule retrieves escalate rule information.
func (c *Client) GetEscalateRule(ctx context.Context, req *GetEscalateRuleRequest) (*EscalateRule, error) {
	result, _, err := doRequestWithResponse[EscalateRule](c, ctx, http.MethodPost, "/channel/escalate/rule/info", req)
	return result, err
}

// UpdateEscalateRule updates an escalate rule.
func (c *Client) UpdateEscalateRule(ctx context.Context, req *UpdateEscalateRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/escalate/rule/update", req)
	return err
}

// DeleteEscalateRule deletes an escalate rule.
func (c *Client) DeleteEscalateRule(ctx context.Context, req *DeleteEscalateRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/escalate/rule/delete", req)
	return err
}

// EnableEscalateRule enables an escalate rule.
func (c *Client) EnableEscalateRule(ctx context.Context, req *EnableEscalateRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/escalate/rule/enable", req)
	return err
}

// DisableEscalateRule disables an escalate rule.
func (c *Client) DisableEscalateRule(ctx context.Context, req *EnableEscalateRuleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/escalate/rule/disable", req)
	return err
}
