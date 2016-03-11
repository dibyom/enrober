//Top Level TODOs go here

//TODO: Make sure all functions have proper description
//TODO: Make sure all functions have proper error handling

package enrober

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	k8sClient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
)

//DeploymentManager is a wrapper type around kubernetes client
type DeploymentManager struct {
	client *k8sClient.Client
}

//ImageDeployment is a collection of necesarry resources for Replication Controller Deployments
type ImageDeployment struct {
	repositoryURI string
	repo          string
	application   string
	revision      string
	virtualHosts  []string
	servedPaths   []string
	podCount      int
}

//CreateDeploymentManager creates an instance of the DeploymentManager from the config passed in, and returns the instance
func CreateDeploymentManager(config restclient.Config) (*DeploymentManager, error) {
	client, err := k8sClient.New(&config)
	if err != nil {
		return nil, err //TODO: Better error handling
	}

	DeploymentManager := &DeploymentManager{
		client: client,
	}
	return DeploymentManager, nil
}

//DeleteReplicationController <description goes here>
//Returns an error
//Delete will only remove a single ReplicationController
func (deploymentManager *DeploymentManager) DeleteReplicationController(imageDeployment ImageDeployment) error {
	err := deploymentManager.client.ReplicationControllers(imageDeployment.repo).Delete(imageDeployment.application)
	if err != nil {
		return err //TODO: Better error handling
	}
	return nil
}

//ListReplicationControllers <description goes here>
//Returns a ReplicationControllerList and an error
//TODO: Should only be passing in label selectors
func (deploymentManager *DeploymentManager) ListReplicationControllers(imageDeployment ImageDeployment) (*api.ReplicationControllerList, error) {
	//TODO: If one option isn't passed in then set it to all or none
	//Need labels to be exclusive, fewer labels == more results
	selector, err := labels.Parse("repo=" + imageDeployment.repo + "," +
		"application=" + imageDeployment.application + "," +
		"revision=" + imageDeployment.revision)
	if err != nil {
		return nil, err //TODO: Better error handling
	}

	options := api.ListOptions{
		LabelSelector: selector,
	}

	controllers, err := deploymentManager.client.ReplicationControllers(imageDeployment.repo).List(options)
	if err != nil {
		return nil, err //TODO: Better error handling
	}

	return controllers, nil
}

//UpdateReplicationController <description goes here>
//Returns a ReplicationController and an error
func (deploymentManager *DeploymentManager) UpdateReplicationController(imageDeployment ImageDeployment) (api.ReplicationController, error) {
	template := ConstructReplicationController(imageDeployment)
	rcResult, err := deploymentManager.client.ReplicationControllers(imageDeployment.repo).Update(&template)
	if err != nil {
		return *rcResult, err //TODO: Better error handling
	}

	return *rcResult, nil
}

//CreateReplicationController <description goes here>
//Returns a ReplicationController and an error
func (deploymentManager *DeploymentManager) CreateReplicationController(imageDeployment ImageDeployment) (api.ReplicationController, error) {
	template := ConstructReplicationController(imageDeployment)
	rcResult, err := deploymentManager.client.ReplicationControllers(imageDeployment.repo).Create(&template)
	if err != nil {
		return *rcResult, err //TODO: Better error handling
	}
	return *rcResult, nil
}

//GetReplicationController <description goes here>
//Returns a ReplicationController and an error
func (deploymentManager *DeploymentManager) GetReplicationController(imageDeployment ImageDeployment) (api.ReplicationController, error) {
	rc, err := deploymentManager.client.ReplicationControllers(imageDeployment.repo).Get(imageDeployment.application)
	if err != nil {
		return *rc, err //TODO: Better error handling
	}
	return *rc, err
}

//CreateNamespace <description goes here>
//Retuns a Namespace and an error
func (deploymentManager *DeploymentManager) CreateNamespace(imageDeployment ImageDeployment) (api.Namespace, error) {
	opt := &api.Namespace{
		ObjectMeta: api.ObjectMeta{
			Name: imageDeployment.repo,
		},
	}
	ns, err := deploymentManager.client.Namespaces().Create(opt)
	if err != nil {
		return *ns, err //TODO: Better error handling
	}
	return *ns, err
}

//TODO: GetNamespace function goes here
//GetNamespace <description goes here>

//DeleteNamespace <description goes here>
func (deploymentManager *DeploymentManager) DeleteNamespace(imageDeployment ImageDeployment) error {
	ns := imageDeployment.repo
	err := deploymentManager.client.Namespaces().Delete(ns)
	if err != nil {
		return err //TODO: Better error handling
	}
	return nil
}

//ConstructReplicationController creates a replication controller object from the passed arguments and default rc template
func ConstructReplicationController(imageDeployment ImageDeployment) api.ReplicationController {
	rcTemplate := api.ReplicationController{
		ObjectMeta: api.ObjectMeta{
			Name: imageDeployment.application, //May take variable
		},
		Spec: api.ReplicationControllerSpec{
			Replicas: imageDeployment.podCount, //Takes imageDeployment.podCount
			Selector: map[string]string{ //ReplicationController Labels go here
				"repo":        imageDeployment.repo,
				"application": imageDeployment.application,
				"revision":    imageDeployment.revision,
			},
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: map[string]string{
						"repo":        imageDeployment.repo,
						"application": imageDeployment.application,
						"revision":    imageDeployment.revision,
					},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Name:  "test1",
							Image: imageDeployment.repo + "/" + imageDeployment.application + ":" + imageDeployment.revision,
						},
					},
				},
			},
		},
		Status: api.ReplicationControllerStatus{},
	}
	return rcTemplate
}
