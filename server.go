package main

import (
	"fmt"
	"net/http"

	"github.com/30x/enrober/wrap"
	"github.com/gorilla/mux"

	"k8s.io/kubernetes/pkg/client/restclient"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/{repo}", RepoHandler).Methods("GET")

	router.HandleFunc("/{repo}/{application}", ApplicationHandler).Methods("GET")

	router.HandleFunc("/{repo}/{application}/{revision}", RevisionHandler).Methods("GET", "PUT")

	http.ListenAndServe(":8080", router)
}

//RepoHandler does stuff
func RepoHandler(w http.ResponseWriter, r *http.Request) {

	//get the variable path
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Path: /%s\n", vars["repo"])

	//get the http verb
	verb := r.Method
	fmt.Fprintf(w, "HTTP Verb: %s\n", verb)

	//blank config so it shold use InClusterConfig
	//TODO: Should have a way to pass client in externally
	clientconfig := restclient.Config{
		Host: "127.0.0.1:8080",
	}

	//manager
	dm, err := wrap.CreateDeploymentManager(clientconfig)
	if err != nil {
		fmt.Fprintf(w, "Shit broke at manager: %v\n", err)
		fmt.Fprintf(w, "In function RepoHandler\n")
		return
	}

	imagedeployment := wrap.ImageDeployment{
		Repo:         vars["repo"],
		Application:  "",
		Revision:     "",
		TrafficHosts: []string{},
		PublicPaths:  []string{},
		PublicPort:   "",
		PodCount:     0,
	}

	//Case statement based on http verb
	switch verb {

	case "GET":

		ns, err := dm.GetNamespace(imagedeployment)
		if err != nil {
			fmt.Fprintf(w, "Shit broke at namespace: %v\n", err)
			fmt.Fprintf(w, "In function RepoHandler\n")
			return
		}
		fmt.Fprintf(w, "Got Namespace %s\n", ns.GetName())

		depList, err := dm.GetDeploymentList(imagedeployment)
		if err != nil {
			fmt.Fprintf(w, "Shit broke at deployment: %v\n", err)
			fmt.Fprintf(w, "In function ApplicationHandler\n")
			return
		}
		for _, dep := range depList.Items {
			fmt.Fprintf(w, "Got Deployment %v\n", dep.GetName())
		}
	}
}

//ApplicationHandler does stuff
func ApplicationHandler(w http.ResponseWriter, r *http.Request) {

	//get the variable path
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Path: /%s\n", vars["repo"])

	//get the http verb
	verb := r.Method
	fmt.Fprintf(w, "HTTP Verb: %s\n", verb)

	//blank config so it shold use InClusterConfig
	//TODO: Should have a way to pass client in externally
	clientconfig := restclient.Config{
		Host: "127.0.0.1:8080",
	}

	//get namespace matching vars["repo"]
	imagedeployment := wrap.ImageDeployment{
		Repo:         vars["repo"],
		Application:  vars["application"],
		Revision:     "",
		TrafficHosts: []string{},
		PublicPaths:  []string{},
		PublicPort:   "",
		PodCount:     0,
	}

	//manager
	dm, err := wrap.CreateDeploymentManager(clientconfig)
	if err != nil {
		fmt.Fprintf(w, "Shit broke at manager: %v\n", err)
		fmt.Fprintf(w, "In function ApplicationHandler\n")
		return
	}

	//Case statement based on http verb
	switch verb {

	case "GET":
		depList, err := dm.GetDeploymentList(imagedeployment)
		if err != nil {
			fmt.Fprintf(w, "Shit broke at deployment: %v\n", err)
			fmt.Fprintf(w, "In function ApplicationHandler\n")
			return
		}
		for _, dep := range depList.Items {
			fmt.Fprintf(w, "Got Deployment %v\n", dep.GetName())
		}
	}
}

//RevisionHandler does stuff
func RevisionHandler(w http.ResponseWriter, r *http.Request) {
	//get the variable path
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Path: /%s\n", vars["repo"])

	//get the http verb
	verb := r.Method
	fmt.Fprintf(w, "HTTP Verb: %s\n", verb)

	//blank config so it shold use InClusterConfig
	//TODO: Should have a way to pass client in externally
	clientconfig := restclient.Config{
		Host: "127.0.0.1:8080",
	}

	imagedeployment := wrap.ImageDeployment{
		Repo:         vars["repo"],
		Application:  vars["application"],
		Revision:     vars["revision"],
		TrafficHosts: []string{},
		PublicPaths:  []string{},
		PublicPort:   "",
		PodCount:     1,
	}

	//manager
	dm, err := wrap.CreateDeploymentManager(clientconfig)
	if err != nil {
		fmt.Fprintf(w, "Shit broke at manager: %v\n", err)
		fmt.Fprintf(w, "In function RevisionHandler\n")
		return
	}

	//Case statement based on http verb
	switch verb {

	case "GET":
		dep, err := dm.GetDeployment(imagedeployment)
		if err != nil {
			fmt.Fprintf(w, "Shit broke at deployment: %v\n", err)
			fmt.Fprintf(w, "In function RevisionHandler\n")
			return
		}
		fmt.Fprintf(w, "Got Deployment %v\n", dep.GetName())

	case "PUT":
		//Check if namespace already exists
		getNs, err := dm.GetNamespace(imagedeployment)

		//Namespace wasn't found so create it
		//TODO: Want to double check the &&
		if err != nil && getNs.GetName() == "" {
			ns, err := dm.CreateNamespace(imagedeployment)
			if err != nil {
				fmt.Fprintf(w, "Shit broke at namespace: %v\n", err)
				fmt.Fprintf(w, "In function RevisionHandler\n")
				return
			}
			fmt.Fprintf(w, "Put Namespace %s\n", ns.GetName())
		}

		//Check if deployment already exists
		getDep, err := dm.GetDeployment(imagedeployment)

		//Deployment wasn't found so create it
		//TODO: Want to double check the &&
		if err != nil && getDep.GetName() == "" {
			dep, err := dm.CreateDeployment(imagedeployment)
			if err != nil {
				fmt.Fprintf(w, "Shit broke at deployment: %v\n", err)
				fmt.Fprintf(w, "In function RevisionHandler\n")
				return
			}
			fmt.Fprintf(w, "Put Deployment %s\n", dep.GetName())
		} else {
			//Deployment was found so modify it
			dep, err := dm.UpdateDeployment(imagedeployment)
			if err != nil {
				fmt.Fprintf(w, "Shit broke at deployment: %v\n", err)
				fmt.Fprintf(w, "In function RevisionHandler\n")
				return
			}
			fmt.Fprintf(w, "Put Deployment %s\n", dep.GetName())
		}
	}
}
