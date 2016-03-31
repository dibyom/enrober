//TODO: Implement proper HTTP response codes

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/30x/enrober/wrap"
	"github.com/gorilla/mux"

	"k8s.io/kubernetes/pkg/client/restclient"
)

func main() {
	router := mux.NewRouter()

	sub := router.PathPrefix("/beeswax/deploy/api/v1").Subrouter()

	sub.HandleFunc("/{repo}", RepoHandler).Methods("GET")

	sub.HandleFunc("/{repo}/{application}", ApplicationHandler).Methods("GET")

	sub.HandleFunc("/{repo}/{application}/{revision}", RevisionHandler).Methods("GET", "PUT")

	http.ListenAndServe(":9000", router)
}

//TODO: Maybe use a config file?
//Global Variables
var clientconfig = restclient.Config{
	// Host: "127.0.0.1:8080", //Local Testing

	Host: "", //In Cluster Testing
}

//RepoHandler does stuff
func RepoHandler(w http.ResponseWriter, r *http.Request) {

	//get the variable path
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Path: /%s\n", vars["repo"])

	//get the http verb
	verb := r.Method
	fmt.Fprintf(w, "HTTP Verb: %s\n", verb)

	//manager
	dm, err := wrap.CreateDeploymentManager(clientconfig)
	if err != nil {
		fmt.Fprintf(w, "Broke at manager: %v\n", err)
		fmt.Fprintf(w, "In function RepoHandler\n")
		return
	}

	imagedeployment := wrap.ImageDeployment{
		Repo:         vars["repo"],
		Application:  "",
		Revision:     "",
		TrafficHosts: []string{},
		PublicPaths:  []string{},
		PathPort:     "",
		PodCount:     0,
	}

	//Case statement based on http verb
	switch verb {

	case "GET":

		ns, err := dm.GetNamespace(imagedeployment)
		if err != nil {
			fmt.Fprintf(w, "Broke at namespace: %v\n", err)
			fmt.Fprintf(w, "In function RepoHandler\n")
			return
		}
		fmt.Fprintf(w, "Got Namespace %s\n", ns.GetName())

		depList, err := dm.GetDeploymentList(imagedeployment)
		if err != nil {
			fmt.Fprintf(w, "Broke at deployment: %v\n", err)
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

	//get namespace matching vars["repo"]
	imagedeployment := wrap.ImageDeployment{
		Repo:         vars["repo"],
		Application:  vars["application"],
		Revision:     "",
		TrafficHosts: []string{},
		PublicPaths:  []string{},
		PathPort:     "",
		PodCount:     0,
	}

	//manager
	dm, err := wrap.CreateDeploymentManager(clientconfig)
	if err != nil {
		fmt.Fprintf(w, "Broke at manager: %v\n", err)
		fmt.Fprintf(w, "In function ApplicationHandler\n")
		return
	}

	//Case statement based on http verb
	switch verb {

	case "GET":
		depList, err := dm.GetDeploymentList(imagedeployment)
		if err != nil {
			fmt.Fprintf(w, "Broke at deployment: %v\n", err)
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

	imagedeployment := wrap.ImageDeployment{
		Repo:         vars["repo"],
		Application:  vars["application"],
		Revision:     vars["revision"],
		TrafficHosts: []string{},
		PublicPaths:  []string{},
		PathPort:     "",
		PodCount:     1,
		EnvVars:      map[string]string{},
		Database:     wrap.DBStruct{},
	}

	//manager
	dm, err := wrap.CreateDeploymentManager(clientconfig)
	if err != nil {
		fmt.Fprintf(w, "Broke at manager: %v\n", err)
		fmt.Fprintf(w, "In function RevisionHandler\n")
		return
	}

	//Case statement based on http verb
	switch verb {

	case "GET":
		dep, err := dm.GetDeployment(imagedeployment)
		if err != nil {
			fmt.Fprintf(w, "Broke at deployment: %v\n", err)
			fmt.Fprintf(w, "In function RevisionHandler\n")
			return
		}
		fmt.Fprintf(w, "Got Deployment %v\n", dep.GetName())

	case "PUT":
		//Check if namespace already exists
		getNs, err := dm.GetNamespace(imagedeployment)

		//Namespace wasn't found so create it
		if err != nil && getNs.GetName() == "" {
			ns, err := dm.CreateNamespace(imagedeployment)
			if err != nil {
				fmt.Fprintf(w, "Broke at namespace: %v\n", err)
				fmt.Fprintf(w, "In function RevisionHandler\n")
				return
			}
			fmt.Fprintf(w, "Put Namespace %s\n", ns.GetName())
		}

		//TODO: Possibly put types somewhere else
		//Where I'm putting the JSON body
		type deploymentVariables struct {
			PodCount        int               `json:"podCount"`
			Image           string            `json:"image"`
			ImagePullSecret string            `json:"imagePullSecret"`
			TrafficHosts    []string          `json:"trafficHosts"`
			PublicPaths     []string          `json:"publicPaths"`
			PathPort        int               `json:"pathPort"`
			EnvVars         map[string]string `json:"envVars"`
			Database        wrap.DBStruct     `json:"database"`
		}

		//TODO: Probably a horrifying amount of input validation
		decoder := json.NewDecoder(r.Body)
		var t deploymentVariables
		err = decoder.Decode(&t)

		//TrafficHosts can't be empty so fail if it is
		//TODO: Should do this checking before we create the namespace
		if t.TrafficHosts[0] == "" {
			fmt.Fprintf(w, "Traffic Hosts cannot be empty")
			return
		}
		imagedeployment.TrafficHosts = t.TrafficHosts

		imagedeployment.PublicPaths = t.PublicPaths
		imagedeployment.PathPort = strconv.Itoa(t.PathPort)

		//Make sure this works
		imagedeployment.EnvVars = t.EnvVars

		imagedeployment.Image = t.Image

		//Parse the Database stuff passed in
		//Verify that the size is either small, medium, large, huge
		sizeMap := map[string]struct{}{
			"small":  {},
			"medium": {},
			"large":  {},
			"huge":   {},
		}
		_, ok := sizeMap[t.Database.Size]
		if ok {
			fmt.Fprintf(w, "Valid DB Size: %v\n", t.Database.Size)
		} else {
			fmt.Fprintf(w, "Invalid DB Size: %v\n", t.Database.Size)
			return
		}
		//TODO: Come back and review this
		//Database creation stuff or something

		//Check if deployment already exists
		getDep, err := dm.GetDeployment(imagedeployment)

		//Deployment wasn't found so create it
		if err != nil && getDep.GetName() == "" {

			//TODO: Should create a secret here?

			dep, err := dm.CreateDeployment(imagedeployment)
			if err != nil {
				fmt.Fprintf(w, "Broke at deployment: %v\n", err)
				fmt.Fprintf(w, "In function RevisionHandler\n")
				return
			}
			fmt.Fprintf(w, "Put Deployment %s\n", dep.GetName())
		} else {
			//Deployment was found so modify it
			dep, err := dm.UpdateDeployment(imagedeployment)
			if err != nil {
				fmt.Fprintf(w, "Broke at deployment: %v\n", err)
				fmt.Fprintf(w, "In function RevisionHandler\n")
				return
			}
			fmt.Fprintf(w, "Put Deployment %s\n", dep.GetName())
		}
	}
}
