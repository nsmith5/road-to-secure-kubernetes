FROM ubuntu
WORKDIR /app
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates openssl golang
COPY . .
RUN go get
RUN go build

# Make application runnable by any user (Openshift compatible)
RUN chgrp -R 0 . && chmod -R g+rwX .

# Switch to arbitrary UID
USER 5678

CMD ["./road-to-secure-kubernetes"]
