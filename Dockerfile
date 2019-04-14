ARG GO_VERSION=1.11

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates git gcc build-base
RUN git config --global http.sslVerify false
RUN git config --global --add remote.origin.proxy ""

WORKDIR /app
COPY ./go.mod ./go.sum ./
COPY src src
RUN go get -insecure gopkg.in/yaml.v2
# RUN go get -insecure gopkg.in/check.v1
# RUN go get -insecure github.com/denisenkom/go-mssqldb
# RUN go get -insecure golang.org/x/crypto
# RUN go get -insecure golang.org/x/sys
RUN go get -insecure github.com/sirupsen/logrus
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /ip-geo-locator src/app/main.go

FROM scratch AS final
# COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /ip-geo-locator /ip-geo-locator
ENTRYPOINT ["/ip-geo-locator"]