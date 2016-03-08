package enrober

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
)

//TODO: I feel like we're wrapping things that don't need to be wrapped!!!

//DeploymentManager is a wrapper type around kubernetes client
type DeploymentManager struct {
	client *unversioned.Client
}

//ImageDeployment is a collection of necesarry resources for Replication Controller Deployments
type ImageDeployment struct {
	repositoryURI string
	repo          string
	application   string
	revision      string
}

//NewDeploymentManager creates an instance of the DeploymentManager from the config passed in, and returns the instance
func NewDeploymentManager(config *restclient.Config) (*DeploymentManager, error) {
	client, err := unversioned.New(config)
	if err != nil {
		return nil, err
	}

	DeploymentManager := &DeploymentManager{
		client: client,
	}
	return DeploymentManager, nil
}

//DeleteReplicationController <description goes here>
func (deploymentManager *DeploymentManager) DeleteReplicationController(imageDeployment *ImageDeployment) error {
	err := deploymentManager.client.ReplicationControllers(imageDeployment.repo).Delete(imageDeployment.application)
	if err != nil {
		return err
	}
	return nil
}

//GetReplicationControllers <description goes here>
func (deploymentManager *DeploymentManager) GetReplicationControllers(imageDeployment *ImageDeployment) (*api.ReplicationControllerList, error) { //Should also return something else
	//Create selector
	selector, err := labels.Parse("repo = " + imageDeployment.repo + "," +
		"app = " + imageDeployment.application + "," +
		"revision = " + imageDeployment.revision)
	if err != nil {
		return nil, err
	}

	options := api.ListOptions{
		LabelSelector: selector,
	}

	controllers, err := deploymentManager.client.ReplicationControllers(imageDeployment.repo).List(options)
	if err != nil {
		return nil, err
	}

	return controllers, nil
}

//UpdateReplicationController <description goes here>
func (deploymentManager *DeploymentManager) UpdateReplicationController(imageDeployment *ImageDeployment) error { //Maybe should return something else
	return nil
}

//CreateReplicationController <description goes here>
func (deploymentManager *DeploymentManager) CreateReplicationController(imageDeployment *ImageDeployment) error { //Maybe should return something else
	return nil
}

//TODO: Variadic functions to support 3 different GETS

//List All Replication Controllers with matching label.repo

//List All Replication Controllers with matching label.repo and label.app

//List All Replication Controllers with matching label.repo, label.app, label.revision

//Variadic function that takes in the below inputs
// repo uri {string}
// repo name {string}
// image name {string}
// image tag {string}
// virtual hosts[] {array of strings}
// served paths[] {array of strings}
// pod count {int}

//Update pod count of Replication Controller

//Delete Replication Controller
// func DeleteReplicationController(repoURI string)
