// Package gcp helps in Performing the certain operation on GKE of google cloud
package gcp

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
)

// GetClusterInput holds the required values to fetch the cluster details
type GetClusterInput struct {
	ClusterName string
	ProjectID   string
	Regions     []string
	JSONPath    string
}

// GetClusters gets the list of kube clusters along with its other details
// This is region specific, mind the region you pass to it.
func (i GetClusterInput) GetClusters() ([]*container.Cluster, error) {

	// Initialization of gcp client
	auth := new(gcloudAuth)
	auth.JSONPath = i.JSONPath
	// adding scopes since it is custom client, this holds good for all other services
	auth.AuthScopes = []string{compute.CloudPlatformScope}
	client, err := auth.getClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	containerService, err := container.New(client)
	if err != nil {
		log.Fatal(err)
	}

	if i.ProjectID == "" {
		return nil, fmt.Errorf("Project ID cannot be nil")
	}

	if i.Regions == nil {
		i.Regions = append(i.Regions, "-")
	}

	// The parent (project and location) where the clusters will be listed.
	// Specified in the format 'projects/*/locations/*'.
	// Location "-" matches all zones and all regions.
	//parent := "projects/core-search-dev/locations/us-central1" // TODO: Update placeholder value.
	clusters := make([]*container.Cluster, 0)
	for _, r := range i.Regions {
		parent := fmt.Sprintf("projects/%s/locations/%s", i.ProjectID, r)
		resp, err := containerService.Projects.Locations.Clusters.List(parent).Context(ctx).Do()
		if err != nil {
			return nil, err
		}
		for _, cluster := range resp.Clusters {
			clusters = append(clusters, cluster)
		}
	}
	return clusters, nil
}

// GetCluster gets the details of specific kube cluster along with its other details
// This is region specific, mind the region you pass to it.
func (i GetClusterInput) GetCluster() (*container.Cluster, error) {

	// Initialization of gcp client
	auth := new(gcloudAuth)
	auth.JSONPath = i.JSONPath
	// adding scopes since it is custom client, this holds good for all other services
	auth.AuthScopes = []string{compute.CloudPlatformScope}
	client, err := auth.getClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	containerService, err := container.New(client)
	if err != nil {
		log.Fatal(err)
	}

	if i.ProjectID == "" {
		return nil, fmt.Errorf("Project ID cannot be empty")
	}

	if len(i.Regions) == 0 {
		return nil, fmt.Errorf("Region cannot be empty")
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", i.ProjectID, i.Regions[0], i.ClusterName)
	resp, err := containerService.Projects.Locations.Clusters.Get(parent).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
	return resp, nil
}
