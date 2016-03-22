//Top Level TODOs go here

//TODO: Decide on better naming scheme
//TODO: Make sure all functions have proper description
//TODO: Make sure all functions have proper error handling

package wrap

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	k8sClient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
)

//DeploymentManager is a wrapper type around kubernetes client
//TODO: Is this needed?
type DeploymentManager struct {
	client *k8sClient.Client
}

//ImageDeployment is a collection of necesarry resources for Replication Controller Deployments
type ImageDeployment struct {
	//RepositoryURI string
	Repo         string
	Application  string
	Revision     string
	TrafficHosts []string
	PublicPaths  []string
	PublicPort   string
	PodCount     int
}

//CreateDeploymentManager creates an instance of the DeploymentManager from the config passed in, and returns the instance
//TODO: Rename DeploymentManager to KubeManager
func CreateDeploymentManager(config restclient.Config) (*DeploymentManager, error) {
	//Function scoping client
	kubeclient := k8sClient.Client{}

	//No given config so use InClusterConfig
	if config.Host == "" {
		c, err := restclient.InClusterConfig()

		if err != nil {
			return nil, err
		}
		client, err := k8sClient.New(c)

		kubeclient = *client

	} else {
		//Creates client based on passed in config
		client, err := k8sClient.New(&config)

		if err != nil {
			return nil, err
		}

		kubeclient = *client
	}

	//Create the DeploymentManager
	DeploymentManager := &DeploymentManager{
		client: &kubeclient,
	}
	return DeploymentManager, nil
}

//CreateNamespace <description goes here>
//Returns a Namespace and an error
func (deploymentManager *DeploymentManager) CreateNamespace(imageDeployment ImageDeployment) (api.Namespace, error) {
	opt := &api.Namespace{
		ObjectMeta: api.ObjectMeta{
			Name: imageDeployment.Repo,
		},
	}
	ns, err := deploymentManager.client.Namespaces().Create(opt)
	if err != nil {
		return *ns, err //TODO: Better error handling
	}
	return *ns, err
}

//GetNamespace <description goes here>
//Returns a Namespace and an error
func (deploymentManager *DeploymentManager) GetNamespace(imageDeployment ImageDeployment) (api.Namespace, error) {
	ns, err := deploymentManager.client.Namespaces().Get(imageDeployment.Repo)
	if err != nil {
		return *ns, err //TODO: Better error handling
	}
	return *ns, err
}

//DeleteNamespace <description goes here>
//Returns an error
func (deploymentManager *DeploymentManager) DeleteNamespace(imageDeployment ImageDeployment) error {
	ns := imageDeployment.Repo
	err := deploymentManager.client.Namespaces().Delete(ns)
	if err != nil {
		return err //TODO: Better error handling
	}
	return nil
}

//constructDeployment creates a deployment object from the passed arguments and a default deployment template
func constructDeployment(imageDeployment ImageDeployment) extensions.Deployment {
	depTemplate := extensions.Deployment{
		ObjectMeta: api.ObjectMeta{
			Name: imageDeployment.Application + "-" + imageDeployment.Revision, //May take variable
		},
		Spec: extensions.DeploymentSpec{
			Replicas: imageDeployment.PodCount, //Takes imageDeployment.podCount
			Selector: &unversioned.LabelSelector{ //Deployment Labels go here
				MatchLabels: map[string]string{
					"Repo":        imageDeployment.Repo,
					"Application": imageDeployment.Application,
					"Revision":    imageDeployment.Revision,
				},
			},
			Template: api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: map[string]string{
						//TODO: Are these at the Pod or Deployment Level?
						"Repo":         imageDeployment.Repo,
						"Application":  imageDeployment.Application,
						"Revision":     imageDeployment.Revision,
						"microservice": "true",
					},
					Annotations: map[string]string{
						//TODO: Make Optional
						"trafficHosts": "", //TODO: Make concatenated string
						//"trafficHosts": imageDeployment.trafficHosts[0], //TODO: Make concatenated string
						//"publicPaths":  imageDeployment.publicPaths[0],  //TODO: Make concatenated string
						//"publicPort":   imageDeployment.publicPort,
					},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Name:  "test1",
							Image: imageDeployment.Repo + "/" + imageDeployment.Application + ":" + imageDeployment.Revision,
							//ReadinessProbe goes here
							//TODO: Implementation details
							/*
								ReadinessProbe: &api.Probe{
									Handler: api.Handler{
										HTTPGet: &api.HTTPGetAction{
											Path: "/ready", //TODO: This should be determined based on an annotation
											Port: intstr.FromInt(8080),
										},
									},
								},*/
						},
					},
				},
			},
		},
		Status: extensions.DeploymentStatus{},
	}
	return depTemplate
}

//GetDeployment <description goes here>
func (deploymentManager *DeploymentManager) GetDeployment(imageDeployment ImageDeployment) (extensions.Deployment, error) {
	dep, err := deploymentManager.client.Deployments(imageDeployment.Repo).Get(imageDeployment.Application + "-" + imageDeployment.Revision)
	if err != nil {
		return *dep, err //TODO: Better error handling
	}
	return *dep, err
}

