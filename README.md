# Road to Secure Kubernetes
_Hardening a containerized application one step at a time_

Welcome to 2.0! Things are looking _slightly_ better at this point. No more
root users in our containers!

## What has changed?

Here are the important changes since 1.0:

```diff
diff --git a/app/Dockerfile b/app/Dockerfile
index 6cfb86c..e9c7ac9 100644
--- a/app/Dockerfile
+++ b/app/Dockerfile
@@ -5,4 +5,11 @@ RUN DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates openssl go
 COPY . .
 RUN go get
 RUN go build
+
+# Make application runnable by any user (Openshift compatible)
+RUN chgrp -R 0 . && chmod -R g+rwX .
+
+# Switch to arbitrary UID
+USER 5678
+
 CMD ["./road-to-secure-kubernetes"]
diff --git a/manifests/deployment.yaml b/manifests/deployment.yaml
index fc55e27..dcde0a3 100644
--- a/manifests/deployment.yaml
+++ b/manifests/deployment.yaml
@@ -12,9 +12,13 @@ spec:
       labels:
         app: road-to-secure-kubernetes
     spec:
+      securityContext:
+        runAsGroup: 4444
+        runAsUser: 1234
+        runAsNonRoot: true
       containers:
       - name: road-to-secure-kubernetes
-        image: nsmith5/road-to-secure-kubernetes:1
+        image: nsmith5/road-to-secure-kubernetes:2
         env:
         - name: REDIS_ADDR
           value: redis:6379
```

First, lets concentrate on the Dockerfile. The most important change is that
we've changed the user. Previously it was inherited from the `ubuntu` image and
was `root`.  Now the user will be UID 5678 by default.

More generally, we've made it possible to run our application as an _arbitrary_
user. The line `RUN chgrp -R 0 . && chmod -R g+rwX` changes the permissions on
the application directory such that any user can run the server. This is great
because some platforms [like
Openshift](https://docs.openshift.com/container-platform/3.11/creating_images/guidelines.html#openshift-specific-guidelines)
set a random, non-root user id for each container at run time. This providers
additional protection against container runtime vulnerabilities.

As you can see by the change in our deployment, we run the container as a
totally different unprivileged user (uid 1234) and group (gid 4444) at run time
to emphasize this.

## What does this prevent?

This most obvious change is that you can't run privileged commands in the
container anymore. For instance, if you'd exploited the RCE before and ran

```
# Use RCE to run `apt-get install -y ssh` in the web server container
$ curl http://localhost/rce/?cmd=apt-get%20install%20-y%20ssh
```

this would have installed SSH on one of the web server pods. Now this is the
result

```
$  curl http://localhost/rce/?cmd=apt-get%20install%20-y%20ssh
E: Could not open lock file /var/lib/dpkg/lock-frontend - open (13: Permission denied)
E: Unable to acquire the dpkg frontend lock (/var/lib/dpkg/lock-frontend), are you root?
exit status 100
```

Hazzah!

There is a less obvious win here that isn't always appreciated: Root _inside_
the container is the same as root _outside_ the container. In Kubernetes, users
are not namespaced. This means that the set of user IDs inside of containers
match the users IDs outside of containers. To get a feeling for what this means
compare the output of `ps -aux` inside and outside of the container:

```
# User inside is 1234 and PID is 1
$ kubectl exec <web-pod> -- ps -aux
USER    PID ....... COMMAND
1234    1           ./road-to-secure-kubernetes

# User outside is 1234 and PID is 555802
$ ps -aux | grep road-to-secure-kubernetes
USER    PID     ............  COMMAND
1234    555802                ./road-to-secure-kubernetes
```

These two processes are exactly the same. The container is in something called
a process namespace. This means the PID is different inside and outside of the
container, but notice that the user is the same.

The implication is that if you're root inside a container and manage to find a
container escape exploit (there have been many over time) you become root on
the host!
