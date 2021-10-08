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
