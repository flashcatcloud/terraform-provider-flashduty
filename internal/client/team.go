package client

import (
	"context"
	"net/http"
)

// Team represents a Flashduty team.
type Team struct {
	TeamID        int64   `json:"team_id"`
	TeamName      string  `json:"team_name"`
	Description   string  `json:"description,omitempty"`
	RefID         string  `json:"ref_id,omitempty"`
	PersonIDs     []int64 `json:"person_ids,omitempty"`
	CreatorID     int64   `json:"creator_id,omitempty"`
	CreatedAt     int64   `json:"created_at,omitempty"`
	UpdatedAt     int64   `json:"updated_at,omitempty"`
	UpdatedBy     int64   `json:"updated_by,omitempty"`
	UpdatedByName string  `json:"updated_by_name,omitempty"`
}

// UpsertTeamRequest represents the request body for creating or updating a team.
type UpsertTeamRequest struct {
	TeamID           int64    `json:"team_id,omitempty"`
	TeamName         string   `json:"team_name"`
	Description      string   `json:"description"`
	RefID            string   `json:"ref_id,omitempty"`
	PersonIDs        []int64  `json:"person_ids,omitempty"`
	Emails           []string `json:"emails,omitempty"`
	Phones           []string `json:"phones,omitempty"`
	ResetIfNameExist bool     `json:"reset_if_name_exist,omitempty"`
}

// UpsertTeamResult represents the response data for team upsert.
type UpsertTeamResult struct {
	TeamID int64 `json:"team_id"`
}

// GetTeamRequest represents the request body for getting team info.
type GetTeamRequest struct {
	TeamID   int64  `json:"team_id,omitempty"`
	TeamName string `json:"team_name,omitempty"`
	RefID    string `json:"ref_id,omitempty"`
}

// DeleteTeamRequest represents the request body for deleting a team.
type DeleteTeamRequest struct {
	TeamID   int64  `json:"team_id,omitempty"`
	TeamName string `json:"team_name,omitempty"`
	RefID    string `json:"ref_id,omitempty"`
}

// UpsertTeam creates or updates a team.
func (c *Client) UpsertTeam(ctx context.Context, req *UpsertTeamRequest) (*UpsertTeamResult, error) {
	result, _, err := doRequestWithResponse[UpsertTeamResult](c, ctx, http.MethodPost, "/team/upsert", req)
	return result, err
}

// GetTeam retrieves team information.
func (c *Client) GetTeam(ctx context.Context, req *GetTeamRequest) (*Team, error) {
	result, _, err := doRequestWithResponse[Team](c, ctx, http.MethodPost, "/team/info", req)
	return result, err
}

// DeleteTeam deletes a team.
func (c *Client) DeleteTeam(ctx context.Context, req *DeleteTeamRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/team/delete", req)
	return err
}

// ListTeamsRequest represents the request body for listing teams.
type ListTeamsRequest struct {
	Page     int    `json:"p,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Query    string `json:"query,omitempty"`
	PersonID int64  `json:"person_id,omitempty"`
}

// ListTeamsResult represents the response data for team list.
type ListTeamsResult struct {
	Items []Team `json:"items"`
	Page  int    `json:"p"`
	Limit int    `json:"limit"`
	Total int    `json:"total"`
}

// ListTeams retrieves a list of teams.
func (c *Client) ListTeams(ctx context.Context, req *ListTeamsRequest) (*ListTeamsResult, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 100
	}
	result, _, err := doRequestWithResponse[ListTeamsResult](c, ctx, http.MethodPost, "/team/list", req)
	return result, err
}
