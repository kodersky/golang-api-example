FROM golang:alpine

RUN apk update && apk upgrade && \
    apk --update add git gcc make glide tzdata && \
    go get -u github.com/golang/dep/cmd/dep

RUN apk --no-cache add ca-certificates

RUN apk add --no-cache bash

WORKDIR /go/src/github.com/kodersky/golang-api-example

COPY glide.yaml glide.lock ./
RUN glide install

COPY ./wait /wait
RUN chmod +x /wait

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/api/main.go

#FROM scratch
#COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
#COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=build /go/src/github.com/kodersky/golang-api-example/main /usr/local/bin/main
#COPY --from=build /bin/sh /bin/sh
#COPY --from=build /bin/bash /bin/bash
#COPY --from=build /go/src/github.com/kodersky/golang-api-example/config.yaml /opt/config.yaml
#COPY --from=build /usr/local/bin/wait /usr/local/bin/wait
