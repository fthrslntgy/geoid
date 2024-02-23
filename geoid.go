package geoid

import (
	"errors"

	"golang.org/x/oauth2"
)

// OAuthConfig returns an oauth2.Config object and user_info_endpoint from the given file. If the file does not exist, an error is returned.
// If the file exists, the client_id, client_secret, and redirect_url are read from the file. If the file does not contain these fields, they can be provided as parameters.
// If the parameters are provided, they can take priority over the file with the params_priority flag.
// If the file does not contain the fields and the parameters are not provided or not prior, these fields will be empty.
// Params should be provided in the following order: client_id, client_secret, redirect_url or be empty.
func OAuthConfig(filepath string, params_priority bool, params ...string) (*oauth2.Config, string, error) {
	if len(params) != 0 && len(params) != 3 {
		return nil, "", errors.New("invalid number of parameters")
	}
	config, user_info_endpoint, err := providerFileParser(filepath)
	if err != nil {
		return nil, "", err
	}

	// If the client_id, client_secret, or redirect_url are not set in the file, use the parameters if they are provided
	if len(params) == 3 {
		if config.ClientID == "" || params_priority {
			config.ClientID = params[0]
		}
		if config.ClientSecret == "" || params_priority {
			config.ClientSecret = params[1]
		}
		if config.RedirectURL == "" || params_priority {
			config.RedirectURL = params[2]
		}
	}
	return config, user_info_endpoint, nil
}
