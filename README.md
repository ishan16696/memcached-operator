## Memcached-Operator
- This is an Memcached-Operator (CRD+controller) using [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) framework.
- It provisions Memcached pods replicas and created a service to access them.

## Overview
- **Kubernetes Controller**: <i>A controller is a loop that reads desired state(spec), observed cluster state (status) and external state and the reconciles cluster state and external state with the desired state, writing any observations down (to our own status)</i>.
- **Custom Resource definition(CRD)**: <i>Defines a custom resource which is not available in default kubernetes implementation</i>.
- **Groups,Versions and kinds**:
   - A **API group** is a collection of related API types.
   - We call each API type a **kind**.
- **API YAML**: It consists of Metadata + Spec + Status + List
   - <i>Metadata</i> holds name/namespace etc.
   - <i>Spec</i> holds desired state.
   - <i>Status</i> holds observed state.
   - <i>List</i> holds many objects.

example: 
```
apiVersion: v1
kind: Pod
metadata:
  name: my-app
  namespace: default
  ...
  ...
spec:
  containers:
  - args: [sh]
    image: gcr.io/bowei-gke/udptest
    imagePullPolicy: Always
    name: client
    ...
  dnsPolicy: ClusterFirst
  ...
  ...
status:
  podIP: 10.8.3.11
  ...
  ...
```
 
### How to run
1. Generate CRD and install the CRD
```
$make manifests
$kubectl apply -f config/crd/bases/cache.example.com_memcacheds.yaml
```

2. Apply the sample
```
$kubectl apply -f config/samples/cache_v1_memcached.yaml
```

3. To start the controller 
```
$make run
```
