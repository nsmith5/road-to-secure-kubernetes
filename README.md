# Road to Secure Kubernetes
_Hardening a containerized application one step at a time_

Welcome to 3.0! Exploits are kind of annoying to pull off now because the
almost all of the filesystem is read only.

## What has changed?

Here are the important changes since 2.0:

```diff
diff --git a/manifests/deployment.yaml b/manifests/deployment.yaml
index dcde0a3..8ae9045 100644
--- a/manifests/deployment.yaml
+++ b/manifests/deployment.yaml
@@ -19,6 +19,8 @@ spec:
       containers:
       - name: road-to-secure-kubernetes
         image: nsmith5/road-to-secure-kubernetes:2
+        securityContext:
+          readOnlyRootFilesystem: true
         env:
         - name: REDIS_ADDR
           value: redis:6379
diff --git a/manifests/redis-sts.yaml b/manifests/redis-sts.yaml
index 7ede128..fb7f22d 100644
--- a/manifests/redis-sts.yaml
+++ b/manifests/redis-sts.yaml
@@ -16,6 +16,8 @@ spec:
       containers:
       - name: redis
         image: redis:latest
+        securityContext:
+          readOnlyRootFilesystem: true
         args:
         - "/usr/local/bin/redis-server"
         - "--appendonly"
```

That's it! Just a one line on each deployment! Changing the root filesystem to
be read only is easy to do and very effective. Some applications need a
directory for some caching or scratch space. In those cases, you can mount an
`EmptyDir` volume in that directory and the application will be able to write
to it.

## What does this prevent?

Read only filesystems make running exploit code a lot more difficult. With this
change the only writable location remaining in our system is the `/data`
directory in the Redis pod.

Its now impossible to modify binaries under `/bin` to something malicious for
example.  Its also difficult to install exploit tools. For instance, if you
know about the RCE in the web server containers you'd probably love to run
`nmap` to map out the network from the server perspective. Because you're not
root you can no longer simply run `apt-get install nmap`, but you can no longer
download the binary using curl etc either because you can't save it to a file.

This change is probably one of the lowest effort changes you can make to harden
a workload.
