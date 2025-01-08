package framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &fwprovider{}

// New returns a new, initialized Terraform Plugin Framework-style provider instance.
// The provider instance is fully configured once the `Configure` method has been called.
func New(primary interface{ Meta() interface{} }) provider.Provider {
	return &fwprovider{
		Primary: primary,
	}
}

type fwprovider struct {
	Primary interface{ Meta() interface{} }
}

func (f *fwprovider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "github"
}

func (f *fwprovider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional: true,
				Description: "The OAuth token used to connect to GitHub. Anonymous mode is enabled if both `token` and " +
					"`app_auth` are not set.",
			},
			"owner": schema.StringAttribute{
				Optional: true,
				Description: "The GitHub owner name to manage. " +
					"Use this field instead of `organization` when managing individual accounts.",
			},
			"retryable_errors": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				Description: "Allow the provider to retry after receiving an error status code, the max_retries should be set for this to work" +
					"Defaults to [500, 502, 503, 504]",
			},
			"max_retries": schema.Int64Attribute{
				Optional: true,
				Description: "Number of times to retry a request after receiving an error status code" +
					"Defaults to 3",
			},
			"organization": schema.StringAttribute{
				Optional: true,
				Description: "The GitHub organization name to manage. " +
					"Use this field instead of `owner` when managing organization accounts.",
				DeprecationMessage: "Use owner (or GITHUB_OWNER) instead of organization (or GITHUB_ORGANIZATION)",
			},
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "The GitHub Base API URL",
			},
			"insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable `insecure` mode for testing purposes",
			},
			"write_delay_ms": schema.Int64Attribute{
				Optional: true,
				Description: "Amount of time in milliseconds to sleep in between writes to GitHub API. " +
					"Defaults to 1000ms or 1s if not set.",
			},
			"read_delay_ms": schema.Int64Attribute{
				Optional: true,
				Description: "Amount of time in milliseconds to sleep in between non-write requests to GitHub API. " +
					"Defaults to 0ms if not set.",
			},
			"retry_delay_ms": schema.Int64Attribute{
				Optional: true,
				Description: "Amount of time in milliseconds to sleep in between requests to GitHub API after an error response. " +
					"Defaults to 1000ms or 1s if not set, the max_retries must be set to greater than zero.",
			},
			"parallel_requests": schema.BoolAttribute{
				Optional: true,
				Description: "Allow the provider to make parallel API calls to GitHub. " +
					"You may want to set it to true when you have a private Github Enterprise without strict rate limits. " +
					"Although, it is not possible to enable this setting on github.com " +
					"because we enforce the respect of github.com's best practices to avoid hitting abuse rate limits" +
					"Defaults to false if not set",
			},
		},
		Blocks: map[string]schema.Block{
			"app_auth": schema.ListNestedBlock{
				Description: "The GitHub App credentials used to connect to GitHub. Conflicts with " +
					"`token`. Anonymous mode is enabled if both `token` and `app_auth` are not set.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "The GitHub App ID.",
						},
						"installation_id": schema.StringAttribute{
							Required:    true,
							Description: "The GitHub App installation instance ID.",
						},
						"pem_file": schema.StringAttribute{
							Required:    true,
							Sensitive:   true,
							Description: "The GitHub App PEM file contents.",
						},
					},
				},
			},
		},
	}
}

func (f *fwprovider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	// Provider's parsed configuration (its instance state) is available through the primary provider's Meta() method.
	v := f.Primary.Meta()
	response.DataSourceData = v
	response.ResourceData = v
	response.EphemeralResourceData = v
}

func (f *fwprovider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (f *fwprovider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}
