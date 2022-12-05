package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const allSeriesDataSourceName = "all_series"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AllSeriessDataSource{}

func NewAllSeriessDataSource() datasource.DataSource {
	return &AllSeriessDataSource{}
}

// AllSeriessDataSource defines the tags implementation.
type AllSeriessDataSource struct {
	client *sonarr.Sonarr
}

// AllSeriess describes the series(es) data model.
type SeriesList struct {
	Series types.Set    `tfsdk:"series"`
	ID     types.String `tfsdk:"id"`
}

func (d *AllSeriessDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + allSeriesDataSourceName
}

func (d *AllSeriessDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Series -->List all available [Series](../resources/series).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"series": schema.SetNestedAttribute{
				MarkdownDescription: "Series list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Series ID.",
							Computed:            true,
						},
						"title": schema.StringAttribute{
							MarkdownDescription: "Series Title.",
							Computed:            true,
						},
						"title_slug": schema.StringAttribute{
							MarkdownDescription: "Series Title in kebab format.",
							Computed:            true,
						},
						"monitored": schema.BoolAttribute{
							MarkdownDescription: "Monitored flag.",
							Computed:            true,
						},
						"season_folder": schema.BoolAttribute{
							MarkdownDescription: "Season Folder flag.",
							Computed:            true,
						},
						"use_scene_numbering": schema.BoolAttribute{
							MarkdownDescription: "Scene numbering flag.",
							Computed:            true,
						},
						"language_profile_id": schema.Int64Attribute{
							MarkdownDescription: "Language Profile ID .",
							Computed:            true,
						},
						"quality_profile_id": schema.Int64Attribute{
							MarkdownDescription: "Quality Profile ID.",
							Computed:            true,
						},
						"tvdb_id": schema.Int64Attribute{
							MarkdownDescription: "TVDB ID.",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							MarkdownDescription: "Series Path.",
							Computed:            true,
						},
						"root_folder_path": schema.StringAttribute{
							MarkdownDescription: "Series Root Folder.",
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (d *AllSeriessDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *AllSeriessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *SeriesList

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get series current value
	response, err := d.client.GetAllSeriesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", allSeriesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+allSeriesDataSourceName)
	// Map response body to resource schema attribute
	series := make([]Series, len(response))
	for i, t := range response {
		series[i].write(ctx, t)
	}

	tfsdk.ValueFrom(ctx, series, data.Series.Type(context.Background()), &data.Series)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
