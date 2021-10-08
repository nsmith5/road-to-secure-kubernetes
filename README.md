# Road to Secure Kubernetes
_Hardening a containerized application on step at a time_

🎉 Ahoy! Welcome to 1.0! Our container application is ready for release.
Our app is a simple webserver that counts the number of requests it has seen.
To deploy the app simply run

```
$ kubectl apply -f manifests
```

And test it out with a couple requests.

```
$ curl http://localhost/
Hello, World for the 1'th time

$ curl http://localhost/
Hello, World for the 2'th time
```

It's not great with grammar, but it sure works!

## Architecture

Our architecture is a classic two-tier system. A simply Go server for the
frontend and Redis to hold the state.

![Architecture Diagram](assets/arch.png)

The Kubernetes manifests can all be found in the `manifests` directory.

## Scary stuff

Unfortuntately our application has some security vulnerabilities 😱. There
seems to be a remote code execution.

```
$ curl http://localhost/rce/?cmd=ls
Dockerfile
go.mod
go.sum
main.go
road-to-secure-kubernetes
```

There also seems to be a server side request forgery vulnerability.

```
$ curl http://localhost/ssrf/?uri=https%3A%2F%2Ficanhazip.com
54.34.12.391
```

If you look in the `app` folder the vulnerabilities will be pretty obvious, but
in a real code base these issues can exist and be extremely hard to detect with
extensive penetration testing.

Instead of assuming we'll be able to catch issues like this in code, lets see
what we can do to limit their impact _without_ changing the code.
