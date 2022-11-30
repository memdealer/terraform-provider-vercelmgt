package vercel

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"vercel-mgt/client"
)

type TeamMember struct {
	UID   types.String `tfsdk:"uid"`
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

func convertResponseToTeamMember(response client.TeamMemberResponse) TeamMember {
	return TeamMember{
		UID:   types.StringValue(response.UID),
		Email: types.StringValue(response.Email),
		Role:  types.StringValue(response.Role),
	}
}

func (p *TeamMember) toUpdateRequest() client.UpdateTeamMemberRequest {
	return client.UpdateTeamMemberRequest{
		Role: toStrPointer(p.Role),
	}
}
