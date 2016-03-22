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
	repo:         "jbowen",
	application:  "testapp",
	revision:     "v0",
	trafficHosts: []string{},
	publicPaths:  []string{},
	publicPort:   "",
	podCount:     1,
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

	gotNs, err := deploymentManager.client.Namespaces().Get(imageDeployment.repo)
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

//Replication Controller Testing

func TestConstructReplicationController(t *testing.T) {
	template := ConstructReplicationController(imageDeployment)

	//TODO: Better assertion
	assert.NotEmpty(t, template)
}

func TestCreateReplicationController(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)
	rc, err := deploymentManager.CreateReplicationController(imageDeployment)
	assert.Nil(t, err)

	getRc, err := deploymentManager.GetReplicationController(imageDeployment)
	assert.Nil(t, err)
	assert.Equal(t, rc, getRc)
}
func TestGetReplicationControllers(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)
	rc, err := deploymentManager.GetReplicationController(imageDeployment)
	assert.Nil(t, err)
	assert.NotEmpty(t, rc)

	// fmt.Printf("%v\n", rc)
	imageName := imageDeployment.repo + "/" + imageDeployment.application + ":" + imageDeployment.revision
	for _, element := range rc.Spec.Template.Spec.Containers {
		// fmt.Printf("%v\n", element.Image)
		assert.Equal(t, element.Image, imageName)
	}
}

//TODO: Have modified ImageDeployment for update to verify it works
//Test UpdateReplicationController
func TestUpdateReplicationController(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)

	var imageDeployment2 = ImageDeployment{
		//repositoryURI: "testURI",
		repo:         "jbowen",
		application:  "testapp",
		revision:     "v0",
		trafficHosts: []string{},
		publicPaths:  []string{},
		publicPort:   "",
		podCount:     3,
	}

	rc, err := deploymentManager.UpdateReplicationController(imageDeployment2)
	assert.Nil(t, err)

	getRc, err := deploymentManager.GetReplicationController(imageDeployment2)
	assert.Nil(t, err)
	assert.Equal(t, rc, getRc)
}

func TestDeleteReplicationController(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)

	//TODO: UpdateReplicationController to have 0 replicas

	err = deploymentManager.DeleteReplicationController(imageDeployment)
	assert.Nil(t, err)

}

func TestListReplicationControllers(t *testing.T) {
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)

	rcList, err := deploymentManager.ListReplicationControllers(imageDeployment)
	assert.Nil(t, err)

	//TODO: Better assertion
	//Mock up container structs?
	// fmt.Printf("%v\n", rcList.Items[0].Labels["application"]) //Gets labels
	imageName := imageDeployment.repo + "/" + imageDeployment.application + ":" + imageDeployment.revision

	for _, element := range rcList.Items[0].Spec.Template.Spec.Containers {
		assert.Equal(t, element.Image, imageName)
	}

	//TODO: Figure out actual assert condition for successful test
}

//End to End test
func TestEndtoEnd(t *testing.T) {
	//CreateDeploymentManager
	deploymentManager, err := CreateDeploymentManager(config)
	assert.Nil(t, err)

	//CreateNamespace
	ns, err := deploymentManager.CreateNamespace(imageDeployment)
	assert.Nil(t, err)

	//GetNamespace
	gotNs, err := deploymentManager.GetNamespace(imageDeployment)
	assert.Nil(t, err)

	assert.Equal(t, ns, gotNs)

	//CreateReplicationController
	rc, err := deploymentManager.CreateReplicationController(imageDeployment)
	assert.Nil(t, err)

	//GetReplicationController
	getRc, err := deploymentManager.GetReplicationController(imageDeployment)
	assert.Nil(t, err)
	assert.Equal(t, rc, getRc)

	//TODO: Add UpdateReplicationController
	var imageDeployment2 = ImageDeployment{
		//repositoryURI: "testURI",
		repo:         "jbowen",
		application:  "testapp",
		revision:     "v0",
		trafficHosts: []string{},
		publicPaths:  []string{},
		publicPort:   "",
		podCount:     3,
	}

	rc, err = deploymentManager.UpdateReplicationController(imageDeployment2)
	assert.Nil(t, err)

	//DeleteReplicationController
	err = deploymentManager.DeleteReplicationController(imageDeployment)
	assert.Nil(t, err)

	//GetReplicationController
	getRc, err = deploymentManager.GetReplicationController(imageDeployment)
	assert.Equal(t, err.Error(), "replicationControllers \"testapp\" not found")

	//DeleteNamespace
	err = deploymentManager.DeleteNamespace(imageDeployment)
	assert.Nil(t, err)

	//GetNamespace
	gotNs, err = deploymentManager.GetNamespace(imageDeployment)
	assert.Nil(t, err)
	assert.Equal(t, string(gotNs.Status.Phase), "Terminating")

}
