package client

import (
	"context"
	"net/http"
	"time"
)

// ScheduleLayer represents a schedule layer configuration.
type ScheduleLayer struct {
	LayerName             string           `json:"layer_name"`
	Mode                  int              `json:"mode"`
	Groups                []ScheduleGroup  `json:"groups"`
	FairRotation          bool             `json:"fair_rotation"`
	HandoffTime           int              `json:"handoff_time"`
	LayerStart            int64            `json:"layer_start"`
	LayerEnd              int64            `json:"layer_end,omitempty"`
	RestrictMode          int              `json:"restrict_mode"`
	RestrictPeriods       []RestrictPeriod `json:"restrict_periods,omitempty"`
	DayMask               *DayMask         `json:"day_mask,omitempty"`
	MaskContinuousEnabled bool             `json:"mask_continuous_enabled"`
	RotationUnit          string           `json:"rotation_unit"`
	RotationValue         int              `json:"rotation_value"`
}

// ScheduleGroup represents a group within a schedule layer.
type ScheduleGroup struct {
	GroupName string                `json:"group_name"`
	Members   []ScheduleGroupMember `json:"members"`
}

// ScheduleGroupMember represents a member within a schedule group.
type ScheduleGroupMember struct {
	RoleID    int64   `json:"role_id"`
	PersonIDs []int64 `json:"person_ids"`
}

// RestrictPeriod represents a time restriction period.
type RestrictPeriod struct {
	RestrictStart int `json:"restrict_start"`
	RestrictEnd   int `json:"restrict_end"`
}

// DayMask represents day of week mask configuration.
type DayMask struct {
	Repeat []int `json:"repeat"`
}

// ScheduleNotify represents notification settings for a schedule.
type ScheduleNotify struct {
	AdvanceInTime int              `json:"advance_in_time,omitempty"`
	FixedTime     *NotifyFixedTime `json:"fixed_time,omitempty"`
	By            *NotifyBy        `json:"by,omitempty"`
	Webhooks      []NotifyWebhook  `json:"webhooks,omitempty"`
}

// NotifyFixedTime represents fixed time notification settings.
type NotifyFixedTime struct {
	Cycle string `json:"cycle,omitempty"`
	Start string `json:"start,omitempty"`
}

// NotifyBy represents personal notification settings.
type NotifyBy struct {
	FollowPreference bool     `json:"follow_preference,omitempty"`
	PersonalChannels []string `json:"personal_channels,omitempty"`
}

// NotifyWebhook represents webhook notification settings.
type NotifyWebhook struct {
	Type     string         `json:"type,omitempty"`
	Settings map[string]any `json:"settings,omitempty"`
}

// Schedule represents a complete schedule object.
type Schedule struct {
	ScheduleID   int64           `json:"schedule_id"`
	ScheduleName string          `json:"schedule_name"`
	Description  string          `json:"description,omitempty"`
	Status       int             `json:"status,omitempty"`
	TeamID       int64           `json:"team_id,omitempty"`
	Layers       []ScheduleLayer `json:"layers,omitempty"`
	Notify       *ScheduleNotify `json:"notify,omitempty"`
	CreatedAt    int64           `json:"created_at,omitempty"`
	UpdatedAt    int64           `json:"updated_at,omitempty"`
}

// CreateScheduleRequest represents the request body for creating a schedule.
type CreateScheduleRequest struct {
	ScheduleName string          `json:"schedule_name"`
	Description  string          `json:"description"`
	Status       int             `json:"status"`
	TeamID       int64           `json:"team_id,omitempty"`
	Layers       []ScheduleLayer `json:"layers"`
	Notify       *ScheduleNotify `json:"notify,omitempty"`
}

// CreateScheduleResult represents the response data for schedule creation.
type CreateScheduleResult struct {
	ScheduleID int64 `json:"schedule_id"`
}

// UpdateScheduleRequest represents the request body for updating a schedule.
type UpdateScheduleRequest struct {
	ScheduleID   int64           `json:"schedule_id"`
	ScheduleName string          `json:"schedule_name,omitempty"`
	Description  string          `json:"description,omitempty"`
	Status       *int            `json:"status,omitempty"`
	TeamID       int64           `json:"team_id,omitempty"`
	Layers       []ScheduleLayer `json:"layers,omitempty"`
	Notify       *ScheduleNotify `json:"notify,omitempty"`
}

// GetScheduleRequest represents the request body for getting schedule info.
type GetScheduleRequest struct {
	ScheduleID int64 `json:"schedule_id"`
	Start      int64 `json:"start"`
	End        int64 `json:"end"`
}

// DeleteScheduleRequest represents the request body for deleting schedules.
type DeleteScheduleRequest struct {
	ScheduleIDs []int64 `json:"schedule_ids"`
}

// CreateSchedule creates a new schedule.
func (c *Client) CreateSchedule(ctx context.Context, req *CreateScheduleRequest) (*CreateScheduleResult, error) {
	result, _, err := doRequestWithResponse[CreateScheduleResult](c, ctx, http.MethodPost, "/schedule/create", req)
	return result, err
}

// GetSchedule retrieves schedule information.
func (c *Client) GetSchedule(ctx context.Context, req *GetScheduleRequest) (*Schedule, error) {
	localReq := *req
	now := time.Now()
	if localReq.Start == 0 {
		localReq.Start = now.Unix()
	}
	if localReq.End == 0 {
		localReq.End = now.Add(7 * 24 * time.Hour).Unix()
	}

	result, _, err := doRequestWithResponse[Schedule](c, ctx, http.MethodPost, "/schedule/info", &localReq)
	return result, err
}

// UpdateSchedule updates a schedule.
func (c *Client) UpdateSchedule(ctx context.Context, req *UpdateScheduleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/schedule/update", req)
	return err
}

// DeleteSchedule deletes schedules by IDs.
func (c *Client) DeleteSchedule(ctx context.Context, req *DeleteScheduleRequest) error {
	_, _, err := doRequestWithResponse[any](c, ctx, http.MethodPost, "/schedule/delete", req)
	return err
}
