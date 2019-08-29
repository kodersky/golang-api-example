FROM golang:alpine

RUN apk update && apk upgrade && \
    apk --update add git gcc make glide tzdata && \
    go get -u github.com/golang/dep/cmd/dep

RUN apk --no-cache add ca-certificates

RUN apk add --no-cache bash

WORKDIR /go/src/github.com/kodersky/golang-api-example

COPY glide.yaml glide.lock ./
RUN glide install

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/api/main.go

#FROM scratch
#COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
#COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=build /go/src/github.com/kodersky/golang-api-example/main /usr/local/bin/main
#COPY --from=build /go/src/github.com/kodersky/golang-api-example/config.yaml /opt/config.yaml
#COPY --from=build /go/src/github.com/kodersky/golang-api-example/wait-for-it.sh /opt/wait-for-it.sh

RUN chmod +x ./wait-for-it.sh

#COPY --from=build /go/src/github.com/kodersky/golang-api-example/error /var/log/error

EXPOSE 8080

#ENTRYPOINT [ "/go/src/github.com/kodersky/golang-api-example/main" ]