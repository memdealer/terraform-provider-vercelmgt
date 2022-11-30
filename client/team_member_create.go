package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type CreateTeamMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type createTeamMemberResponse struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (c *Client) CreateTeamMember(ctx context.Context, request CreateTeamMemberRequest) (r TeamMemberResponse, err error) {
	url := fmt.Sprintf("%s/v1/teams/%s/members", c.baseURL, c._teamID)

	payload := string(mustMarshal(request))
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url,
		strings.NewReader(payload),
	)
	if err != nil {
		return r, err
	}

	tflog.Trace(ctx, "creating team member", map[string]interface{}{
		"url":     url,
		"payload": payload,
	})
	var teamMemberResponse createTeamMemberResponse
	err = c.doRequest(req, &teamMemberResponse)
	if err != nil {
		return r, err
	}

	return TeamMemberResponse{
		UID:   teamMemberResponse.UID,
		Email: teamMemberResponse.Email,
		Role:  teamMemberResponse.Role,
	}, nil

}
