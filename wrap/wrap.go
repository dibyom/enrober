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

type DBStruct struct {
	Name string
	Size string
}

//TODO: May have to add a secret name here?
//ImageDeployment is a collection of necesarry resources for Replication Controller Deployments
type ImageDeployment struct {
	Repo            string
	Application     string
	Revision        string
	TrafficHosts    []string
	PublicPaths     []string
	PathPort        string
	PodCount        int
	Image           string
	ImagePullSecret string
	EnvVars         map[string]string
	//Database stuff
	Database DBStruct //Is this even needed here?
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
	return *ns, nil
}

//GetNamespace <description goes here>
//Returns a Namespace and an error
func (deploymentManager *DeploymentManager) GetNamespace(imageDeployment ImageDeployment) (api.Namespace, error) {
	ns, err := deploymentManager.client.Namespaces().Get(imageDeployment.Repo)
	if err != nil {
		return *ns, err //TODO: Better error handling
	}
	return *ns, nil
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

	//TODO: Maybe dont use this?
	reg := os.Getenv("DOCKER_REGISTRY_URL")
	regString := ""
	if reg != "" {
		regString = reg + "/"
	} else {
		regString = ""
	}

	//Need to make sure the EnvVars map[string]string into a []api.EnvVar
	//TODO: This may be really fucking stupid
	var keys []string
	var values []string
	for k, v := range imageDeployment.EnvVars {
		keys = append(keys, k)
		values = append(values, v)
	}

	envVarTemp := make([]api.EnvVar, len(keys))

	for index, value := range keys {
		envVarTemp[index].Name = value
	}
	for index, value := range values {
		envVarTemp[index].Value = value
	}

	portEnvVar := api.EnvVar{
		Name:  "PORT",
		Value: imageDeployment.PathPort,
	}
	envVarFinal := append(envVarTemp, portEnvVar)

	//TODO: Do we want this
	var imageTemp string
	// fmt.Printf("Image String: %v\n", imageDeployment.Image)
	if imageDeployment.Image == "" { //No passed in Image
		imageTemp = regString + imageDeployment.Repo + "/" + imageDeployment.Application + ":" + imageDeployment.Revision
	} else {
		imageTemp = imageDeployment.Image
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
						//TODO: Make sure this is valid calico policy
						"projectcalico.org/policy": "allow tcp from label Application=" + imageDeployment.Application + " to ports " + imageDeployment.PathPort + "; allow tcp from app=nginx-ingress",
						"trafficHosts":             trafficHosts,
						"publicPaths":              publicPaths,
						"pathPort":                 imageDeployment.PathPort,
					},
				},
				Spec: api.PodSpec{
					//TODO: Come back to this
					ImagePullSecrets: []api.LocalObjectReference{
						api.LocalObjectReference{
							Name: imageDeployment.ImagePullSecret,
						},
					},
					Containers: []api.Container{
						api.Container{
							Name: imageDeployment.Application + "-" + imageDeployment.Revision,
							//TODO: How would we get default images?
							Image: imageTemp,
							Env:   envVarFinal,
							// []api.EnvVar{
							// 	api.EnvVar{
							// 		Name:  "PORT",
							// 		Value: imageDeployment.PathPort,
							// 	},
							// 	//TODO: Add support for passed in env var map
							// 	// for i, v := range imageDeployment
							// },
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
	return *dep, nil
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
	return *depList, nil
}

//CreateDeployment <description goes here>
//Returns a Deployment and an error
func (deploymentManager *DeploymentManager) CreateDeployment(imageDeployment ImageDeployment) (extensions.Deployment, error) {
	template := constructDeployment(imageDeployment)
	dep, err := deploymentManager.client.Deployments(imageDeployment.Repo).Create(&template)
	if err != nil {
		return *dep, err //TODO: Better error handling
	}
	return *dep, nil
}

//UpdateDeployment <description goes here>
func (deploymentManager *DeploymentManager) UpdateDeployment(imageDeployment ImageDeployment) (extensions.Deployment, error) {
	template := constructDeployment(imageDeployment)
	dep, err := deploymentManager.client.Deployments(imageDeployment.Repo).Update(&template)
	if err != nil {
		return *dep, err //TODO: Better error handling
	}
	return *dep, nil
}

// //ConstructSecret <description goes here>
// func (deploymentManager *DeploymentManager) ConstructSecret(name string) (api.Secret, error) {
// 	secretTemplate := api.Secret{
// 		ObjectMeta: api.ObjectMeta{
// 			Name: name,
// 		},
// 		Data: map[string][]byte{},
// 		Type: "Opaque",
// 	}
// 	return secretTemplate, nil
// }

// //CreateSecret <description goes here>
// func (deploymentManager *DeploymentManager) CreateSecret(imageDeployment ImageDeployment) (api.Secret, error) {
// 	template := ConstructSecret() //TODO
// 	secret, err := deploymentManager.client.Secrets(imageDeployment.Repo).Create(&template)
// 	if err != nil {
// 		return *dep, err //TODO: Better error handling
// 	}
// 	return *dep, nil
// }
