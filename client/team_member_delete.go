package client

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
)

func (c *Client) DeleteTeamMember(ctx context.Context, memberUID string) (err error) {

	url := fmt.Sprintf("%s/v1/teams/%s/members/%s", c.baseURL, c._teamID, memberUID)

	tflog.Trace(ctx, "deleting team member", map[string]interface{}{
		"url": url,
	})

	req, err := http.NewRequest(
		"DELETE",
		url,
		nil,
	)
	if err != nil {
		return err
	}

	tflog.Trace(ctx, "deleting team member", map[string]interface{}{
		"url": url,
	})
	err = c.doRequest(req, nil)

	return err
}
