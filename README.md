# Road to Secure Kubernetes
_Hardening a containerized application one step at a time_

Welcome to 6.0! We're starting get into marginal gains here, but we've still
got some actions we can take. This time we've set resource limits and requests
to limit our disruptive our service can be to other services under a denial of
service attack.

## What has changed?

We've added resource requests and limits to Redis and the web server

```diff
diff --git a/manifests/redis/statefulset.yaml b/manifests/redis/statefulset.yaml
index fb7f22d..a0780df 100644
--- a/manifests/redis/statefulset.yaml
+++ b/manifests/redis/statefulset.yaml
@@ -26,6 +26,13 @@ spec:
         - "$(REDIS_PASSWD)"
         ports:
         - containerPort: 6379
+        resources:
+          requests:
+            memory: 100Mi
+            cpu: 200m
+          limits:
+            memory: 512Mi
+            cpu: 500m
         env:
         - name: REDIS_PASSWD
           valueFrom:
diff --git a/manifests/web/deployment.yaml b/manifests/web/deployment.yaml
index 0e8f31a..120bc7c 100644
--- a/manifests/web/deployment.yaml
+++ b/manifests/web/deployment.yaml
@@ -31,4 +31,11 @@ spec:
               key: password
         ports:
         - containerPort: 8080
+        resources:
+          requests:
+            memory: 20Mi
+            cpu: 100m
+          limits:
+            memory: 100Mi
+            cpu: 300m
```

This is generally considered best practice in Kubernetes, but its not always
highlighted as a security advantage. We'll talk about what this helps prevent
in a second, but first lets talk briefly about how to set these variables.

We can set the requests and limits of CPU and memory. The requests are what
your pod is _guaranteed_ to get and limits are what your pod cannot exceed.

## What does this prevent?

First of all, by setting requests we've made it less likely that the Kubernetes
scheduler will rescheduler or pod under resource contention. Kubernetes kills
pods that use less resources than their requests _last_. It kills pods using
more than their requests first. If you've no requests than you're always in
excess of your requests (Kubernetes treats no requests like requesting zero) so
this is the worst position to be in. This improves availability and that is
one part of the security triangle.

Setting limits on resources limits the damage your service can do to other
services under a denial of service attack. This doesn't stop a denial of
service attack by any means, but it can contain the damage to a single service
in some cases. Without limits, the pod can use CPU or memory until it saturates
the nodes available resources and Kubernetes starts to throttle or kill pods.
In the best case this simply kills the pod that is using all these resources,
but in the worst case it might kill a pod from a different service that forgot
to set requests.

Circuit breaking and rate limiting are much better approaches to limit the
threat of DOS attacks, but this is all we can do with the Dockerfile and
Kubernetes manifests alone.
