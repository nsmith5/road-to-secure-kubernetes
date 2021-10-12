# Road to Secure Kubernetes
_Hardening a containerized application one step at a time_

Welcome to 10.0! We've dropped service account credentials
from the container.
 
## What has changed?

By default, all pods have the credentials of their service account mounted
at `/var/run/secrets/kubernetes.io/serviceaccount`. We've removed these
credentials from Redis and the web server.

```diff
diff --git a/manifests/redis/statefulset.yaml b/manifests/redis/statefulset.yaml
index 2346dfc..827fb4a 100644
--- a/manifests/redis/statefulset.yaml
+++ b/manifests/redis/statefulset.yaml
@@ -13,6 +13,7 @@ spec:
       labels:
         app: redis
     spec:
+      automountServiceAccountToken: false
       containers:
       - name: redis
         image: redis:latest
diff --git a/manifests/web/deployment.yaml b/manifests/web/deployment.yaml
index 417e673..d1bb321 100644
--- a/manifests/web/deployment.yaml
+++ b/manifests/web/deployment.yaml
@@ -12,6 +12,7 @@ spec:
       labels:
         app: road-to-secure-kubernetes
     spec:
+      automountServiceAccountToken: false
       securityContext:
         runAsGroup: 4444
         runAsUser: 1234
```

Previously the following was possible, for instance:

```
$ kubectl exec redis-0 -- ls /var/run/secrets/kubernetes.io/seviceaccount/
ca.crt
namespace
token
```

and now, this directory doesn't exist.

## What does this prevent?

These service account credentials allow applications to talk the Kubernetes
API. Applications like the Nginx ingress controller we're running in our
cluster need these to function, but our application doesn't use the Kubernetes
API at all.

While the default service account doesn't have many privileges someone might
accidentally bind privileges to this service account. This becomes a pathway to
privilege escalation in Kubernetes and in the worst case scenario an exploit
can run arbitrary workloads in your cluster. Its best to remove these
credentials all together if your workload doesn't leverage the Kubernetes API.
