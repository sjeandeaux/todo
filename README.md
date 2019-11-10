# todo

[![CircleCI](https://circleci.com/gh/sjeandeaux/todo.svg?style=svg)](https://circleci.com/gh/sjeandeaux/todo)
[![Coverage Status](https://coveralls.io/repos/github/sjeandeaux/todo/badge.svg?branch=master)](https://coveralls.io/github/sjeandeaux/todo?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sjeandeaux/todo)](https://goreportcard.com/report/github.com/sjeandeaux/todo)
[![](https://images.microbadger.com/badges/image/sjeandeaux/todo.svg)](https://microbadger.com/images/sjeandeaux/todo)

> Manage your todos with a daemon *todod* and a client *todo-cli*.

## Tools and development

### Structure

> The folder `todo-grpc` contains the schema of the microservice. The schema describes a simple CRUD of todos.
> The folder `todod` contains the default implementation of todo management based on mongo.
> The folder `todo-cli` contains the client which calls the daemon **todod**.
> The folder `pkg` contains the source code.
> The file `Dockerfile` use the multi-staging to generate an image docker with the binaries (todod and todo-cli)

### Tools

* make helps to manage common command lines.
* go the go module must be activated GO111MODULE=on.
* docker-compose for local tests.
* docker for the containerization.
* protoc to generate code from protobuff.

### Development

```sh
make help
clean                          clean the target folder
cover-html                     show the coverage in an HTML page
docker-compose-build           builds the application image with docker-compose.
docker-compose-up              spawns the containers.
fmt                            go fmt
generate                       generate the go from protobuf
gocyclo                        check cyclomatic
help                           prints help.
it-test                        integration test
lint                           go lint on packages
misspell                       misspellpackages
test                           test
tools                          download tools
vet                            go vet on packages
```

### Run the unit tests and integration tests

The tests are written in a Behavior Driven Development way and most of the time in TDD.

```sh
make tools #install the requirements for the test
make dependencies #download the go dependencies
make test
make it-test
```

The integration tests are not considered as [short test](https://golang.org/pkg/testing/#hdr-Skipping).

### Client

```bash
#Create
➜  todo-cli git:(develop) ✗ todo-cli create --title=ori --description='12factor apps' --state=NOT_STARTED --tags="job" --reminder=$(date +%s)
INFO[0000] ID:"5dc58ee5d954d9bc69be5523"
➜  todo-cli git:(develop) ✗ todo-cli create --title=ori --description='12factor eventstore' --state=NOT_STARTED --tags="job" --reminder=$(date +%s)
INFO[0000] ID:"5dc58f08d954d9bc69be5524"

#Read
➜  todo-cli git:(develop) ✗ todo-cli read --id 5dc58f08d954d9bc69be5524
INFO[0000] Todo:&{5dc58f08d954d9bc69be5524 ori 12factor eventstore NOT_STARTED [job golang] 1573228462}

#Update
➜  todo-cli git:(develop) ✗ todo-cli update --id 5dc58f08d954d9bc69be5524 --title=ori --description='12factor eventstore' --state=NOT_STARTED --tags="job,golang" --reminder=$(date +%s)
INFO[0000] Updated:true

#Delete
➜  todo-cli git:(develop) ✗ todo-cli delete --id 5dc587ecd954d9bc69be5522                                                                                                   INFO[0000] Deleted:true

#Search
➜  todo-cli git:(develop) ✗ todo-cli search
INFO[0000] Todo:{5dc586b0a766c6e6404f01b5 Read - Challenge - todo Read - Should create a micro service with 12factor DONE [golang 12factor k8s] 1573046180}
INFO[0000] Todo:{5dc58ee5d954d9bc69be5523 ori 12factor apps NOT_STARTED [job] 1573228261}
INFO[0000] Todo:{5dc58f08d954d9bc69be5524 ori 12factor eventstore NOT_STARTED [job] 1573228296}

➜  todo-cli git:(develop) ✗ todo-cli search --pattern '.event.*'
INFO[0000] Todo:{5dc58f08d954d9bc69be5524 ori 12factor eventstore NOT_STARTED [job] 1573228296}
```


## CI/CD

The stack for the CI/CD:
* circle ci
* coveralls
* docker registry

The circle ci uses:
 * CODECOV_TOKEN
 * COVERALLS_REPO_TOKEN
 * DOCKER_LOGIN
 * DOCKER_PASSWORD

## Questions

### Prove how it aligns to 12factor app best practices

|Factor|Why does it fit|
|:----------:|:----------:|
|Codebase|The SCM follows the gitflow pattern. The `master` is for the production, `develop` for the staging. |
|Dependencies|It uses go mod which manages dependencies. Same happens with helm where we can spevify the dependencies the requirements.yaml file|
|Config|The `todod` application can be configured with environment variables|
|Backing services|The `todod` service has a resource which is the mongo database and configuration with `MONGO_URL`|
|Build, release, run|Build: docker image, Release: kubernetes|
|Processes|The `todod` application is stateless|
|Port binding|The `todod` application use the port 8080|
|Concurrency|The `todod` application is easy to scale|
|Disposability|If the database is up and ready, the `todod` application will start fast |
|Dev/Prod parity|The gitflow pattern allows to have the parity between dev and prod|
|Logs|The logs are redirected in stdout for the container and in a json file.|
|Admin processes|The adminstrotator can use the metrics from `/metrics` and the logs|


### Prove how it fits and uses the best cloud native understanding

* Containerized: The application `todod` is deployed in a container.
* Dynamically orchestrated: A load-balancer (traefik, nginx, ...) can be used above the services of todod.
* Microservices-oriented: The `todod` application manages only todos
* Statelessness: The `todod` can die. Nothing is lost.

### How would you expand on this service to allow for the use of an eventstore?




### How would this service be accessed and used from an external client from the cluster?

Multiple solutions exist : Nodeport, Ingress or a load-balancer.

## Todos

- [ ] better documentation in go
- [ ] better coverage (ex: cli)
- [ ] use a proxy
- [ ] to keep it simple see KO to run go project in kubernetes
- [ ] Kubernates operator instead of helm chart for mongodb


### Useful command lines

## Minikube

```bash
# Install the release
cd helm
helm dep update

export NAME=todo

helm install stable/nginx-ingress --name nginx --namespace dev
helm install --name ${NAME} --namespace dev .


openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /tmp/tls.key -out /tmp/tls.crt -subj "/CN=${NAME}-todod.io"
kubectl create secret tls ${NAME}-todod-secret --key /tmp/tls.key --cert /tmp/tls.crt  -n dev


```




