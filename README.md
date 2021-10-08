# Road to Secure Kubernetes
_Hardening a containerized application one step at a time_

Welcome to 4.0! We've beefed up our network security _significantly_

## What has changed?

We've added network policies now to lock down networking in our namespace.

> Aside: we've also organized the manifests in folders so run `kubectl apply -R
> -f manifests` to update now

At this point only the following traffic is allowed inside the `default`
namespace is the following:

| Source            | Destination      | Reason                            |
|:------------------|:-----------------|:----------------------------------|
| Nginx ingress pod | Web server :8080 | Load balance HTTP requests        |
| Web server        | Core DNS :53     | DNS Lookup for Redis service      |
| Web server        | Redis :6379      | Redis connection for count lookup |

This is acheived with three network policies. The first blocks all ingress
and egress traffic by default in the `default` namespace.

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
```

The second allows the Nginx ingress to web server and the DNS egress and Redis egress.

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: road-to-secure-kubernetes
spec:
  podSelector:
    matchLabels:
      app: road-to-secure-kubernetes
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          kubernetes.io/metadata.name: ingress-nginx
    - podSelector:
        matchLabels:
          app.kubernetes.io/component: controller
          app.kubernetes.io/name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: redis
    ports:
    - protocol: TCP
      port: 6379
  - to:
    - namespaceSelector:
        matchLabels:
          kubernetes.io/metadata.name: kube-system
    - podSelector:
        matchLabels:
          k8s-app: kube-dns
    ports:
    - protocol: UDP
      port: 53
    - protocol: TCP
      port: 53
```

Finally, the third allows ingress from the web server to Redis.

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: redis
spec:
  podSelector:
    matchLabels:
      app: redis
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: road-to-secure-kubernetes
    ports:
    - protocol: TCP
      port: 6379
```

## What does this prevent?

Limiting network like this is extremely effective. Here are a few things you 
can't do anymore

**SSRF Can't reach anything anymore**

If you try to exploit the SSRF endpoint now you'll notice that everything times out.

```
$ curl http://localhost/ssrf/?uri=http://google.com
Get http://google.com: dial tcp 172.217.14.206:80: i/o timeout
```

In a production scenario this might prevent a leak of cloud credentials (the
Capitol One break was SSRF against cloud metadata endpoint), request forgery to
an unintended microservice, privilege escalation inside Kubernetes (the
Kubernetes API is just an HTTP call away from every pod).

**Its difficult to start another server in the web container now**

A naive attacker might install SSH and run an SSH server in our web container
to access the pod. There is only 1 port accessible on the pod and that is 8080,
but its already in use by the web server. If you kill the web server the pod
crashes so its now impossible to add new network listeners to the pod and
connect to them.

**Its difficult to create a reverse shell now**

A more sophisticated attacker might try to open a reverse shell in the web
container. In this scenerio the server is hosted elsewhere and the reverse
shell connects to it as a client. This uses the ephemeral port range and
connects out bound. But now the only outbound connections allowed are to Redis
and CoreDNS.

In a realistic workload it can take some time and research to make a proper
network policy for all microservices in a system. Its easy, for instance, to
forget the DNS looks need to be allowed. Once the work is done though, the
benefits are clear. Minimal networking makes attacks a _lot_ more difficult.
