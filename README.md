# Road to Secure Kubernetes
_Hardening a containerized application one step at a time_

Welcome to 7.0! Small, easy change on this one: drop linux capabilities.
 
## What has changed?

Here is the diff:

```diff
diff --git a/manifests/redis/statefulset.yaml b/manifests/redis/statefulset.yaml
index a0780df..4dcb900 100644
--- a/manifests/redis/statefulset.yaml
+++ b/manifests/redis/statefulset.yaml
@@ -18,6 +18,9 @@ spec:
         image: redis:latest
         securityContext:
           readOnlyRootFilesystem: true
+          capabilities:
+            drop:
+            - ALL
         args:
         - "/usr/local/bin/redis-server"
         - "--appendonly"
diff --git a/manifests/web/deployment.yaml b/manifests/web/deployment.yaml
index 120bc7c..4254eb7 100644
--- a/manifests/web/deployment.yaml
+++ b/manifests/web/deployment.yaml
@@ -21,6 +21,9 @@ spec:
         image: nsmith5/road-to-secure-kubernetes:5
         securityContext:
           readOnlyRootFilesystem: true
+          capabilities:
+            drop:
+            - ALL
         env:
         - name: REDIS_ADDR
           value: redis:6379
```

Linux capabilities are a per-thread attribute that enumerate the privileged actions
possible on a system. For instance CAP_NET_RAW allows a process to

> * Use RAW and PACKET sockets;
> * bind to any address for transparent proxying.

The root user (or processes with effective user id 0) have all capabilities, but
non-root users may have the ability to use some capabilities. In this change we
drop all capabilities from our containers.

## What does this prevent?

By default containers have quite a few capabilities including the CAP_NET_RAW
capability described above. CAP_NET_RAW has been leveraged in a few container
escape vulnerabilies and isn't needed by our application at all. Generally,
we want to restrict all uneeded privileges.

If your service binds to a port < 1024 you may want to add `CAP_BIND_SERVICE`
or better yet, change it to bind to a unprivileged port!
