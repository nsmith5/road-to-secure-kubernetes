# Road to Secure Kubernetes
_Hardening a containerized application one step at a time_

Welcome to 5.0! In this step we've hardened the web container by removing
almost _everything_ from it.

## What has changed?

We've modified the Dockerfile for our web server significantly.

```Dockerfile
FROM golang:1.16 as build
WORKDIR /build
COPY . .
RUN go get
RUN CGO_ENABLED=0 go build
RUN chgrp 0 road-to-secure-kubernetes && chmod g+X road-to-secure-kubernetes

FROM scratch
COPY --from=build /build/road-to-secure-kubernetes .
USER 5678
ENTRYPOINT ["/road-to-secure-kubernetes"]
```

Lets highlight a few things:

- `FROM golang:1.16 as build` shows we're using a multistage build. The first
stage is simply used to build the binary and the second stage is how the final
container is constructed. The binary is copied from the first stage into the second
- `CGO_ENABLED=0 go build` By turning off CGO, the resulting binary is completely
statically linked with zero dependency on glibc etc.
- `RUN chgrp 0` and `chmod g+X` has the same affect as before: The binary can now
be executed by any user ID.
- `FROM scratch` this final container starts as _completely_ empty

All put together, the final container contains exactly one file: the web server
binary. It can be executed by an user. There is _nothing_ else in the
container.

## What does this prevent?

The RCE exploit is now significantly less of a threat. While it exists, there
aren't any programs other than the web server to run!

```
$ curl http://localhost/rce/?cmd=ls
exec: "ls": executable file not found in $PATH

$ curl http://localhost/rce/?cmd=/road-to-secure-kubernetes
listen tcp :8080: bind: address already in use
```

While creating a statically linked binary with zero dependencies isn't possible
for all programming languages, its certainly possible remove a lot from most
containers. From best to worst here are the image bases you should be using for
containers:

- `scratch` absolutely empty. Gold star for you.
- [`distroless`](https://github.com/GoogleContainerTools/distroless) minimal
  containers _without_ the operating system. Supports Java, Python, NodeJS and
other popular languages. These images are signed by
[cosign](https://github.com/sigstore/cosign), which is also an awesome way to
protect your supply chain.
- `busybox` If you absolutely _need_ a shell or other simple CLI tools this is
  the image to use. Thankfully no really bad things like a package manager in
this one.
- `alpine` Minimal linux. Unfortunately this one has a full blown package
  manager. Convenient for making a container but also convenient for exploiting
one.
- `ubuntu / fedora / debian` Absolute worse case scenario. These have a full
  package manager and come bloated with all kinds of tools that help an
attacker.
