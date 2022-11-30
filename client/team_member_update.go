package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type UpdateTeamMemberRequest struct {
	Role *string `json:"role"`
}

type UpdateTeamMemberResponse struct {
	TeamID *string `json:"id"`
}

func (c *Client) UpdateTeamMember(ctx context.Context, memberUID string, request UpdateTeamMemberRequest) (r UpdateTeamMemberResponse, err error) {
	url := fmt.Sprintf("%s/v1/teams/%s/members/%s", c.baseURL, c._teamID, memberUID)

	tflog.Trace(ctx, "updating team member", map[string]interface{}{
		"url": url,
	})
	payload := string(mustMarshal(request))
	req, err := http.NewRequestWithContext(
		ctx,
		"PATCH",
		url,
		strings.NewReader(payload),
	)
	if err != nil {
		return r, err
	}

	tflog.Trace(ctx, "updating team member", map[string]interface{}{
		"url":     url,
		"payload": payload,
	})
	err = c.doRequest(req, &r)
	return r, err
}
