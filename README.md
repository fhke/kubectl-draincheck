# Kubernetes drain pre-checker

This tool checks whether pods in a Kubernetes cluster are eligible to be evicted by `kubectl drain`.

## Why?

There are multiple reasons why pods cannot be evicted by `kubectl drain`:

* There are multiple pod disruption budgets acting on a pod
* A pod disruption budget is misconfigured, allowing no disruptions
* Multiple pods in a replicaset/statefulset are crashing, causing the pod disruption budget to prevent evictions
* A pod has no owner references

This can make cluster maintenance difficult as manual intervention is required to drain nodes, particularly in multi-tenant clusters where administrators may not have access to fix misconfigured or broken applications. By using this tool, cluster administrators can identify misconfigured workloads & check whether pods are able to be evicted prior to draining nodes.

## How do I use it?

To install the tool, run `kubectl install github.com/fhke/kubectl-draincheck`. Once it has been installed, you can use it as follows:

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
