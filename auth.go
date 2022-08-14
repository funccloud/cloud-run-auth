package cloudrunauth

import (
	"context"
	"fmt"

	"cloud.google.com/go/compute/metadata"
)

type MetadataServerToken struct {
	ServiceURL string
}

// GetRequestMetadata is called on every request, so we are sure that token is always not expired
func (t MetadataServerToken) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	// based on https://cloud.google.com/run/docs/authenticating/service-to-service#go
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", t.ServiceURL)
	idToken, err := metadata.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("cannot query id token for gRPC: %w", err)
	}

	return map[string]string{
		"authorization": "Bearer " + idToken,
	}, nil
}

func (MetadataServerToken) RequireTransportSecurity() bool {
	return true
}
