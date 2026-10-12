# XYZ Corporation — Docker, Kubernetes & CI/CD Projects

A collection of five containerization and orchestration projects demonstrating Docker image optimization, Docker Swarm deployments, Kubernetes network isolation, persistent file uploads, and Jenkins CI/CD automation.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Technology Stack](#technology-stack)
3. [Problem 1 — Optimized Java Docker Image](#problem-1--optimized-java-docker-image)
4. [Problem 2 — Node.js Microservices with Docker Swarm](#problem-2--nodejs-microservices-with-docker-swarm)
5. [Problem 3 — Go Microservices Network Isolation](#problem-3--go-microservices-network-isolation)
6. [Problem 4 — Persistent File Uploads with Flask](#problem-4--persistent-file-uploads-with-flask)
7. [Problem 5 — Jenkins CI/CD Pipeline](#problem-5--jenkins-cicd-pipeline)
8. [Security and Optimization](#security-and-optimization)
9. [Known Limitations](#known-limitations)
10. [Repository Structure](#repository-structure)
11. [Author and Project Links](#author-and-project-links)

## Project Overview

The project demonstrates five practical DevOps scenarios:

| Problem | Objective | Main Technologies |
|---|---|---|
| 1 | Build a small, secure Java container image | Java, Maven, multi-stage Docker build |
| 2 | Deploy and scale a multi-service application | Node.js, Express, PostgreSQL, Docker Swarm |
| 3 | Isolate Go microservices using network policies | Go, Kubernetes, Calico, NetworkPolicy |
| 4 | Accept file uploads and preserve data across container restarts | Python, Flask, Docker Compose, named volumes |
| 5 | Automate testing, image building, publishing, deployment, and verification | Jenkins, Docker Hub, Docker Swarm |

## Technology Stack

- Docker and Docker Compose
- Docker Swarm services, overlay networks, and rolling updates
- Java and Maven
- Node.js and Express
- PostgreSQL
- Go
- Kubernetes and Calico network policies
- Python and Flask
- Jenkins Declarative Pipeline
- Docker Hub
- Git and GitHub
- AWS EC2

## Problem 1 — Optimized Java Docker Image

**Directory:** `java-service/`

### Objective

Build and run a Java microservice using a multi-stage Dockerfile to separate compilation dependencies from the production runtime.

### Implementation

- Uses Maven with Eclipse Temurin Java 17 for the build stage.
- Uses Eclipse Temurin Java 17 JRE for the runtime stage.
- Copies the built application artifact into the runtime image.
- Runs the application as a non-root user with UID `10001`.
- Keeps the runtime image separate from the Maven build environment.

### Build and run

From the repository root:

```bash
cd java-service
docker build -t xyz-java-service:1.0.0 .
docker run --rm xyz-java-service:1.0.0
```

### Verified output

```text
XYZ Corporation Java Microservice is running!
```

The published image was tested and pushed to Docker Hub.

**Published image:** `jayashree0009/xyz-java-service:1.0.0`

## Problem 2 — Node.js Microservices with Docker Swarm

**Directory:** `node-swarm/`

### Objective

Deploy a multi-service application using Docker Swarm, separate frontend and backend traffic, persist database data, and demonstrate scaling and rolling updates.

### Architecture

The stack contains three main services:

- **Web:** NGINX frontend with two replicas.
- **API:** Node.js and Express backend with three replicas.
- **Database:** PostgreSQL with one replica and persistent storage.

### Networking

Two overlay networks separate application traffic:

- `xyz_frontend` — frontend-facing network.
- `xyz_backend` — backend and database network configured as internal.

The separation limits unnecessary network exposure between application tiers.

### Persistence

PostgreSQL uses a named volume called `pgdata`, so database files are stored outside the database container's writable layer.

### Deployment and verification

From the repository root, use the stack configuration supplied in `node-swarm/`:

```bash
cd node-swarm
docker stack deploy -c docker-compose.yml xyz
docker stack services xyz
docker service ls
```

Check individual service tasks:

```bash
docker service ps xyz_web
docker service ps xyz_api
docker service ps xyz_db
```

The stack was verified with two web replicas, three API replicas, and one database replica.

### Scaling and rolling updates

The project demonstrated changing service replica counts and updating a service image. For example, to scale the API:

```bash
docker service scale xyz_api=3
docker service ls
```

Use the actual service image and configuration from the project when performing a rolling update.

### High availability consideration

The lab ran on a single Swarm manager. Multiple service replicas on one node demonstrate service orchestration and scaling, but do not provide node-level high availability. A production cluster should distribute workloads across suitable nodes and use multiple managers with quorum.

## Problem 3 — Go Microservices Network Isolation

**Directory:** `go-network-isolation/`

### Objective

Deploy three Go services in Kubernetes and use NetworkPolicy resources to control which services can communicate.

### Components

- Service A
- Service B
- Service C
- Kubernetes Deployments and Services
- Namespace: `xyz-network-isolation`
- Calico network-policy enforcement

### Implementation

The project includes Kubernetes resources for the application workloads and NetworkPolicy rules:

- `restrict-service-a`
- `restrict-service-b`
- `isolate-service-c`

The policies were tested against the intended allow and deny cases. Actual communication permissions depend on the policy selectors, ports, and ingress or egress rules defined in the manifests.

### Inspect the deployment

```bash
kubectl config current-context
kubectl get namespace xyz-network-isolation
kubectl get deployments,services -n xyz-network-isolation
kubectl get networkpolicy -n xyz-network-isolation
```

Inspect individual policy rules:

```bash
kubectl describe networkpolicy restrict-service-a -n xyz-network-isolation
kubectl describe networkpolicy restrict-service-b -n xyz-network-isolation
kubectl describe networkpolicy isolate-service-c -n xyz-network-isolation
```

### Go tests

From the Go project directory:

```bash
cd go-network-isolation
CGO_ENABLED=0 go test ./...
```

The tests pass with CGO disabled.

### CNI and policy enforcement

The project contains a custom Go CNI-related simulation for learning purposes. It is not a replacement for a production CNI plugin. In the Kubernetes cluster, Calico provides the actual network-policy enforcement.

## Problem 4 — Persistent File Uploads with Flask

**Directory:** `python-uploads/`

### Objective

Build a Flask service that accepts file uploads and stores uploaded files in a Docker named volume so that data survives container restarts.

### API endpoints

| Endpoint | Purpose |
|---|---|
| `POST /upload` | Upload a file |
| `GET /files` | List stored files |

The application stores uploads under `/data` and runs as a non-root user with UID `10001`.

### Start the application

From the repository root:

```bash
cd python-uploads
docker compose up -d --build
docker compose ps
```

The service is published on port `5000`.

### Upload a test file

Create a sample file:

```bash
printf 'XYZ persistent storage test' > sample.txt
```

Upload it:

```bash
curl -X POST \
  -F "file=@sample.txt" \
  http://localhost:5000/upload
```

List uploaded files:

```bash
curl http://localhost:5000/files
```

### Verify persistence

Restart the application container:

```bash
docker compose restart upload-app
```

Then list the files again:

```bash
curl http://localhost:5000/files
```

The test upload was verified to remain available after a container restart because it was stored in a named volume.

The Compose volume is named `uploads_data` in the configuration and is created with a project-scoped name by Docker Compose.

## Problem 5 — Jenkins CI/CD Pipeline

**Directory:** `jenkins-cicd/`

### Objective

Automate the application testing, Docker image build, Docker Hub publishing, Docker Swarm deployment, and deployment verification through Jenkins.

### Pipeline stages

1. **Test** — runs Python unit tests.
2. **Build image** — builds the application container image using the current Jenkins build number as the tag.
3. **Push image** — authenticates to Docker Hub using Jenkins credentials and pushes the image.
4. **Deploy to Swarm** — deploys the stack using the newly built image tag.
5. **Verify deployment** — waits for the expected two replicas and sends an HTTP request to the application.

### Application

The application uses Python's standard-library HTTP server and listens on port `8000`. The Docker image uses Python 3.12 Alpine and runs as a non-root user with UID `10001`.

### Pipeline configuration

- Jenkinsfile: `jenkins-cicd/Jenkinsfile`
- Swarm deployment configuration: `jenkins-cicd/deploy-compose.yml`
- Application: `jenkins-cicd/app/app.py`
- Unit tests: `jenkins-cicd/app/test_app.py`

The pipeline uses the Jenkins credential ID `dockerhub-creds` for Docker Hub authentication. The token is stored in Jenkins credentials rather than hard-coded in the Jenkinsfile.

### Run the pipeline

1. Open the Jenkins web interface.
2. Select the configured pipeline job.
3. Click **Build Now**.
4. Open the build's **Console Output**.
5. Confirm that every stage passes and the build ends with `Finished: SUCCESS`.

### Successful verification

The pipeline was successfully run through all five stages. The final build pushed image tag `6`, deployed the Swarm service, observed the service settle at two replicas, and verified the HTTP response:

```text
XYZ Corporation CI/CD deployment is running!
Deployment verification successful
Finished: SUCCESS
```

The Swarm service is named `xyz-cicd_web` and publishes host port `8082` to container port `8000`.

Useful verification commands:

```bash
docker service ls --filter name=xyz-cicd_web
docker service ps xyz-cicd_web
curl http://127.0.0.1:8082/
```

The Jenkins verification stage uses the EC2 host's private IP to reach the published service from inside the Jenkins container. If the host's private IP changes, update the verification URL in `jenkins-cicd/Jenkinsfile`.

## Security and Optimization

The projects demonstrate the following practices:

- Multi-stage Docker builds to separate build and runtime dependencies.
- Non-root execution for the Java, Flask, and CI/CD application containers.
- Docker named volumes for persistent application data.
- Separate Docker overlay networks for frontend and backend traffic.
- Kubernetes NetworkPolicy resources for workload isolation.
- Jenkins credentials for Docker Hub authentication.
- Versioned CI/CD image tags based on Jenkins build numbers.
- Automated tests and post-deployment HTTP verification.
- Git ignore rules for generated artifacts, local environment files, caches, and test uploads.

### Production hardening

Before production use:

- Restrict Jenkins network access and use HTTPS.
- Protect and rotate credentials and access tokens.
- Avoid exposing the Docker socket to Jenkins unless required. Access to the Docker daemon socket can grant host-root-equivalent control.
- Configure Docker credential storage rather than relying on unencrypted CLI credential files.
- Use a multi-node Swarm cluster for node-level availability.
- Configure monitoring, logging, backups, resource limits, and health checks.
- Review Kubernetes policies for least-privilege ingress and egress rules.

## Known Limitations

- **Swarm availability:** the demonstrated Swarm manager setup uses one node; it is not a highly available manager cluster.
- **Kubernetes policy enforcement:** the custom CNI-related code is a simulation. Calico enforces network policies in the actual cluster.
- **Go test environment:** the documented test command disables CGO; environments requiring CGO may need a C compiler and development headers.
- **Jenkins verification address:** the pipeline currently uses the EC2 private IP `172.31.13.158`. This must be updated if the address changes.
- **Local lab configuration:** port numbers, namespaces, service names, and credentials IDs reflect this project and may need adaptation in other environments.

## Repository Structure

```text
xyz-docker-project/
├── java-service/
│   └── Dockerfile and Java/Maven application
├── node-swarm/
│   └── Node.js, NGINX, PostgreSQL and Swarm configuration
├── go-network-isolation/
│   └── Go services and Kubernetes network-isolation resources
├── python-uploads/
│   └── Flask application and Compose configuration
├── jenkins-cicd/
│   ├── Jenkinsfile
│   ├── deploy-compose.yml
│   └── app/
├── .gitignore
└── README.md
```

## Author and Project Links

### GitHub repository

https://github.com/jayashree3349-ops/xyz-docker-project

### Docker Hub images

- Java microservice: https://hub.docker.com/r/jayashree0009/xyz-java-service
- Go Service A: https://hub.docker.com/r/jayashree0009/xyz-service-a
- Go Service B: https://hub.docker.com/r/jayashree0009/xyz-service-b
- Go Service C: https://hub.docker.com/r/jayashree0009/xyz-service-c
- Jenkins CI/CD application: https://hub.docker.com/r/jayashree0009/xyz-cicd-app
