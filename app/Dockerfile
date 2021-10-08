FROM ubuntu
WORKDIR /app
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates openssl golang
COPY . .
RUN go get
RUN go build
CMD ["./road-to-secure-kubernetes"]
