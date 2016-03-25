//Top Level TODOs go here

//TODO: Decide on better naming scheme
//TODO: Make sure all functions have proper description
//TODO: Make sure all functions have proper error handling

package wrap

import (
	"os"
	"strconv"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
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
	//RepositoryURI string
	Repo         string
	Application  string
	Revision     string
	TrafficHosts []string
	PublicPaths  []string
	PathPort     string
	PodCount     int
}

//CreateDeploymentManager creates an instance of the DeploymentManager from the config passed in, and returns the instance
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
	//Concatenate Annotations
	//TODO: This should be reviewed
	trafficHosts := ""
	for index, value := range imageDeployment.TrafficHosts {
		if index != 0 {
			trafficHosts += " " + value
		} else {
			trafficHosts += value
		}
	}

	publicPaths := ""
	for index, value := range imageDeployment.PublicPaths {
		if index != 0 {
			publicPaths += " " + value
		} else {
			publicPaths += value
		}
	}

	//TODO: Handle this error down the line
	intPathPort, err := strconv.Atoi(imageDeployment.PathPort)
	if err != nil {
		return extensions.Deployment{}
	}

	reg := os.Getenv("DOCKER_REGISTRY_URL")
	regString := ""
	if reg != "" {
		regString = reg + "/"
	} else {
		regString = ""
	}

	depTemplate := extensions.Deployment{
		ObjectMeta: api.ObjectMeta{
			Name: imageDeployment.Application + "-" + imageDeployment.Revision, //May take variable
		},
		Spec: extensions.DeploymentSpec{
			Replicas: imageDeployment.PodCount,
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
						"Repo":         imageDeployment.Repo,
						"Application":  imageDeployment.Application,
						"Revision":     imageDeployment.Revision,
						"microservice": "true",
					},
					Annotations: map[string]string{
						//TODO: Make Optional
						"trafficHosts": trafficHosts,
						"publicPaths":  publicPaths,
						"pathPort":     imageDeployment.PathPort,
					},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Name: imageDeployment.Application + "-" + imageDeployment.Revision,
							//TODO: How would we get default images?
							Image: regString + imageDeployment.Repo + "/" + imageDeployment.Application + ":" + imageDeployment.Revision,
							Env: []api.EnvVar{
								api.EnvVar{
									Name:  "PORT",
									Value: imageDeployment.PathPort,
								},
							},
							Ports: []api.ContainerPort{
								api.ContainerPort{
									ContainerPort: intPathPort,
								},
							},
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
