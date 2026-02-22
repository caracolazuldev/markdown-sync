package google

import (
    "context"
    "net/http"
)

// Placeholder: Auth client provider. Implement OAuth2 and service-account flows here.
func NewHTTPClient(ctx context.Context, authMode string, credentialsPath string) (*http.Client, error) {
    // TODO: implement OAuth2 + service account token sources.
    return http.DefaultClient, nil
}
