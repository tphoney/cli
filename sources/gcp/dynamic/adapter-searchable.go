package dynamic

import (
	"context"
	"fmt"
	"strings"

	"github.com/overmindtech/cli/discovery"
	"github.com/overmindtech/cli/sdp-go"
	gcpshared "github.com/overmindtech/cli/sources/gcp/shared"
)

// SearchableAdapter implements discovery.SearchableAdapter for GCP dynamic adapters.
type SearchableAdapter struct {
	customSearchMethodDesc string
	searchURLFunc          gcpshared.EndpointFunc
	Adapter
}

// NewSearchableAdapter creates a new GCP dynamic adapter.
func NewSearchableAdapter(searchURLFunc gcpshared.EndpointFunc, config *AdapterConfig, customSearchMethodDesc string) (discovery.SearchableAdapter, error) {
	a := Adapter{
		projectID:           config.ProjectID,
		scope:               config.Scope,
		httpCli:             config.HTTPClient,
		getURLFunc:          config.GetURLFunc,
		sdpAssetType:        config.SDPAssetType,
		sdpAdapterCategory:  config.SDPAdapterCategory,
		terraformMappings:   config.TerraformMappings,
		linker:              config.Linker,
		potentialLinks:      potentialLinksFromBlasts(config.SDPAssetType, gcpshared.BlastPropagations),
		uniqueAttributeKeys: config.UniqueAttributeKeys,
	}

	if a.httpCli == nil {
		gcpHTTPCliWithOtel, err := gcpshared.GCPHTTPClientWithOtel()
		if err != nil {
			return nil, err
		}

		a.httpCli = gcpHTTPCliWithOtel
	}

	return SearchableAdapter{
		customSearchMethodDesc: customSearchMethodDesc,
		searchURLFunc:          searchURLFunc,
		Adapter:                a,
	}, nil
}

func (g SearchableAdapter) Metadata() *sdp.AdapterMetadata {
	return &sdp.AdapterMetadata{
		Type:            g.sdpAssetType.String(),
		Category:        g.sdpAdapterCategory,
		DescriptiveName: g.sdpAssetType.Readable(),
		SupportedQueryMethods: &sdp.AdapterSupportedQueryMethods{
			Get:               true,
			GetDescription:    getDescription(g.sdpAssetType, g.scope, g.uniqueAttributeKeys),
			Search:            true,
			SearchDescription: searchDescription(g.sdpAssetType, g.scope, g.uniqueAttributeKeys, g.customSearchMethodDesc),
		},
		TerraformMappings: g.terraformMappings,
		PotentialLinks:    g.potentialLinks,
	}
}

func (g SearchableAdapter) Search(ctx context.Context, scope, query string, ignoreCache bool) ([]*sdp.Item, error) {
	if scope != g.scope {
		return nil, &sdp.QueryError{
			ErrorType:   sdp.QueryError_NOSCOPE,
			ErrorString: fmt.Sprintf("requested scope %v does not match any adapter scope %v", scope, g.Scopes()),
		}
	}

	if strings.HasPrefix(query, "projects/") {
		// This must be a terraform query in the format of:
		// projects/{{project}}/datasets/{{dataset}}/tables/{{name}}
		// projects/{{project}}/serviceAccounts/{{account}}/keys/{{key}}
		return terraformMappingViaSearch(ctx, g.Adapter, query)
	}

	// This is a regular SEARCH call
	url := g.searchURLFunc(query)
	if url == "" {
		return nil, &sdp.QueryError{
			ErrorType: sdp.QueryError_OTHER,
			ErrorString: fmt.Sprintf(
				"failed to construct the URL for the query \"%s\". SEARCH method description: %s",
				query,
				g.Metadata().GetSupportedQueryMethods().GetSearchDescription(),
			),
		}
	}

	return aggregateSDPItems(ctx, g.Adapter, url)
}

func (g SearchableAdapter) SearchStream(ctx context.Context, scope, query string, ignoreCache bool, stream discovery.QueryResultStream) {
	if scope != g.scope {
		stream.SendError(&sdp.QueryError{
			ErrorType:   sdp.QueryError_NOSCOPE,
			ErrorString: fmt.Sprintf("requested scope %v does not match any adapter scope %v", scope, g.Scopes()),
		})
		return
	}

	if strings.HasPrefix(query, "projects/") {
		// This must be a terraform query in the format of:
		// projects/{{project}}/datasets/{{dataset}}/tables/{{name}}
		// projects/{{project}}/serviceAccounts/{{account}}/keys/{{key}}
		items, err := terraformMappingViaSearch(ctx, g.Adapter, query)
		if err != nil {
			stream.SendError(&sdp.QueryError{
				ErrorType:   sdp.QueryError_OTHER,
				ErrorString: fmt.Sprintf("failed to execute terraform mapping search for query \"%s\": %v", query, err),
			})
			return
		}

		// There should only be one item in the result, so we can send it directly
		stream.SendItem(items[0])
		return
	}

	searchURL := g.searchURLFunc(query)
	if searchURL == "" {
		stream.SendError(&sdp.QueryError{
			ErrorType: sdp.QueryError_OTHER,
			ErrorString: fmt.Sprintf(
				"failed to construct the URL for the query \"%s\". SEARCH method description: %s",
				query,
				g.Metadata().GetSupportedQueryMethods().GetSearchDescription(),
			),
		})
		return
	}

	streamSDPItems(ctx, g.Adapter, searchURL, stream)
}
