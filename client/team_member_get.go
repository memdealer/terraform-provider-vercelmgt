package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TeamMemberResponse struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type TeamMemberResponseFlattened struct {
	Members []TeamMemberResponse `json:"members"`
}

func (c *Client) GetTeamMember(ctx context.Context, memberEmail string) (r TeamMemberResponse, err error) {
	url := fmt.Sprintf("%s/v2/teams/%s/members?search=%s&limit=1", c.baseURL, c._teamID, memberEmail)
	var flatResp TeamMemberResponseFlattened

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		url,
		nil,
	)
	if err != nil {
		return r, fmt.Errorf("creating request: %s", err)
	}
	tflog.Trace(ctx, "getting team member", map[string]interface{}{
		"url":   url,
		"email": memberEmail,
	})

	err = c.doRequest(req, &flatResp)

	tflog.Trace(ctx, "TeamMember response", map[string]interface{}{
		"ResponseMarker": flatResp.Members,
	})

	// During import, if none found -> none should be returned
	// The check prevents panic.
	if len(flatResp.Members) > 1 {
		return TeamMemberResponse{
			UID:   flatResp.Members[0].UID,
			Email: flatResp.Members[0].Email,
			Role:  flatResp.Members[0].Role,
		}, err
	}
	return TeamMemberResponse{}, err
}
