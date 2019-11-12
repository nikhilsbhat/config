package gcp

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/config/decode"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/compute/v1"
)

/*type auth struct {
	AuthScopes []string
}*/

type gcpSVCred struct {
	Type                string `json:"type,omitempty"`
	ProjectID           string `json:"project_id,omitempty"`
	PrivateKeyID        string `json:"private_key_id,omitempty"`
	PrivateKey          string `json:"private_key,omitempty"`
	ClientEmail         string `json:"client_email,omitempty"`
	ClientID            string `json:"client_id,omitempty"`
	AuthURI             string `json:"auth_uri,omitempty"`
	TokenURI            string `json:"token_uri,omitempty"`
	AuthProviderCertURL string `json:"auth_provider_x509_cert_url,omitempty"`
	ClientCertURL       string `json:"client_x509_cert_url,omitempty"`
}

type gcloudAuth struct {
	GCPSVCauth *gcpSVCred
	ProjectID  string
	AuthScopes []string
	JSONPath   string
	Zone       string
	RawJSON    []byte
	Client     *http.Client
}

func (auth *gcloudAuth) getClient() (*http.Client, error) {
	if auth.JSONPath == "" {
		client, err := auth.getDefalutClient()
		if err != nil {
			return nil, fmt.Errorf("Unable to initialize the default gcp client")
		}
		return client, nil
	}

	if err := auth.getCred(); err != nil {
		return nil, err
	}
	if client := auth.getCustomClient(); client != nil {
		return client, nil
	}
	return nil, fmt.Errorf("Unable to initialize the gcp client")
}

func (auth *gcloudAuth) getDefalutClient() (*http.Client, error) {

	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (auth *gcloudAuth) getCustomClient() *http.Client {

	conf := &jwt.Config{
		Email:      auth.GCPSVCauth.ClientEmail,
		PrivateKey: []byte(auth.GCPSVCauth.PrivateKey),
		Scopes:     auth.AuthScopes,
		TokenURL:   auth.GCPSVCauth.TokenURI,
		Subject:    auth.GCPSVCauth.ClientEmail,
	}

	client := conf.Client(oauth2.NoContext)
	return client
}

func (auth *gcloudAuth) getCred() error {
	jsonCont, err := decode.ReadFile(auth.JSONPath)
	if err != nil {
		return err
	}

	var jsonAuth gcpSVCred

	if decodneuerr := decode.JsonDecode(jsonCont, &jsonAuth); decodneuerr != nil {
		fmt.Println("Error Decoding JSON to gcpSVCred")
		return decodneuerr
	}

	auth.GCPSVCauth = &jsonAuth
	return nil
}
