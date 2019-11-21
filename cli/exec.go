package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nikhilsbhat/config/decode"
	"github.com/nikhilsbhat/config/gcp"
	"github.com/nikhilsbhat/config/version"

	"github.com/nikhilsbhat/neuron/cli/ui"
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
	jsonPath            string
	k8clusterName       string
	regions             []string
	version             string
	shell               *bufio.Reader
}

var (
	jsonAuth gcloudAuth
)

// This is temporary function created for testing and will be removed soon if its purpose fails.
func init() {
	jsonAuth.shell = bufio.NewReader(os.Stdin)
}

func versionConfig(cmd *cobra.Command, args []string) error {
	fmt.Println("Config", version.GetVersion())
	return nil
}

func configSet(auth gcloudAuth) error {

	/*if path == "" {
		return fmt.Errorf("Could not find the variable what you are searching for")
	}*/

	if auth.jsonPath != "" {
		auth.fillGcloudAuth()
	} else {
		path, err := getJSONPathFromEnv()
		if err != nil {
			return err
		}
		auth.jsonPath = path
		auth.fillGcloudAuth()
	}

	if jsErr := auth.setServiceAccount(); jsErr != nil {
		return fmt.Errorf(fmt.Sprintf("An Error occurred while setting service account: %s\n", getStringOfMessage(jsErr)))
	}
	if spErr := auth.setProject(); spErr != nil {
		return fmt.Errorf(fmt.Sprintf("An Error occurred while setting up gcp project: %s\n", getStringOfMessage(spErr)))
	}
	if gcErr := auth.getClusterName(); gcErr != nil {
		return fmt.Errorf(fmt.Sprintf("An Error occurred while fetching cluster name: %s\n", getStringOfMessage(gcErr)))
	}
	if scErr := auth.setContainerCredentials(); scErr != nil {
		return fmt.Errorf(fmt.Sprintf("An Error occurred while setting cluster credentials %s\n", getStringOfMessage(scErr)))
	}
	return nil
}

func (g *gcloudAuth) setServiceAccount() error {
	_, err := exec.Command("gcloud", "auth", "activate-service-account", g.ClientEmail, fmt.Sprintf("--key-file=%s", g.jsonPath)).Output()
	if err != nil {
		return err
	}
	cm.NeuronSaysItsInfo("Service account is set successfully\n")
	return nil
}

func (g *gcloudAuth) setProject() error {
	_, err := exec.Command("gcloud", "config", "set", "project", g.ProjectID).Output()
	if err != nil {
		return err
	}
	cm.NeuronSaysItsInfo("Project is set successfully\n")
	return nil
}

func (g *gcloudAuth) getClusterName() error {
	cluster := new(gcp.GetClusterInput)
	cluster.ProjectID = g.ProjectID
	cluster.JSONPath = g.jsonPath
	cluster.ClusterName = g.k8clusterName
	cluster.Regions = g.regions

	// Fetches the details of the selected cluster.
	if len(cluster.ClusterName) != 0 {
		clusters, err := cluster.GetCluster()
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("Selected cluster is :%s in the region :%s\n", ui.Info(clusters.Name), ui.Info(clusters.Location)))
		g.k8clusterName = clusters.Name
		g.regions = []string{clusters.Location}
		return nil
	}

	// Fetches the details of all available cluster in the selected region.
	clusters, err := cluster.GetClusters()
	if err != nil {
		return err
	}

	// Avoiding shell if only one cluster exists.
	if len(clusters) == 1 {
		g.k8clusterName = clusters[0].Name
		g.regions = []string{clusters[0].Location}
		return nil
	}

	clust := make(map[string]string)
	for _, cluster := range clusters {
		fmt.Println(fmt.Sprintf("The cluster is: %s in the region: %s\n", ui.Info(cluster.Name), ui.Info(cluster.Location)))
		if clust[cluster.Name] != cluster.Location {
			clust[cluster.Name] = cluster.Location
		}
	}

	var clusterselec string
	for ok := true; ok; ok = (len(clust[clusterselec]) == 0) {
		clusterselec, err = getClusterFromIntr()
		if err != nil {
			return err
		}
		if len(clust[clusterselec]) == 0 {
			cm.NeuronSaysItsWarn("The cluster selected doesnot exists, please make a valid selection\n")
		}
	}

	fmt.Println(fmt.Sprintf("Selected cluster is :%s in the region :%s\n", ui.Info(clusterselec), ui.Info(clust[clusterselec])))
	if stat := getConfirmOfCLuster(); stat == false {
		cm.NeuronSaysItsInfo("you opted no, I'm backing off")
		os.Exit(1)
	}

	g.k8clusterName = clusterselec
	g.regions = []string{clust[clusterselec]}

	return nil
}

func (g *gcloudAuth) setContainerCredentials() error {
	_, err := exec.Command("gcloud", "container", "clusters", "get-credentials", g.k8clusterName, "--region", g.regions[0]).Output()
	if err != nil {
		return err
	}
	cm.NeuronSaysItsInfo("K8 cluster credentials is set successfully\n")
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
	jsonCont, err := decode.ReadFile(g.jsonPath)
	if err != nil {
		return err
	}

	if decodneuerr := decode.JsonDecode(jsonCont, &g); decodneuerr != nil {
		cm.NeuronSaysItsError("Error Decoding JSON to gcloudAuth")
		return decodneuerr
	}
	return nil
}

func getConfirmOfCLuster() bool {

	for {
		fmt.Print(ui.Debug("$config>> "))
		fmt.Print(ui.Debug("you want to switch to the cluster selected ? [yes/no]: "))
		cmdString, err := (jsonAuth.shell).ReadString('\n')
		if err != nil {
			return false
		}

		if len(cmdString) <= 1 {
			cm.NeuronSaysItsWarn("did not get any valid input")
			return false
		}
		// Have to implement the wait function until valid input is passed
		cmnds := getArrayOfEntries(cmdString)
		if (cmnds[0] == "yes") || (cmnds[0] == "y") {
			return true
		} else if (cmnds[0] == "no") || (cmnds[0] == "n") {
			return false
		} else {
			cm.NeuronSaysItsWarn("opt eithier yes/no")
			return false
		}
	}
}

func getClusterFromIntr() (string, error) {

	for {
		fmt.Print(ui.Debug("$config>> "))
		fmt.Print(ui.Debug("select cluster from above list [multiple entry not accepted]: "))
		cmdString, err := (jsonAuth.shell).ReadString('\n')
		if err != nil {
			return "", err
		}

		if len(cmdString) <= 1 {
			return "", fmt.Errorf("Selection of cluster cannot be empty")
		}
		if cmdlen := len(getArrayOfEntries(cmdString)) > 1; cmdlen == true {
			return "", fmt.Errorf("Cannot select multiple cluster for this operation")
		}
		return strings.Join(getArrayOfEntries(cmdString)[:1], ""), nil
	}
}

func getArrayOfEntries(commandStr string) []string {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	return strings.Fields(commandStr)
}
