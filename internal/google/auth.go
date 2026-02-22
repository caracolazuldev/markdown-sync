package google

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "os"
    "path/filepath"

    "golang.org/x/oauth2"
    oauthGoogle "golang.org/x/oauth2/google"
)

var (
    // default scopes used by the CLI; callers can request narrower scopes if needed
    defaultScopes = []string{"https://www.googleapis.com/auth/documents"}
)

// NewHTTPClient returns an authenticated *http.Client.
// authMode: "oauth" (interactive) or "service" (service account). When using
// service mode, credentialsPath must point to a service account JSON key.
func NewHTTPClient(ctx context.Context, authMode string, credentialsPath string) (*http.Client, error) {
    switch authMode {
    case "service":
        if credentialsPath == "" {
            // allow GOOGLE_APPLICATION_CREDENTIALS fallback
            credentialsPath = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
            if credentialsPath == "" {
                return nil, errors.New("service auth requires credentials path or GOOGLE_APPLICATION_CREDENTIALS")
            }
        }
        data, err := os.ReadFile(credentialsPath)
        if err != nil {
            return nil, fmt.Errorf("reading service account key: %w", err)
        }
        cfg, err := oauthGoogle.JWTConfigFromJSON(data, defaultScopes...)
        if err != nil {
            return nil, fmt.Errorf("parsing service account key: %w", err)
        }
        return cfg.Client(ctx), nil
    case "oauth":
        // interactive OAuth2 flow. Client ID/Secret may be supplied via env vars
        clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
        clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
        if clientID == "" || clientSecret == "" {
            return nil, errors.New("GOOGLE_OAUTH_CLIENT_ID and GOOGLE_OAUTH_CLIENT_SECRET must be set for oauth auth mode")
        }
        conf := &oauth2.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            Endpoint:     oauthGoogle.Endpoint,
            RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
            Scopes:       defaultScopes,
        }

        tokenFile := tokenCacheFile()
        tok, err := tokenFromFile(tokenFile)
        if err != nil {
            tok, err = getTokenFromWeb(conf)
            if err != nil {
                return nil, err
            }
            if err := saveToken(tokenFile, tok); err != nil {
                // non-fatal: warn but continue
                fmt.Fprintf(os.Stderr, "warning: saving token: %v\n", err)
            }
        }
        return conf.Client(ctx, tok), nil
    default:
        return nil, fmt.Errorf("unknown auth mode: %s", authMode)
    }
}

func tokenCacheFile() string {
    cfg := os.Getenv("XDG_CONFIG_HOME")
    if cfg == "" {
        home := os.Getenv("HOME")
        cfg = filepath.Join(home, ".config")
    }
    dir := filepath.Join(cfg, "markdown-sync")
    _ = os.MkdirAll(dir, 0o700)
    return filepath.Join(dir, "tokens.json")
}

func tokenFromFile(path string) (*oauth2.Token, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var tok oauth2.Token
    if err := json.NewDecoder(f).Decode(&tok); err != nil {
        return nil, err
    }
    return &tok, nil
}

func saveToken(path string, token *oauth2.Token) error {
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
    if err != nil {
        return err
    }
    defer f.Close()
    return json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(conf *oauth2.Config) (*oauth2.Token, error) {
    authURL := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Fprintf(os.Stderr, "Go to the following link in your browser then type the authorization code:\n%v\n", authURL)
    var code string
    if _, err := fmt.Scan(&code); err != nil {
        return nil, err
    }
    tok, err := conf.Exchange(context.Background(), code)
    if err != nil {
        return nil, err
    }
    return tok, nil
}
