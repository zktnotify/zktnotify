FROM golang:1.12.9 AS go-builder

ENV APP_NAME zktnotify
ENV GOPROXY=https://goproxy.io
ENV GO111MODULE on

WORKDIR /$APP_NAME/builds

COPY go.mod /$APP_NAME/builds
COPY go.sum /$APP_NAME/builds
RUN go mod download

COPY . /$APP_NAME/builds

# build
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
     go build -ldflags="-s -w" -installsuffix cgo -o /$APP_NAME/$APP_NAME /$APP_NAME/builds/main.go \
    && rm -rf /$APP_NAME/builds


FROM centos:8 as prod
MAINTAINER "leaftree <leaftree@github.com>"

WORKDIR /zktnotify

COPY --from=go-builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=go-builder /zktnotify/zktnotify /zktnotify/zktnotify
COPY ./run.sh /zktnotify/run.sh

EXPOSE 4567 4567

CMD ["/zktnotify/zktnotify","start","-f"]
