package client

import (
	"context"
	"errors"
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
	// Sometimes, the invited user has no account and vercel logic is plain simple is to not return UUID if user has not been created.
	// I must not argue that such a logic is faulty, as then you should provide either a method to see if user exists on vercel, or at least do not add it to /getTeamMembers call.
	// Anyhow, I got to invent this bicycle, just to cover cases where UUID is not returned, meaning that the user is not just there, therefore invite is discarded the error issued and all burns.
	// The logic is: Invite user, asses if UUID is empty, if so -> find UUID of the invite, delete invite, issue an error.
	// PS: if you know better way, plz issue a PR, otherwise feel free to laugh.
	//var invitationCode EmailInvitationCodesResponse
	// I guess this is the correct way to see if the field in struct is empty.
	if teamMemberResponse.UID == "" {
		tflog.Info(ctx, "The user HAS NO vercel account, dropping invitation:", map[string]interface{}{})
		// then we gotta find the UUID of the request we just send
		invitationCode, err := c.GetTeamMemberInviteCode(ctx, request.Email)
		// Check again if it is empty, just in case of sudden #YOLO
		if invitationCode.ID == "" {
			err = errors.New("failed to get invite code")
			return TeamMemberResponse{}, err
		}

		tflog.Info(ctx, "[INVITE CODE] Got this one:", map[string]interface{}{
			"email": invitationCode.Email,
			"id":    invitationCode.ID,
		})
		err = c.DeleteTeamMemberInviteCode(ctx, invitationCode.ID)

		if err != nil {
			return TeamMemberResponse{}, err
		}

		// drop the mic and cover.
		return TeamMemberResponse{},
			errors.New(fmt.Sprintf("the vercel user must exist before adding it to the team. \n"+
				"user in question is: %s", request.Email))
	}

	return TeamMemberResponse{
		UID:   teamMemberResponse.UID,
		Email: teamMemberResponse.Email,
		Role:  teamMemberResponse.Role,
	}, nil

}
