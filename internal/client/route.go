package client

import (
	"context"
	"net/http"
)

// RouteFilter represents a filter condition in route rules.
type RouteFilter struct {
	Key  string   `json:"key"`
	Oper string   `json:"oper"`
	Vals []string `json:"vals"`
}

// RouteCase represents a conditional routing rule.
type RouteCase struct {
	If               []RouteFilter `json:"if"`
	ChannelIDs       []int64       `json:"channel_ids,omitempty"`
	Fallthrough      bool          `json:"fallthrough,omitempty"`
	RoutingMode      string        `json:"routing_mode,omitempty"`
	NameMappingLabel string        `json:"name_mapping_label,omitempty"`
}

// RouteDefault represents default routing configuration.
type RouteDefault struct {
	ChannelIDs []int64 `json:"channel_ids"`
}

// RouteSection represents a section divider in route rules.
type RouteSection struct {
	Name     string `json:"name,omitempty"`
	Position int    `json:"position,omitempty"`
}

// Route represents routing configuration for an integration.
type Route struct {
	IntegrationID int64          `json:"integration_id,omitempty"`
	Cases         []RouteCase    `json:"cases,omitempty"`
	Sections      []RouteSection `json:"sections,omitempty"`
	Default       *RouteDefault  `json:"default,omitempty"`
	Status        string         `json:"status,omitempty"`
	Version       int            `json:"version,omitempty"`
	CreatorID     int64          `json:"creator_id,omitempty"`
	UpdatedBy     int64          `json:"updated_by,omitempty"`
	CreatedAt     int64          `json:"created_at,omitempty"`
	UpdatedAt     int64          `json:"updated_at,omitempty"`
}

// RouteHistory represents a historical version of route configuration.
type RouteHistory struct {
	Cases     []RouteCase    `json:"cases,omitempty"`
	Sections  []RouteSection `json:"sections,omitempty"`
	Default   *RouteDefault  `json:"default,omitempty"`
	Version   int            `json:"version,omitempty"`
	UpdatedBy int64          `json:"updated_by,omitempty"`
	UpdatedAt int64          `json:"updated_at,omitempty"`
}

// ListRoutesRequest represents a request to list routes.
type ListRoutesRequest struct {
	IntegrationIDs []int64 `json:"integration_ids"`
}

// ListRoutesResult represents the result of listing routes.
type ListRoutesResult struct {
	Items []Route `json:"items"`
}

// GetRouteRequest represents a request to get route info.
type GetRouteRequest struct {
	IntegrationID int64 `json:"integration_id"`
}

// UpsertRouteRequest represents a request to create or update a route.
type UpsertRouteRequest struct {
	IntegrationID int64          `json:"integration_id"`
	Version       int            `json:"version"`
	Cases         []RouteCase    `json:"cases,omitempty"`
	Sections      []RouteSection `json:"sections,omitempty"`
	Default       *RouteDefault  `json:"default,omitempty"`
}

// ListRouteHistoryRequest represents a request to list route history.
type ListRouteHistoryRequest struct {
	IntegrationID int64  `json:"integration_id"`
	Page          int    `json:"p,omitempty"`
	Limit         int    `json:"limit,omitempty"`
	Asc           bool   `json:"asc,omitempty"`
	OrderBy       string `json:"orderby,omitempty"`
}

// ListRouteHistoryResult represents the result of listing route history.
type ListRouteHistoryResult struct {
	Items []RouteHistory `json:"items"`
}

// RoutePreviewRequest represents a request to preview route matching.
type RoutePreviewRequest struct {
	AlertIDs []string `json:"alert_ids"`
	Route    struct {
		Cases   []RouteCase   `json:"cases,omitempty"`
		Default *RouteDefault `json:"default,omitempty"`
	} `json:"route"`
}

// RoutePreviewItem represents a single preview result.
type RoutePreviewItem struct {
	AlertID      string   `json:"alert_id"`
	AlertTitle   string   `json:"alert_title"`
	StartTime    int64    `json:"start_time"`
	MatchedCases []string `json:"matched_cases"`
}

// RoutePreviewResult represents the result of route preview.
type RoutePreviewResult struct {
	Items []RoutePreviewItem `json:"items"`
}

// ListRoutes retrieves routes for the given integration IDs.
func (c *Client) ListRoutes(ctx context.Context, req *ListRoutesRequest) (*ListRoutesResult, error) {
	result, _, err := doRequestWithResponse[ListRoutesResult](c, ctx, http.MethodPost, "/route/list", req)
	return result, err
}

// GetRoute retrieves route configuration for an integration.
func (c *Client) GetRoute(ctx context.Context, req *GetRouteRequest) (*Route, error) {
	result, _, err := doRequestWithResponse[Route](c, ctx, http.MethodPost, "/route/info", req)
	return result, err
}

// UpsertRoute creates or updates route configuration.
func (c *Client) UpsertRoute(ctx context.Context, req *UpsertRouteRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/route/upsert", req)
	return err
}

// ListRouteHistory retrieves route history for an integration.
func (c *Client) ListRouteHistory(ctx context.Context, req *ListRouteHistoryRequest) (*ListRouteHistoryResult, error) {
	result, _, err := doRequestWithResponse[ListRouteHistoryResult](c, ctx, http.MethodPost, "/route/history/list", req)
	return result, err
}

// PreviewRoute previews route matching for given alerts.
func (c *Client) PreviewRoute(ctx context.Context, req *RoutePreviewRequest) (*RoutePreviewResult, error) {
	result, _, err := doRequestWithResponse[RoutePreviewResult](c, ctx, http.MethodPost, "/route/preview", req)
	return result, err
}
