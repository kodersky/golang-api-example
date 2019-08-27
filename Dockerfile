FROM golang:alpine as build
RUN apk update && apk upgrade && \
    apk --update add git gcc make glide tzdata && \
    go get -u github.com/golang/dep/cmd/dep

ENV TZ=Asia/Bangkok

WORKDIR /go/src/github.com/kodersky/golang-api-example

COPY glide.yaml glide.lock ./
RUN glide install

COPY . ./


#RUN go build -o main cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/api/main.go

FROM scratch
COPY --from=build /go/src/github.com/kodersky/golang-api-example/main /usr/local/bin/main
COPY --from=build /go/src/github.com/kodersky/golang-api-example/config.yaml /opt/config.yaml
ENTRYPOINT [ "/usr/local/bin/main" ]

#ENTRYPOINT ["/go/src/github.com/kodersky/golang-api-example/main"]