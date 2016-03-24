package wrap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/client/restclient"
)

//Global Variables
var config = restclient.Config{
	Host: "127.0.0.1:8080",
}

var config2 = restclient.Config{
	Host: "",
}

var imageDeployment = ImageDeployment{
	//repositoryURI: "testURI",
	Repo:         "jbowen",
	Application:  "testapp",
	Revision:     "v0",
	TrafficHosts: []string{},
	PublicPaths:  []string{},
	PublicPort:   "",
	PodCount:     1,
}

func TestCreateDeploymentManager(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)

	//TODO: Better assertion
	assert.NotEmpty(t, deploymentManager)
}

//Testing CreateDeploymentManager with empty config
func TestCreateInDeployment(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config2)
	assert.Nil(t, err)

	assert.NotEmpty(t, deploymentManager)
}

//Namespace Testing
func TestCreateNamespace(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)

	ns, err := deploymentManager.CreateNamespace(imageDeployment)
	assert.Nil(t, err)

	gotNs, err := deploymentManager.client.Namespaces().Get(imageDeployment.Repo)
	assert.Nil(t, err)

	assert.Equal(t, ns, *gotNs)
}

func TestGetNamespace(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)

	ns, err := deploymentManager.CreateNamespace(imageDeployment)
	assert.Nil(t, err)

	gotNs, err := deploymentManager.GetNamespace(imageDeployment)
	assert.Nil(t, err)

	assert.Equal(t, ns, gotNs)
}
func TestDeleteNamespace(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)

	err = deploymentManager.DeleteNamespace(imageDeployment)
	assert.Nil(t, err)
}
func TestCreateandDeleteNamespace(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)

	_, err = deploymentManager.CreateNamespace(imageDeployment)
	assert.Nil(t, err)

	err = deploymentManager.DeleteNamespace(imageDeployment)
	assert.Nil(t, err)

}

func TestConstructDeployment(t *testing.T) {
	template := constructDeployment(imageDeployment)

	assert.NotEmpty(t, template)
	fmt.Println(template)
}

//Deployment Testing

func TestCreateDeployment(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)

	dep, err := deploymentManager.CreateDeployment(imageDeployment)
	assert.Nil(t, err)

	fmt.Println(dep)
}
