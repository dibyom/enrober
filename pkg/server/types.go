package server

import "k8s.io/kubernetes/pkg/api"

type collectionObject struct {
	selfLink string
	isA      string
	name     string
}

type shipyardGetResponse struct {
	isA          string
	environments string
}

type environmentRequest struct {
	isA           string
	sharingSet    string
	name          string
	publicSecret  bool
	privateSecret bool
	hostNames     []string
	deployments   string
}

type environmentResponse struct {
	isA           string
	sharingSet    string
	name          string
	publicSecret  string
	privateSecret string
	hostNames     []string
	deployments   string
}

type deploymentRequest struct {
	isA             string
	name            string
	deploymentName  string
	publicHosts     string
	privateHosts    string
	replicas        int
	environment     string
	podTemplateSpec *api.PodTemplateSpec
}

type deploymentResponse struct {
	isA             string
	name            string
	deploymentName  string
	publicHosts     string
	privateHosts    string
	replicas        int
	environment     string
	podTemplateSpec *api.PodTemplateSpec
}
