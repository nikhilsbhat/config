package cli

import (
	"config/decode"
	"config/gcp"
	"config/version"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type gcloudAuth struct {
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
	JSONPath            string
	k8clusterName       string
	region              string
	version             string
}

func init() {

	path := os.Getenv("CONFIG_DATA")
	if path == "" {
		cm.NeuronSaysItsError("Could not find the variable what you are searching for")
	}

}

func versionConfig(cmd *cobra.Command, args []string) error {
	fmt.Println("Config", version.GetVersion())
	return nil
}

func configSet(auth gcloudAuth) error {

	if auth.JSONPath != "" {
		auth.fillGcloudAuth()
	} else {
		path, err := getJSONPathFromEnv()
		if err != nil {
			return err
		}
		auth.JSONPath = path
		auth.fillGcloudAuth()
	}

	if jsErr := auth.setServiceAccount(); jsErr != nil {
		cm.NeuronSaysItsError(fmt.Sprintf("An Error occured while setting service account %s", getStringOfMessage(jsErr)))
	}
	if spErr := auth.setProject(); spErr != nil {
		cm.NeuronSaysItsError(fmt.Sprintf("An Error occured while setting up gcp project %s", getStringOfMessage(spErr)))
	}
	if gcErr := auth.getClusterName(); gcErr != nil {
		cm.NeuronSaysItsError(fmt.Sprintf("An Error occured while fetching cluster name %s", getStringOfMessage(gcErr)))
	}
	/*if scErr := auth.setContainerCredentials(); scErr != nil {
		cm.NeuronSaysItsError(fmt.Sprintf("An Error occured while setting cluster credentials %s", getStringOfMessage(scErr)))
	}*/
	return nil
}

func (g *gcloudAuth) setServiceAccount() error {
	_, err := exec.Command("gcloud", "auth", "activate-service-account", g.ClientEmail, fmt.Sprintf("--key-file=%s", g.JSONPath)).Output()
	if err != nil {
		return err
	}
	cm.NeuronSaysItsInfo("Service account is set successfully")
	return nil
}

func (g *gcloudAuth) setProject() error {
	_, err := exec.Command("gcloud", "config", "set", "project", g.ProjectID).Output()
	if err != nil {
		return err
	}
	cm.NeuronSaysItsInfo("Project is set successfully")
	return nil
}

func (g *gcloudAuth) getClusterName() error {
	cluster := new(gcp.GetClusterInput)
	cluster.ProjectID = g.ProjectID
	cluster.JSONPath = g.JSONPath

	clusters, err := cluster.GetClusters()
	if err != nil {
		return err
	}

	ver, err := strconv.Atoi(g.version)
	if err != nil {
		return fmt.Errorf("Unable to parse the version you passed, please pass a valid one")
	}
	if ver == 1 {
		for _, cluster := range clusters {
			fmt.Println(g.ProjectID, cluster.Name)
			if pro := strings.Contains(g.ProjectID, cluster.Name); pro == true {
				fmt.Println(fmt.Sprintf("The cluster is: %s in the region: %s, %s", cluster.Name, cluster.Location, cluster.Locations))
			}
		}
	} else {
		for _, cluster := range clusters {
			fmt.Println(fmt.Sprintf("The cluster is: %s in the region: %s, %s", cluster.Name, cluster.Location, cluster.Locations))
		}
	}
	return nil
}

func (g *gcloudAuth) setContainerCredentials() error {
	_, err := exec.Command("gcloud", "container", "clusters", "get-credentials", g.k8clusterName, "--region", g.region).Output()
	if err != nil {
		return err
	}
	cm.NeuronSaysItsInfo("K8 cluster credentials is set successfully")
	return nil
}

func getStringOfMessage(g interface{}) string {
	switch g.(type) {
	case string:
		return g.(string)
	case error:
		return g.(error).Error()
	default:
		return "unknown messagetype"
	}
}

func getJSONPathFromEnv() (string, error) {

	path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if path == "" {
		cm.NeuronSaysItsError("Could not find the variable what you are searching for")
		return "", fmt.Errorf("Unable to find the credentials")
	}

	return path, nil
}

func (g *gcloudAuth) fillGcloudAuth() error {
	jsonCont, err := decode.ReadFile(g.JSONPath)
	if err != nil {
		return err
	}

	if decodneuerr := decode.JsonDecode(jsonCont, &g); decodneuerr != nil {
		cm.NeuronSaysItsError("Error Decoding JSON to gcloudAuth")
		return decodneuerr
	}
	return nil
}
