package provider

import (
	"context"
	"os"

	"terraform-provider-flashduty/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &FlashdutyProvider{}

type FlashdutyProvider struct {
	version string
}

type FlashdutyProviderModel struct {
	AppKey  types.String `tfsdk:"app_key"`
	BaseURL types.String `tfsdk:"base_url"`
}

func (p *FlashdutyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "flashduty"
	resp.Version = p.version
}

func (p *FlashdutyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Flashduty provider is used to interact with Flashduty API resources.",
		Attributes: map[string]schema.Attribute{
			"app_key": schema.StringAttribute{
				MarkdownDescription: "The APP Key for Flashduty API. Can be obtained from Flashduty Console -> Account Settings -> APP Key. Can also be set via `FLASHDUTY_APP_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The base URL for Flashduty API. Defaults to `https://api.flashcat.cloud`. Can also be set via `FLASHDUTY_BASE_URL` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *FlashdutyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config FlashdutyProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for unknown values
	if config.AppKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("app_key"),
			"Unknown Flashduty APP Key",
			"The provider cannot create the Flashduty API client as there is an unknown configuration value for the Flashduty APP key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FLASHDUTY_APP_KEY environment variable.",
		)
	}

	if config.BaseURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Unknown Flashduty Base URL",
			"The provider cannot create the Flashduty API client as there is an unknown configuration value for the Flashduty base URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FLASHDUTY_BASE_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Get values from config or environment variables
	appKey := os.Getenv("FLASHDUTY_APP_KEY")
	baseURL := os.Getenv("FLASHDUTY_BASE_URL")

	if !config.AppKey.IsNull() {
		appKey = config.AppKey.ValueString()
	}

	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}

	// Validate required configuration
	if appKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("app_key"),
			"Missing Flashduty APP Key",
			"The provider cannot create the Flashduty API client as there is a missing or empty value for the Flashduty APP key. "+
				"Set the app_key value in the configuration or use the FLASHDUTY_APP_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create client with options
	var opts []client.ClientOption
	if baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}

	c := client.NewClient(appKey, p.version, opts...)

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *FlashdutyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMemberInviteResource,
		NewTeamResource,
		NewChannelResource,
		NewScheduleResource,
		NewIncidentResource,
		NewEscalateRuleResource,
		NewSilenceRuleResource,
		NewInhibitRuleResource,
		NewFieldResource,
		NewRouteResource,
		NewAlertPipelineResource,
		NewTemplateResource,
	}
}

func (p *FlashdutyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTeamDataSource,
		NewTeamsDataSource,
		NewChannelDataSource,
		NewChannelsDataSource,
		NewMemberDataSource,
		NewMembersDataSource,
		NewFieldDataSource,
		NewFieldsDataSource,
		NewRouteDataSource,
		NewRouteHistoryDataSource,
		NewTemplateDataSource,
		NewTemplatesDataSource,
		NewAlertPipelineDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FlashdutyProvider{
			version: version,
		}
	}
}
