package geoid

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
)

// Parse the OpenID Connect Provider file and create an oauth2.Config
func providerFileParser(file_path string) (*oauth2.Config, string, error) {
	// Check if the file exists
	info, err := os.Stat(file_path)
	if err != nil {
		return nil, "", err
	}
	if info.IsDir() {
		return nil, "", os.ErrNotExist
	}

	// Read the file
	file_content, err := os.ReadFile(file_path)
	if err != nil {
		return nil, "", err
	}

	// Parse the file
	var provider_content ProviderContent
	err = json.Unmarshal(file_content, &provider_content)
	if err != nil {
		return nil, "", err
	}

	// Create the oauth2.Config
	config := &oauth2.Config{
		ClientID:     provider_content.ClientID,
		ClientSecret: provider_content.ClientSecret,
		RedirectURL:  provider_content.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:       provider_content.AuthorizationEndpoint,
			TokenURL:      provider_content.TokenEndpoint,
			DeviceAuthURL: provider_content.DeviceAuthorizationEndpoint,
			AuthStyle:     oauth2.AuthStyleInParams,
		},
		Scopes: provider_content.ScopesSupported,
	}

	return config, provider_content.UserinfoEndpoint, nil
}
