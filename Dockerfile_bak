FROM golang:1.11.5
ARG version
COPY go.* /go/src/github.com/mobingi/oceand/
COPY *.go /go/src/github.com/mobingi/oceand/
COPY vendor/ /go/src/github.com/mobingi/oceand/vendor/
WORKDIR /go/src/github.com/mobingi/oceand/
RUN GO111MODULE=on GOFLAGS=-mod=vendor CGO_ENABLED=0 GOOS=linux go build -v -ldflags "-X github.com/mobingi/oceand/main.Version=$version" -a -installsuffix cgo -o oceand .

FROM ubuntu:18.04
RUN apt-get update -y && apt-get install -y ca-certificates
WORKDIR /oceand/
COPY --from=0 /go/src/github.com/mobingi/oceand .
ENTRYPOINT ["/oceand/oceand"]
CMD ["run", "--logtostderr"]
