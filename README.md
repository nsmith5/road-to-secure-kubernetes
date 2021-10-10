# Road to Secure Kubernetes
_Hardening a containerized application one step at a time_

Welcome to 8.0! One more small, easy change: disable privilege escalation.
 
## What has changed?

Here is the diff:

```diff
diff --git a/manifests/redis/statefulset.yaml b/manifests/redis/statefulset.yaml
index 4dcb900..2346dfc 100644
--- a/manifests/redis/statefulset.yaml
+++ b/manifests/redis/statefulset.yaml
@@ -18,6 +18,7 @@ spec:
         image: redis:latest
         securityContext:
           readOnlyRootFilesystem: true
+          allowPrivilegeEscalation: false
           capabilities:
             drop:
             - ALL
diff --git a/manifests/web/deployment.yaml b/manifests/web/deployment.yaml
index 4254eb7..f7dce8c 100644
--- a/manifests/web/deployment.yaml
+++ b/manifests/web/deployment.yaml
@@ -21,6 +21,7 @@ spec:
         image: nsmith5/road-to-secure-kubernetes:5
         securityContext:
           readOnlyRootFilesystem: true
+          allowPrivilegeEscalation: false
           capabilities:
             drop:
             - ALL
```

## What does this prevent?

Its possible for a process in Linux to create children that have more
privileges than it has. This is a privilege escalation. As an example, consider
a program like `sudo`.

This change sets the `no_new_privs` attribute on our container process which
makes this behaviour impossible.
