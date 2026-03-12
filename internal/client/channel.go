package client

import (
	"context"
	"net/http"
)

// ChannelGroup configures how alerts are aggregated into incidents.
type ChannelGroup struct {
	Method            string          `json:"method"`
	Cases             []GroupCaseRule `json:"cases,omitempty"`
	Equals            [][]string      `json:"equals,omitempty"`
	AllEqualsRequired bool            `json:"all_equals_required,omitempty"`
	TimeWindow        int             `json:"time_window,omitempty"`
	IKeys             []string        `json:"i_keys,omitempty"`
	IScoreThreshold   float64         `json:"i_score_threshold,omitempty"`
	StormThresholds   []int           `json:"storm_thresholds,omitempty"`
}

// GroupCaseRule is a branch-level grouping rule with filter conditions and merge dimensions.
type GroupCaseRule struct {
	If     []Filter `json:"if"`
	Equals []string `json:"equals"`
}

// ChannelFlapping configures incident flap detection to reduce noise from rapidly changing alerts.
type ChannelFlapping struct {
	IsDisabled bool `json:"is_disabled"`
	MaxChanges int  `json:"max_changes,omitempty"`
	InMins     int  `json:"in_mins,omitempty"`
	MuteMins   int  `json:"mute_mins,omitempty"`
}

// Channel represents a Flashduty collaboration space (channel).
type Channel struct {
	ChannelID               int64            `json:"channel_id"`
	ChannelName             string           `json:"channel_name"`
	Description             string           `json:"description,omitempty"`
	TeamID                  int64            `json:"team_id,omitempty"`
	ManagingTeamIDs         []int64          `json:"managing_team_ids,omitempty"`
	AutoResolveTimeout      int              `json:"auto_resolve_timeout,omitempty"`
	AutoResolveMode         string           `json:"auto_resolve_mode,omitempty"`
	IsPrivate               bool             `json:"is_private,omitempty"`
	DisableOutlierDetection bool             `json:"disable_outlier_detection,omitempty"`
	DisableAutoClose        bool             `json:"disable_auto_close,omitempty"`
	Group                   *ChannelGroup    `json:"group,omitempty"`
	Flapping                *ChannelFlapping `json:"flapping,omitempty"`
	Status                  string           `json:"status,omitempty"`
	CreatedAt               int64            `json:"created_at,omitempty"`
	UpdatedAt               int64            `json:"updated_at,omitempty"`
}

// CreateChannelRequest represents the request body for creating a channel.
type CreateChannelRequest struct {
	ChannelName             string           `json:"channel_name"`
	Description             string           `json:"description,omitempty"`
	TeamID                  int64            `json:"team_id"`
	ManagingTeamIDs         []int64          `json:"managing_team_ids,omitempty"`
	AutoResolveTimeout      int              `json:"auto_resolve_timeout,omitempty"`
	AutoResolveMode         string           `json:"auto_resolve_mode,omitempty"`
	IsPrivate               *bool            `json:"is_private,omitempty"`
	DisableOutlierDetection *bool            `json:"disable_outlier_detection,omitempty"`
	DisableAutoClose        *bool            `json:"disable_auto_close,omitempty"`
	Group                   *ChannelGroup    `json:"group,omitempty"`
	Flapping                *ChannelFlapping `json:"flapping,omitempty"`
}

// CreateChannelResult represents the response data for channel creation.
type CreateChannelResult struct {
	ChannelID   int64  `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

// UpdateChannelRequest represents the request body for updating a channel.
type UpdateChannelRequest struct {
	ChannelID               int64            `json:"channel_id"`
	ChannelName             string           `json:"channel_name,omitempty"`
	Description             string           `json:"description,omitempty"`
	TeamID                  int64            `json:"team_id,omitempty"`
	ManagingTeamIDs         []int64          `json:"managing_team_ids,omitempty"`
	AutoResolveTimeout      *int             `json:"auto_resolve_timeout,omitempty"`
	AutoResolveMode         string           `json:"auto_resolve_mode,omitempty"`
	IsPrivate               *bool            `json:"is_private,omitempty"`
	DisableOutlierDetection *bool            `json:"disable_outlier_detection,omitempty"`
	DisableAutoClose        *bool            `json:"disable_auto_close,omitempty"`
	Group                   *ChannelGroup    `json:"group,omitempty"`
	Flapping                *ChannelFlapping `json:"flapping,omitempty"`
}

// GetChannelRequest represents the request body for getting channel info.
type GetChannelRequest struct {
	ChannelID int64 `json:"channel_id"`
}

// DeleteChannelRequest represents the request body for deleting a channel.
type DeleteChannelRequest struct {
	ChannelID int64 `json:"channel_id"`
}

// ChannelStatusRequest represents the request body for enabling or disabling a channel.
type ChannelStatusRequest struct {
	ChannelID int64 `json:"channel_id"`
}

// CreateChannel creates a new channel.
func (c *Client) CreateChannel(ctx context.Context, req *CreateChannelRequest) (*CreateChannelResult, error) {
	result, _, err := doRequestWithResponse[CreateChannelResult](c, ctx, http.MethodPost, "/channel/create", req)
	return result, err
}

// GetChannel retrieves channel information.
func (c *Client) GetChannel(ctx context.Context, req *GetChannelRequest) (*Channel, error) {
	result, _, err := doRequestWithResponse[Channel](c, ctx, http.MethodPost, "/channel/info", req)
	return result, err
}

// UpdateChannel updates a channel.
func (c *Client) UpdateChannel(ctx context.Context, req *UpdateChannelRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/update", req)
	return err
}

// DeleteChannel deletes a channel.
func (c *Client) DeleteChannel(ctx context.Context, req *DeleteChannelRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/delete", req)
	return err
}

// EnableChannel enables a channel.
func (c *Client) EnableChannel(ctx context.Context, req *ChannelStatusRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/enable", req)
	return err
}

// DisableChannel disables a channel.
func (c *Client) DisableChannel(ctx context.Context, req *ChannelStatusRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/channel/disable", req)
	return err
}

// ListChannelsRequest represents the request body for listing channels.
type ListChannelsRequest struct {
	Page    int     `json:"p,omitempty"`
	Limit   int     `json:"limit,omitempty"`
	Query   string  `json:"query,omitempty"`
	TeamIDs []int64 `json:"team_ids,omitempty"`
	IsBrief bool    `json:"is_brief,omitempty"`
}

// ListChannelsResult represents the response data for channel list.
type ListChannelsResult struct {
	Items       []Channel `json:"items"`
	HasNextPage bool      `json:"has_next_page"`
	Total       int       `json:"total"`
}

// ListChannels retrieves a list of channels.
func (c *Client) ListChannels(ctx context.Context, req *ListChannelsRequest) (*ListChannelsResult, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 100
	}
	result, _, err := doRequestWithResponse[ListChannelsResult](c, ctx, http.MethodPost, "/channel/list", req)
	return result, err
}
