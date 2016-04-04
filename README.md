#enrober

>###Warning: There is currently minimal input validation on the json.

This project consists of a wrapper library around the kubernetes client api as well as an API server that exposes said library. The server can be deployed both locally and as a docker container. 

###Local Deployment

```sh
go build
./enrober
```

The server will be accesible at `localhost:9000/beeswax/deploy/api/v1`

###Docker Deployment

To deploy the server as a docker container you must modify the client configuration code on lines 31 through 35  in `server.go` to `Host: "",` and then build the project as a static binary before building the docker image. 

To build a static binary:

```sh
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o enrober .
```

To build a docker image:

```sh
docker build -t enrober .
```

A prebuilt docker image is available with:
 
```sh
docker pull jbowen/enrober
```

The server will be accesible at `<docker-ip>/beeswax/deploy/api/v1`

##API Design

The current implementation of the API is simple. It allows for a `GET` request to be made at any of the three levels `repo`, `application`, `revision` or a `PUT` request to be made at the `revision` level. 

**Inputs:** 

`{repo}` is a `string`

`{application}` is a `string`

`{revision}` is a `string`

A `json` body with a corresponding header of `Content-Type: application/json` is required for `PUT` requests. The body has the following contents:

```
{
	"PodCount": 			Int,  				//Required 
	"Image":				string, 			//Optional (for now)
	"ImagePullSecret":		string				//Optional
	"TrafficHosts": 		[string],			//Required
	"PublicPaths":  		[string],			//Required 
	"PathPort": 			Int	,				//Required
	"EnvVars": { 			string: string,		//Optional
	    
	    ...,
	    ...
	},
	"Database": {
		"Name": string,
		"Size": "small/medium/large/huge"
	}
}
```

###Usage

> Assuming you are running the server locally

Get all deployments that match a given repository name:

```sh
curl localhost:9000/beeswax/deploy/api/v1/{repo}
```

Get all deployments that match a given repository name and application name:

```sh
curl localhost:9000/beeswax/deploy/api/v1/{repo}/{application}
```

Check if a deployment exists that matches a given repository name, application name, and revision tag:

```sh
curl localhost:9000/beeswax/deploy/api/v1/{repo}/{application}/{revision}
```

Create a deployment in the given repo, with the given application name and revision tag. If a namespace matching the given repo doesn't exist it will be created. 

>Doesn't cover all options, only required ones.

```sh
curl -X PUT -H "Content-Type: application/json" -d '{"PodCount": 1, "TrafficHosts": ["test.k8s.local"], "PublicPaths": ["/app"], "PathPort": 9000}' "http://localhost:9000/beeswax/deploy/api/v1/{repo}/{application}/{revision}"
```

The Image to be deployed is determined by concatenating the path parameters as follows:

```
{repo} + "/" + {application} + ":" + {revision}
```

If an image matching this pattern isn't found then the deployed pods will fail with a `PullImageError` 
