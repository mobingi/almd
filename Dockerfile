FROM golang:1.12-alpine as builder
WORKDIR /src
COPY . .
RUN go build -mod=vendor -o /bin/oceand

FROM alpine as release
COPY --from=builder /bin/oceand /oceand
ENV CONFIG_PATH /etc/oceand/config/config.yaml
ENV BACKEND_URL https://service[dev|qa].mobingi.com/m/ocean/oceand/
ENTRYPOINT ["/oceand"]
