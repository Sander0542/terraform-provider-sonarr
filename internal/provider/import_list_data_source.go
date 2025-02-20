package provider

import (
	"context"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const importListDataSourceName = "import_list"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ImportListDataSource{}

func NewImportListDataSource() datasource.DataSource {
	return &ImportListDataSource{}
}

// ImportListDataSource defines the import_list implementation.
type ImportListDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

func (d *ImportListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListDataSourceName
}

func (d *ImportListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Download Clients -->\nSingle [Download Client](../resources/import_list).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_add": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic add flag.",
				Computed:            true,
			},
			"season_folder": schema.BoolAttribute{
				MarkdownDescription: "Season folder flag.",
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Computed:            true,
			},
			"should_monitor": schema.StringAttribute{
				MarkdownDescription: "Should monitor.",
				Computed:            true,
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Root folder path.",
				Computed:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "ImportList implementation name.",
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "ImportList configuration template.",
				Computed:            true,
			},
			"series_type": schema.StringAttribute{
				MarkdownDescription: "Series type.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Import List name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Import List ID.",
				Computed:            true,
			},
			// Field values
			"limit": schema.Int64Attribute{
				MarkdownDescription: "Limit.",
				Computed:            true,
			},
			"trakt_list_type": schema.Int64Attribute{
				MarkdownDescription: "Trakt list type.",
				Computed:            true,
			},
			"list_type": schema.Int64Attribute{
				MarkdownDescription: "Simkl list type.",
				Computed:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Computed:            true,
				Sensitive:           true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Computed:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Computed:            true,
				Sensitive:           true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth User.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Computed:            true,
			},
			"rating": schema.StringAttribute{
				MarkdownDescription: "Rating.",
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL.",
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Computed:            true,
			},
			"listname": schema.StringAttribute{
				MarkdownDescription: "List name.",
				Computed:            true,
			},
			"list_id": schema.StringAttribute{
				MarkdownDescription: "List ID.",
				Computed:            true,
			},
			"genres": schema.StringAttribute{
				MarkdownDescription: "Genres.",
				Computed:            true,
			},
			"years": schema.StringAttribute{
				MarkdownDescription: "Years.",
				Computed:            true,
			},
			"trakt_additional_parameters": schema.StringAttribute{
				MarkdownDescription: "Trakt additional parameters.",
				Computed:            true,
			},
			"language_profile_ids": schema.SetAttribute{
				MarkdownDescription: "Language profile IDs.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"quality_profile_ids": schema.SetAttribute{
				MarkdownDescription: "Quality profile IDs.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"tag_ids": schema.SetAttribute{
				MarkdownDescription: "Tag IDs.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (d *ImportListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *ImportListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ImportList

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get importList current value
	response, _, err := d.client.ImportListAPI.ListImportList(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListDataSourceName, err))

		return
	}

	data.find(ctx, data.Name.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+importListDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (i *ImportList) find(ctx context.Context, name string, importLists []sonarr.ImportListResource, diags *diag.Diagnostics) {
	for _, list := range importLists {
		if list.GetName() == name {
			i.write(ctx, &list, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(importListDataSourceName, "name", name))
}
