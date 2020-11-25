package utils

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	configmapNS string = "razeedeploy"
	configmapRegex string = "cluster-.*-metadata"
)

type ClusterMetadata struct {
	Region string
	Zone string
}

type Cluster struct {
	Name string
	Metadata ClusterMetadata
}

type ClusterManager interface {
	List() ([]Cluster, error)
}

type ClusterManagerImpl struct {
	Client client.Client
}

/* Assuming the configmaps for the clusters have the structure cluster-<name>-metadata
 * These configmaps should at least contain the following data:
 * {
 *		 ClusterName (string)
 *		 Region (string)
 *		 Zone 	(string)
 * }
 */
func (cm *ClusterManagerImpl) List() ([]Cluster, error) {
	var configmapList corev1.ConfigMapList
	ctx := context.Background()
	ns := client.InNamespace(configmapNS)
	if err := cm.Client.List(ctx, &configmapList, ns); err != nil {
		return nil, err
	}
	var clusters []Cluster
	for _, configmap := range configmapList.Items {
		configmapName := configmap.Name
		if match, _ := regexp.MatchString(configmapRegex, configmapName); match {
			cluster := Cluster {
				Name: configmap.Data["ClusterName"],
				Metadata: ClusterMetadata {
					Region: configmap.Data["Region"],
					Zone: configmap.Data["Zone"],
				},
			}
			clusters = append(clusters, cluster)
		}
	}
	return clusters, nil
}

func CreateClusterManagerImpl (client client.Client) ClusterManager {
	return &ClusterManagerImpl{
		Client: client,
	}
}