# NOTE it's only a template
FROM golang:1.12-alpine as builder
WORKDIR /src
COPY . .
RUN go build -mod=vendor -o /bin/oceand

FROM alpine as release
COPY --from=builder /bin/oceand /oceand
ENV CONFIG_PATH /xxx/xxx
ENV BACKEND_URL xxx/xxx/xxx
ENTRYPOINT ["/oceand"]