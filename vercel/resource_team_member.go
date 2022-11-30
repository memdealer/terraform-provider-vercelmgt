package vercel

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"vercel-mgt/client"
)

func newTeamMemberResource() resource.Resource {
	return &teamMemberResource{}
}

type teamMemberResource struct {
	client *client.Client
}

func (r *teamMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_member"
}

func (r *teamMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r teamMemberResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Retrieves information about member user in given team`,
		Attributes: map[string]tfsdk.Attribute{

			"uid": {
				Type:          types.StringType,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown()},
			},
			"email": {
				Required: true,
				Type:     types.StringType,
			},
			"role": {
				Required: true,
				Type:     types.StringType,
			},
		},
	}, nil
}

type resourceTeamMember struct {
	p vercelProvider
}

func (r *teamMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeamMember

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.client.GetTeamMember(ctx, state.Email.ValueString())
	if client.NotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading team member",
			fmt.Sprintf("Could not get team member %s %s, unexpected error: %s",
				state.UID.ValueString(),
				state.Email.ValueString(),
				err,
			),
		)
		return
	}

	result := convertResponseToTeamMember(out)
	tflog.Trace(ctx, "read team member", map[string]interface{}{
		"email": result.Email.ValueString(),
		"role":  result.Role.ValueString(),
		"UID":   result.UID.ValueString(),
	})

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *teamMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeamMember

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.client.CreateTeamMember(ctx, client.CreateTeamMemberRequest{
		Email: plan.Email.ValueString(),
		Role:  plan.Role.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating team-member",
			"Could not create team member, unexpected error: "+err.Error(),
		)
		return
	}

	result := convertResponseToTeamMember(out)
	tflog.Trace(ctx, "created team member", map[string]interface{}{
		"email": plan.Email.ValueString(),
		"role":  plan.Role.ValueString(),
		"UID":   result.UID.ValueString(),
	})

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *teamMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeamMember

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTeamMember(ctx, state.UID.ValueString())

	if client.NotFound(err) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting team member",
			fmt.Sprintf(
				"Could not delete team member %s, unexpected error: %s",
				state.UID.ValueString(),
				err,
			),
		)
		return
	}

	tflog.Trace(ctx, "deleted team member", map[string]interface{}{
		"UUID":  state.UID.ValueString(),
		"EMAIL": state.Email.ValueString(),
	})
}

func (r *teamMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeamMember
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateTeamMember(
		ctx,
		plan.UID.ValueString(),
		plan.toUpdateRequest(),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating team member",
			fmt.Sprintf("Could not update team member  %s  unexpected error: %s",
				plan.UID.ValueString(),
				err,
			),
		)
		return
	}

	result := TeamMember{ // The Vercel returns a team ID from a request,
		// I couldn't figure out more advanced way to overcome it,
		// rather than just append it to state like that.
		UID:   plan.UID,
		Email: plan.Email,
		Role:  plan.Role,
	}
	tflog.Trace(ctx, "update project domain", map[string]interface{}{
		"Role":  result.Role.ValueString(),
		"Email": result.Email.ValueString(),
		"UID":   result.UID.ValueString(),
	})

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *teamMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	email := req.ID

	out, err := r.client.GetTeamMember(ctx, email)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading team member",
			fmt.Sprintf("Could not get team member %s, unexpected error: %s",
				email,
				err,
			),
		)
		return
	}

	result := convertResponseToTeamMember(out)

	tflog.Trace(ctx, "imported team member", map[string]interface{}{
		"UUID":  result.UID.ValueString(),
		"EMAIL": result.Email.ValueString(),
	})

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
