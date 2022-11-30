package vercel

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"vercel-mgt/client"
)

func newMembersDataSource() datasource.DataSource {
	return &membersDataSource{}
}

type membersDataSource struct {
	client *client.Client
}

func (d *membersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_members"
}

func (d *membersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (r membersDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides information about an existing team members`,
		Attributes: map[string]tfsdk.Attribute{
			"members": {
				Computed: true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"Uid": {
						Computed: true,
						Type:     types.StringType,
					},
					"Email": {
						Computed: true,
						Type:     types.StringType,
					},
					"Role": {
						Computed: true,
						Type:     types.StringType,
					},
				}),
			},
		},
	}, nil
}

type TeamData struct {
	Members []Team `tfsdk:"members"`
}

type Team struct {
	Email types.String      `tfsdk:"email"`
	Uid   types.String      `tfsdk:"uid"`
	Role  map[string]string `tfsdk:"role"`
}

// Read will read a file from the filesytem and provide terraform with information about it.
// It is called by the provider whenever data source values should be read to update state.
func (d *membersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeamData

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := d.client.GetTeamMembers(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading members data source",
			fmt.Sprintf("Could not read team members unexpected error: %s",
				err,
			),
		)
		return
	}

	//result := convertResponseToProject(out, config.coercedFields(), types.SetNull(envVariableElemType))
	//tflog.Trace(ctx, "read project", map[string]interface{}{
	//	"team_id":    result.TeamID.ValueString(),
	//	"project_id": result.ID.ValueString(),
	//})

	diags = resp.State.Set(ctx, out)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
