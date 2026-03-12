package client

import (
	"context"
	"fmt"
	"net/http"
)

// MemberInvite represents the request body for inviting a member.
type MemberInvite struct {
	Email       string  `json:"email,omitempty"`
	Phone       string  `json:"phone,omitempty"`
	CountryCode string  `json:"country_code,omitempty"`
	MemberName  string  `json:"member_name,omitempty"`
	RefID       string  `json:"ref_id,omitempty"`
	RoleIDs     []int64 `json:"role_ids,omitempty"`
}

// MemberInviteResult represents the response data for member invitation.
type MemberInviteResult struct {
	Items []MemberShort `json:"items"`
}

// MemberShort represents a brief member info returned from invite API.
type MemberShort struct {
	MemberID   int64  `json:"member_id"`
	MemberName string `json:"member_name"`
}

// Member represents a full member object.
type Member struct {
	MemberID       int64   `json:"member_id"`
	MemberName     string  `json:"member_name"`
	Email          string  `json:"email,omitempty"`
	Phone          string  `json:"phone,omitempty"`
	CountryCode    string  `json:"country_code,omitempty"`
	RefID          string  `json:"ref_id,omitempty"`
	Status         string  `json:"status,omitempty"`
	AccountRoleIDs []int64 `json:"account_role_ids,omitempty"`
	PhoneVerified  bool    `json:"phone_verified,omitempty"`
	EmailVerified  bool    `json:"email_verified,omitempty"`
	CreatedAt      int64   `json:"created_at,omitempty"`
	UpdatedAt      int64   `json:"updated_at,omitempty"`
}

// MemberListResult represents the response data for member list.
type MemberListResult struct {
	Page  int      `json:"p"`
	Limit int      `json:"limit"`
	Total int      `json:"total"`
	Items []Member `json:"items"`
}

// GetMemberRequest represents the request body for getting a member by various identifiers.
type GetMemberRequest struct {
	MemberID   int64  `json:"member_id,omitempty"`
	MemberName string `json:"member_name,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	RefID      string `json:"ref_id,omitempty"`
}

// UpdateMemberRequest represents the request body for updating member info.
type UpdateMemberRequest struct {
	MemberID   int64               `json:"member_id,omitempty"`
	MemberName string              `json:"member_name,omitempty"`
	Email      string              `json:"email,omitempty"`
	Phone      string              `json:"phone,omitempty"`
	RefID      string              `json:"ref_id,omitempty"`
	Updates    MemberUpdatePayload `json:"updates"`
}

// MemberUpdatePayload represents the fields that can be updated.
type MemberUpdatePayload struct {
	Phone       string `json:"phone,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	Email       string `json:"email,omitempty"`
	MemberName  string `json:"member_name,omitempty"`
	TimeZone    string `json:"time_zone,omitempty"`
	Locale      string `json:"locale,omitempty"`
	RefID       string `json:"ref_id,omitempty"`
}

// DeleteMemberRequest represents the request body for deleting a member.
type DeleteMemberRequest struct {
	MemberID int64  `json:"member_id,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	RefID    string `json:"ref_id,omitempty"`
}

// InviteMembers invites one or more members to Flashduty.
// The "from": "api" parameter bypasses member verification (requires whitelist).
func (c *Client) InviteMembers(ctx context.Context, members []MemberInvite) (*MemberInviteResult, error) {
	reqBody := map[string]any{
		"members": members,
		"from":    "api",
	}

	result, _, err := doRequestWithResponse[MemberInviteResult](c, ctx, http.MethodPost, "/member/invite", reqBody)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ListMembers retrieves a list of members.
func (c *Client) ListMembers(ctx context.Context, page, limit int, query string) (*MemberListResult, error) {
	reqBody := map[string]any{
		"p":     page,
		"limit": limit,
	}
	if query != "" {
		reqBody["query"] = query
	}

	result, _, err := doRequestWithResponse[MemberListResult](c, ctx, http.MethodPost, "/member/list", reqBody)
	return result, err
}

// GetMemberByID retrieves a member by querying all pages of the member list.
func (c *Client) GetMemberByID(ctx context.Context, memberID int64) (*Member, error) {
	page := 1
	limit := 100
	for {
		result, err := c.ListMembers(ctx, page, limit, "")
		if err != nil {
			return nil, err
		}

		for i := range result.Items {
			if result.Items[i].MemberID == memberID {
				return &result.Items[i], nil
			}
		}

		if len(result.Items) < limit {
			break
		}
		page++
	}

	return nil, fmt.Errorf("%w: member_id=%d", ErrNotFound, memberID)
}

// UpdateMember updates member information.
func (c *Client) UpdateMember(ctx context.Context, req *UpdateMemberRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/member/info/reset", req)
	return err
}

// DeleteMember deletes a member.
func (c *Client) DeleteMember(ctx context.Context, req *DeleteMemberRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/member/delete", req)
	return err
}
