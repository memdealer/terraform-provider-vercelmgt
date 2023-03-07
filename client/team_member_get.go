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

type EmailInvitationCodesResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type TeamMemberResponseFlattened struct {
	Members []TeamMemberResponse `json:"members"`
}

type InvitationCodesFlattened struct {
	InvitationCodes []EmailInvitationCodesResponse `json:"emailInviteCodes"`
}

func (c *Client) GetTeamMemberInviteCode(ctx context.Context, memberEmail string) (r EmailInvitationCodesResponse, err error) {
	url := fmt.Sprintf("%s/v2/teams/%s/members?search=%s&limit=1", c.baseURL, c._teamID, memberEmail)
	var flatResp InvitationCodesFlattened

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		url,
		nil,
	)

	if err != nil {
		return r, fmt.Errorf("creating request: %s", err)
	}
	tflog.Trace(ctx, "getting team member email codes", map[string]interface{}{
		"url":   url,
		"email": memberEmail,
	})

	err = c.doRequest(req, &flatResp)

	tflog.Trace(ctx, "Get codes response", map[string]interface{}{
		"InviteCodes:": flatResp.InvitationCodes,
	})

	if err != nil {
		tflog.Error(ctx, "Could not retrieve email codes invitation.", map[string]interface{}{
			"url":   url,
			"email": memberEmail,
		})
		return EmailInvitationCodesResponse{}, err
	}
	// iterate over email codes, find the needed one where email == matches with what we are finding.
	for _, v := range flatResp.InvitationCodes {
		tflog.Info(ctx, "TeamMember response", map[string]interface{}{
			"code":  v.ID,
			"email": v.Email,
		})
		if v.Email == memberEmail {
			return EmailInvitationCodesResponse{
				ID:    v.ID,
				Email: v.Email,
				Role:  v.Role,
			}, nil
		}
	}
	// if nothing is found, then something wrong and gotta shut down
	return EmailInvitationCodesResponse{}, nil
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
	if len(flatResp.Members) > 0 {
		return TeamMemberResponse{
			UID:   flatResp.Members[0].UID,
			Email: flatResp.Members[0].Email,
			Role:  flatResp.Members[0].Role,
		}, err
	}
	return TeamMemberResponse{}, err
}
