<H1>Kubernetes Custom Resource Definition (CRD) Example</H1>

This repository is a fork of https://github.com/martin-helmich/kubernetes-crd-example <br>

The yaml files and the golang codes updated to align with below versions 

* golang (>= 1.20.3) 
* kubernetes (>= 1.26.11)

Table of contents
- [Download](#download)
- [Setup](#setup)
- [(Re) Generate golang source code (Optional)](#re-generate-golang-source-code-optional)
- [Build](#build)
- [Run](#run)
  - [Create crd and Watch](#create-crd-and-watch)
  - [Make change to the CRD object](#make-change-to-the-crd-object)

# Download
```
git clone https://github.com/HaesungSeo/kubernetes-crd-example.git
```

# Setup
```
go mod init github.com/martin-helmich/kubernetes-crd-example
go mod tidy
```

# (Re) Generate golang source code (Optional)
```
cd api/types/v1alpha1/
go generate
```

# Build
```
go build
```

# Run

## Create crd and Watch
Before run the watcher, execute below sequences
```
cd kubernetes
kubectl apply -f crd.yaml
```

Run the watcher!
```
./kubernetes-crd-example
```

## Make change to the CRD object
Open another terminal then run it!
```
cd kubernetes
cat project.yaml | kubectl apply -f-
cat project.yaml | sed 's/replicas:.*/replicas: 2/g' | kubectl apply -f-
cat project.yaml | sed 's/replicas:.*/replicas: 3/g' | kubectl apply -f-
cat project.yaml | sed 's/replicas:.*/replicas: 4/g' | kubectl apply -f-
```

See the watcher's screen dump!
```
$ ./kubernetes-crd-example
2024/04/11 13:53:39 Using kubeconfig file:  /home/rocky/.kube/config
0 projects found:

2024/04/11 13:53:46 ADDED: example-project
2024/04/11 13:53:46   spec.replicas 1
2024/04/11 13:55:36 MODIFIED: example-project
2024/04/11 13:55:36   spec.replicas 2
2024/04/11 13:55:40 MODIFIED: example-project
2024/04/11 13:55:40   spec.replicas 3
2024/04/11 13:55:47 MODIFIED: example-project
2024/04/11 13:55:47   spec.replicas 4
```
