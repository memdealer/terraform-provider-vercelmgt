package client

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
)

type ListTeamMembersResponse struct {
	Members    []TeamMemberResponse `json:"members"`
	Pagination struct {
		Count int `json:"count"`
		Next  int `json:"next"`
		Prev  int `json:"prev"`
	} `json:"pagination"`
}

func (c *Client) GetTeamMembers(ctx context.Context) (r TeamMemberResponse, err error) {
	url := fmt.Sprintf("%s/v2/teams/%s/members", c.baseURL, c._teamID)
	//var ListOfMembers []TeamMemberResponse

	tflog.Trace(ctx, "getting team members list", map[string]interface{}{
		"url": url,
	})

	for true {
		var flatResp ListTeamMembersResponse

		req, err := http.NewRequestWithContext(
			ctx,
			"GET",
			url,
			nil,
		)
		if err != nil {
			return r, fmt.Errorf("creating request: %s", err)
		}

		err = c.doRequest(req, &flatResp)

		tflog.Trace(ctx, "TeamMember response", map[string]interface{}{
			"ResponseMarker": flatResp.Members,
			"Paginator":      flatResp.Pagination.Next,
		})

		for _, v := range flatResp.Members {
			tflog.Trace(ctx, "TeamMember response", map[string]interface{}{
				"Member": v,
			})

		}

	}

	return TeamMemberResponse{}, err
}
