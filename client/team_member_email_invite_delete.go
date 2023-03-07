package client

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
)

func (c *Client) DeleteTeamMemberInviteCode(ctx context.Context, inviteID string) (err error) {

	url := fmt.Sprintf("%s/v1/teams/%s/invites/%s", c.baseURL, c._teamID, inviteID)

	tflog.Info(ctx, "[DELETING TEAM MEMBER INVITE CODE]", map[string]interface{}{
		"url":      url,
		"inviteID": inviteID,
	})

	req, err := http.NewRequest(
		"DELETE",
		url,
		nil,
	)
	if err != nil {
		return err
	}

	err = c.doRequest(req, nil)

	tflog.Trace(ctx, "[DELETED TEAM MEMBER INVITE CODE]", map[string]interface{}{
		"url":      url,
		"inviteID": inviteID,
	})
	return err
}
