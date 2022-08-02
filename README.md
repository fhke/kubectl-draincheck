# Kubernetes drain pre-checker

This tool checks whether pods in a Kubernetes cluster are eligible to be evicted by `kubectl drain`.

## Why?

There are multiple reasons why pods cannot be evicted by `kubectl drain`:

* There are multiple pod disruption budgets acting on a pod
* A pod disruption budget is misconfigured, allowing no disruptions
* Multiple pods in a replicaset/statefulset are crashing, causing the pod disruption budget to prevent evictions
* A pod has no owner references

This can make cluster maintenance difficult as manual intervention is required to drain nodes, particularly in multi-tenant clusters where administrators may not have access to fix misconfigured or broken applications. By using this tool, cluster administrators can identify misconfigured workloads & check whether pods are eligible to be evicted prior to draining nodes.

## How does it work?

Firstly, the tool checks that all selected pods have at least one owner reference, as `kubectl drain` will reject pods that do not have an owner reference. Any pods without owner references are reported.

It then uses the [eviction API](https://kubernetes.io/docs/concepts/scheduling-eviction/api-eviction/) to create an eviction resource with dry-run mode enabled for each selected pod. This is the same mechanism that `kubectl drain` uses to evict pods. If there are any errors blocking the pod from being evicted, these are reported.

## Installation guide

### Local install

To install the tool, run the following command:

```console
$ go install github.com/fhke/kubectl-draincheck
```

### Run in Docker

To run in a Docker container, run the following command:

_n.b. - this assumes that your kubeconfig is in the default location `~/.kube/config`. If it is in a different location, you will need to change the volume mount source._

```console
$ docker run --rm -v ~/.kube/config:/kubeconfig -e KUBECONFIG=/kubeconfig quay.io/fhke97/kubectl-draincheck
```

## Usage

### Check all pods in a cluster

```console
$ kubectl draincheck --all-namespaces
```

### Check all pods in a namespace

```console
$ kubectl draincheck --namespace foo
```

### Check specified pods

```console
$ kubectl draincheck --namespace foo bar-pod baz-pod
```
