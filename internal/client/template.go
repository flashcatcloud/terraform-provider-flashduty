package client

import (
	"context"
	"net/http"
)

// Template represents a FlashDuty notification template.
type Template struct {
	TemplateID   string `json:"template_id"`
	TeamID       int64  `json:"team_id,omitempty"`
	TemplateName string `json:"template_name"`
	Description  string `json:"description,omitempty"`
	Email        string `json:"email,omitempty"`
	SMS          string `json:"sms,omitempty"`
	Dingtalk     string `json:"dingtalk,omitempty"`
	Wecom        string `json:"wecom,omitempty"`
	Feishu       string `json:"feishu,omitempty"`
	FeishuApp    string `json:"feishu_app,omitempty"`
	DingtalkApp  string `json:"dingtalk_app,omitempty"`
	WecomApp     string `json:"wecom_app,omitempty"`
	TeamsApp     string `json:"teams_app,omitempty"`
	SlackApp     string `json:"slack_app,omitempty"`
	Slack        string `json:"slack,omitempty"`
	Zoom         string `json:"zoom,omitempty"`
	Telegram     string `json:"telegram,omitempty"`
	Status       string `json:"status,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
	UpdatedAt    int64  `json:"updated_at,omitempty"`
}

type CreateTemplateRequest struct {
	TeamID       int64  `json:"team_id,omitempty"`
	TemplateName string `json:"template_name"`
	Description  string `json:"description,omitempty"`
	Email        string `json:"email,omitempty"`
	SMS          string `json:"sms,omitempty"`
	Dingtalk     string `json:"dingtalk,omitempty"`
	Wecom        string `json:"wecom,omitempty"`
	Feishu       string `json:"feishu,omitempty"`
	FeishuApp    string `json:"feishu_app,omitempty"`
	DingtalkApp  string `json:"dingtalk_app,omitempty"`
	WecomApp     string `json:"wecom_app,omitempty"`
	TeamsApp     string `json:"teams_app,omitempty"`
	SlackApp     string `json:"slack_app,omitempty"`
	Slack        string `json:"slack,omitempty"`
	Zoom         string `json:"zoom,omitempty"`
	Telegram     string `json:"telegram,omitempty"`
}

type CreateTemplateResult struct {
	TemplateID   string `json:"template_id"`
	TemplateName string `json:"template_name"`
}

type UpdateTemplateRequest struct {
	TemplateID   string `json:"template_id"`
	TeamID       int64  `json:"team_id,omitempty"`
	TemplateName string `json:"template_name"`
	Description  string `json:"description,omitempty"`
	Email        string `json:"email,omitempty"`
	SMS          string `json:"sms,omitempty"`
	Dingtalk     string `json:"dingtalk,omitempty"`
	Wecom        string `json:"wecom,omitempty"`
	Feishu       string `json:"feishu,omitempty"`
	FeishuApp    string `json:"feishu_app,omitempty"`
	DingtalkApp  string `json:"dingtalk_app,omitempty"`
	WecomApp     string `json:"wecom_app,omitempty"`
	TeamsApp     string `json:"teams_app,omitempty"`
	SlackApp     string `json:"slack_app,omitempty"`
	Slack        string `json:"slack,omitempty"`
	Zoom         string `json:"zoom,omitempty"`
	Telegram     string `json:"telegram,omitempty"`
}

type GetTemplateRequest struct {
	TemplateID string `json:"template_id"`
}

type DeleteTemplateRequest struct {
	TemplateID string `json:"template_id"`
}

type ListTemplatesRequest struct {
	Page    int     `json:"p,omitempty"`
	Limit   int     `json:"limit,omitempty"`
	Query   string  `json:"query,omitempty"`
	TeamIDs []int64 `json:"team_ids,omitempty"`
}

type ListTemplatesResult struct {
	Items       []Template `json:"items"`
	HasNextPage bool       `json:"has_next_page"`
}

func (c *Client) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*CreateTemplateResult, error) {
	result, _, err := doRequestWithResponse[CreateTemplateResult](c, ctx, http.MethodPost, "/template/create", req)
	return result, err
}

func (c *Client) GetTemplate(ctx context.Context, req *GetTemplateRequest) (*Template, error) {
	result, _, err := doRequestWithResponse[Template](c, ctx, http.MethodPost, "/template/info", req)
	return result, err
}

func (c *Client) UpdateTemplate(ctx context.Context, req *UpdateTemplateRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/template/update", req)
	return err
}

func (c *Client) DeleteTemplate(ctx context.Context, req *DeleteTemplateRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/template/delete", req)
	return err
}

func (c *Client) ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResult, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 100
	}
	result, _, err := doRequestWithResponse[ListTemplatesResult](c, ctx, http.MethodPost, "/template/list", req)
	return result, err
}