//GetDeploymentList <description goes here>
func (deploymentManager *DeploymentManager) GetDeploymentList(imageDeployment ImageDeployment) (extensions.DeploymentList, error) {
	//TODO: Make sure this isn't dumb
	depList := &extensions.DeploymentList{}
	selector, err := labels.Parse("Application=" + imageDeployment.Application)

	//No application is passed
	if imageDeployment.Application == "" {
		depList, err = deploymentManager.client.Deployments(imageDeployment.Repo).List(api.ListOptions{
			LabelSelector: labels.Everything(),
		})
		if err != nil {
			return *depList, err //TODO: Better error handling
		}
	} else {
		depList, err = deploymentManager.client.Deployments(imageDeployment.Repo).List(api.ListOptions{
			LabelSelector: selector,
		})
		if err != nil {
			return *depList, err //TODO: Better error handling
		}
	}
	return *depList, err
}

//CreateDeployment <description goes here>
//Returns a Deployment and an error
func (deploymentManager *DeploymentManager) CreateDeployment(imageDeployment ImageDeployment) (extensions.Deployment, error) {
	template := constructDeployment(imageDeployment)
	dep, err := deploymentManager.client.Deployments(imageDeployment.Repo).Create(&template)
	if err != nil {
		return *dep, err //TODO: Better error handling
	}
	return *dep, err
}

//UpdateDeployment <description goes here>
func (deploymentManager *DeploymentManager) UpdateDeployment(imageDeployment ImageDeployment) (extensions.Deployment, error) {
	template := constructDeployment(imageDeployment)
	dep, err := deploymentManager.client.Deployments(imageDeployment.Repo).Update(&template)
	if err != nil {
		return *dep, err //TODO: Better error handling
	}
	return *dep, err
}

//CreateReplicationController <description goes here>
//Returns a ReplicationController and an error
func (deploymentManager *DeploymentManager) CreateReplicationController(imageDeployment ImageDeployment) (api.ReplicationController, error) {
	template := ConstructReplicationController(imageDeployment)
	rcResult, err := deploymentManager.client.ReplicationControllers(imageDeployment.Repo).Create(&template)
	if err != nil {
		return *rcResult, err //TODO: Better error handling
	}
	return *rcResult, nil
}

//UpdateReplicationController <description goes here>
//Returns a ReplicationController and an error
func (deploymentManager *DeploymentManager) UpdateReplicationController(imageDeployment ImageDeployment) (api.ReplicationController, error) {
	template := ConstructReplicationController(imageDeployment)
	rcResult, err := deploymentManager.client.ReplicationControllers(imageDeployment.Repo).Update(&template)
	if err != nil {
		return *rcResult, err //TODO: Better error handling
	}

	return *rcResult, nil
}

//DeleteReplicationController <description goes here>
//Returns an error
//Delete will only remove a single ReplicationController
func (deploymentManager *DeploymentManager) DeleteReplicationController(imageDeployment ImageDeployment) error {
	err := deploymentManager.client.ReplicationControllers(imageDeployment.Repo).Delete(imageDeployment.Application)
	if err != nil {
		return err //TODO: Better error handling
	}
	return nil
}

//GetReplicationController <description goes here>
//Returns a ReplicationController and an error
func (deploymentManager *DeploymentManager) GetReplicationController(imageDeployment ImageDeployment) (api.ReplicationController, error) {
	rc, err := deploymentManager.client.ReplicationControllers(imageDeployment.Repo).Get(imageDeployment.Application)
	if err != nil {
		return *rc, err //TODO: Better error handling
	}
	return *rc, err
}

//ListReplicationControllers <description goes here>
//Returns a ReplicationControllerList and an error
//TODO: Should only be passing in label selectors
func (deploymentManager *DeploymentManager) ListReplicationControllers(imageDeployment ImageDeployment) (*api.ReplicationControllerList, error) {
	//TODO: If one option isn't passed in then set it to all or none
	//Need labels to be exclusive, fewer labels == more results
	selector, err := labels.Parse("Repo=" + imageDeployment.Repo + "," +
		"Application=" + imageDeployment.Application + "," +
		"Revision=" + imageDeployment.Revision)
	if err != nil {
		return nil, err //TODO: Better error handling
	}

	options := api.ListOptions{
		LabelSelector: selector,
	}

	controllers, err := deploymentManager.client.ReplicationControllers(imageDeployment.Repo).List(options)
	if err != nil {
		return nil, err //TODO: Better error handling
	}

	return controllers, nil
}

//ConstructReplicationController creates a replication controller object from the passed arguments and default rc template
func ConstructReplicationController(imageDeployment ImageDeployment) api.ReplicationController {
	rcTemplate := api.ReplicationController{
		ObjectMeta: api.ObjectMeta{
			Name: imageDeployment.Application, //May take variable
		},
		Spec: api.ReplicationControllerSpec{
			Replicas: imageDeployment.PodCount, //Takes imageDeployment.podCount
			/*
				Selector: map[string]string{ //ReplicationController Labels go here
					"Repo":        imageDeployment.Repo,
					"Application": imageDeployment.Application,
					"revision":    imageDeployment.revision,
				},*/
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: map[string]string{
						"Repo":        imageDeployment.Repo,
						"Application": imageDeployment.Application,
						"Revision":    imageDeployment.Revision,
					},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Name:  "test1",
							Image: imageDeployment.Repo + "/" + imageDeployment.Application + ":" + imageDeployment.Revision,
						},
					},
				},
			},
		},
		Status: api.ReplicationControllerStatus{},
	}
	return rcTemplate
}
