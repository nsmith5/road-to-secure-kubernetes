# Road to Secure Kubernetes
_Hardening a containerized application one step at a time_

Welcome to 9.0! We're starting to hit some niche stuff here: kernel protections
with seccomp.
 
## What has changed?

Back when you created your `kind` cluster at the start of the exercise a
seccomp profile was loaded into each node in the cluster. The seccomp profile
is a simple JSON file (checkout `cluster/seccomp/fine-grain.json`).

We've modified our web server to use this seccomp profile here

```diff
diff --git a/manifests/web/deployment.yaml b/manifests/web/deployment.yaml
index f7dce8c..417e673 100644
--- a/manifests/web/deployment.yaml
+++ b/manifests/web/deployment.yaml
@@ -16,6 +16,9 @@ spec:
         runAsGroup: 4444
         runAsUser: 1234
         runAsNonRoot: true
+        seccompProfile:
+          type: Localhost
+          localhostProfile: profiles/fine-grain.json
       containers:
       - name: road-to-secure-kubernetes
         image: nsmith5/road-to-secure-kubernetes:5
```

This profile was specially created for this workload and had to be uploaded to
every node that the workload could have been scheduled to so the kubelet could
load it. Pretty tricky to maintain, but is it worth it?

## What does this prevent?

Seccomp limits the system calls a process can make to the Linux kernel. By
limiting the system calls a process can make we're reducing the attack surface
of the kernel via this process. The syscall surface is very large (there > 300
syscalls on modern Linux) and processes typically use only a handful of these
calls.

Seccomp profile have a few weak points though:
- Hard to create: Tooling like `strace`, `perf trace -s` and reading audit logs
  help, but its a long manual process to get it right
- Fragile: Minor changes in application / runtime can change the syscalls used
  and break the app if the original seccomp profile is used.
- Difficult to deploy: As you can see these need to deployed to each node
  instead of ship with the application manifests. This process is often owned
by a different team in many organizations. Who can a team own its seccomp
profile? 

A reasonable compromise is use the `RuntimeDefault` profile that ships with 
container runtimes.

```diff
  securityContext:
+   seccompProfile:
+     type: RuntimeNative
```

This isn't least privileged (e.g it allows syscalls your application doesn't
use), but it does drop many obscure and legacy syscalls to reduce the kernel
attack surface a bit.
