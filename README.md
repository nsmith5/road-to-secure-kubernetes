# Road to Secure Kubernetes
_Hardening a containerized application one step at a time_

This repository hosts a tutorial on security hardening a containerized workload
in Kubernetes. Its a self-guided, hands on guide from the "default" settings we
see in Kubernetes to a relatively well configured workload. The mitigations
described are by no means exhaustive but show a lot of low hanging fruit anyone
can take advantage of to harden a workload.

## Video Walk-through

I recorded a walk-through of this entire tutorial for folks that want a 
video guide: https://www.youtube.com/watch?v=fe_6UZG8Hlo

## Prerequistes

To run through the tutorial you'll need

- [Docker](https://docker.io)
- [`kind`](https://kind.sigs.k8s.io/) to run a Kubernetes cluster on your laptop with Docker
- `kubectl` the Kubernetes CLI to interact with the cluster
- `helm` to install [Cilium](https://cilium.io/) in our cluster

Before you begin, install the `kind` cluster as follows:

```bash
$ cd cluster

# Install kind cluster
$ kind create cluster --config config.yaml

# Install Cilium into kind cluster
$ helm repo add cilium https://helm.cilium.io/
$ helm install cilium cilium/cilium --version 1.9.10 \
   --namespace kube-system \
   --set nodeinit.enabled=true \
   --set kubeProxyReplacement=partial \
   --set hostServices.enabled=false \
   --set externalIPs.enabled=true \
   --set nodePort.enabled=true \
   --set hostPort.enabled=true \
   --set bpf.masquerade=false \
   --set image.pullPolicy=IfNotPresent \
   --set ipam.mode=kubernetes

# Wait to be installed
$ kubectl wait --for=condition=available deployment.apps/cilium-operator -n kube-system

# Install Nginx Ingress controller
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
$ kubectl wait --for=condition=available deployment.apps/ingress-nginx-controller -n ingress-nginx

```

Once you can run `curl http://localhost` and get back a 404 like this one from Nginx, you're ready
to start

```html
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx</center>
</body>
</html>
```

## How-to

The tutorial shows the step by step progression of an application configuration. Each configuration or step
has a corresponding git tag from `1` to `10`. Start at `1` and move from tag to tag. For every change there
is a detailed explaination of whats been changed and what the change mitigates.

- [Step 1](https://github.com/nsmith5/road-to-secure-kubernetes/tree/1) is our
  starting point. If I was to hazard a guess, about 95% of Kubernetes
application are deployed in this state. Its a functioning application with some
vulnerabilities as you'll see.
- [Step 2](https://github.com/nsmith5/road-to-secure-kubernetes/tree/2) uses a non-root user in the container
- [Step 3](https://github.com/nsmith5/road-to-secure-kubernetes/tree/3) leverages read-only filesystems
- [Step 4](https://github.com/nsmith5/road-to-secure-kubernetes/tree/4) adds network policies
- [Step 5](https://github.com/nsmith5/road-to-secure-kubernetes/tree/5) uses a `scratch` container
- [Step 6](https://github.com/nsmith5/road-to-secure-kubernetes/tree/6) adds resource requests and limits
- [Step 7](https://github.com/nsmith5/road-to-secure-kubernetes/tree/7) drops linux capabilities
- [Step 8](https://github.com/nsmith5/road-to-secure-kubernetes/tree/8) disables privilege escalation
- [Step 9](https://github.com/nsmith5/road-to-secure-kubernetes/tree/9) adds seccomp profile
- [Step 10](https://github.com/nsmith5/road-to-secure-kubernetes/tree/10) removes service account credentials

Navigate to each tag to learn more!
