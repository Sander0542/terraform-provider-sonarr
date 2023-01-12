package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const importListExclusionDataSourceName = "import_list_exclusion"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ImportListExclusionDataSource{}

func NewImportListExclusionDataSource() datasource.DataSource {
	return &ImportListExclusionDataSource{}
}

// ImportListExclusionDataSource defines the importListExclusion implementation.
type ImportListExclusionDataSource struct {
	client *sonarr.APIClient
}

func (d *ImportListExclusionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListExclusionDataSourceName
}

func (d *ImportListExclusionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Import Lists -->Single [ImportListExclusion](../resources/import_list_exclusion).",
		Attributes: map[string]schema.Attribute{
			"tvdb_id": schema.Int64Attribute{
				MarkdownDescription: "Series TVDB ID.",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Series to be excluded.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "ImportListExclusion ID.",
				Computed:            true,
			},
		},
	}
}

func (d *ImportListExclusionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ImportListExclusionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var importListExclusion *ImportListExclusion

	resp.Diagnostics.Append(req.Config.Get(ctx, &importListExclusion)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get importListExclusions current value
	response, _, err := d.client.ImportListExclusionApi.ListImportListExclusion(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListExclusionDataSourceName, err))

		return
	}

	value, err := findImportListExclusion(importListExclusion.TVDBID.ValueInt64(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", importListExclusionDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListExclusionDataSourceName)
	importListExclusion.write(value)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &importListExclusion)...)
}

func findImportListExclusion(tvID int64, importListExclusions []*sonarr.ImportListExclusionResource) (*sonarr.ImportListExclusionResource, error) {
	for _, t := range importListExclusions {
		if t.GetTvdbId() == int32(tvID) {
			return t, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(importListExclusionDataSourceName, "tvdb_id", strconv.Itoa(int(tvID)))
}
