package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"golang.org/x/oauth2"
)

type configHelperResp struct {
	Credential struct {
		AccessToken string `json:"access_token"`
		TokenExpiry string `json:"token_expiry"`
	} `json:"credential"`
	Configuration struct {
		Properties struct {
			Core struct {
				Project string `json:"project"`
			} `json:"core"`
		} `json:"properties"`
	} `json:"configuration"`
}

type configHelper func() (*configHelperResp, error)

// An SDKConfig provides access to tokens from an account already
// authorized via the Google Cloud SDK.
type SDKConfig struct {
	helper configHelper
}

func NewSDKConfig() (*SDKConfig, error) {
	helper := func() (*configHelperResp, error) {
		cmd := exec.Command("gcloud", "config", "config-helper", "--format=json")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return nil, err
		}
		var resp configHelperResp
		if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
			return nil, fmt.Errorf("failure parsing the output from the Cloud SDK config helper: %v", err)
		}
		return &resp, nil
	}
	return &SDKConfig{helper}, nil
}

// Token returns an oauth2.Token retrieved from the Google Cloud SDK.
func (c *SDKConfig) Token() (*oauth2.Token, error) {
	resp, err := c.helper()
	if err != nil {
		return nil, fmt.Errorf("failure invoking the Cloud SDK config helper: %v", err)
	}
	expiry, err := time.Parse(time.RFC3339, resp.Credential.TokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("failure parsing the access token expiration time: %v", err)
	}
	return &oauth2.Token{
		AccessToken: resp.Credential.AccessToken,
		Expiry:      expiry,
	}, nil
}

// ProjectID returns the project id retrieved from the Google Cloud SDK.
func (c *SDKConfig) ProjectID() (string, error) {
	resp, err := c.helper()
	if err != nil {
		return "", fmt.Errorf("failure invoking the Cloud SDK config helper: %v", err)
	}

	return resp.Configuration.Properties.Core.Project, nil
}
