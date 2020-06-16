FROM golang:1.14-alpine AS build

WORKDIR /go/src/kubernetes-job/
COPY src/ ./src/
RUN apk add --no-cache git && cd src/ &&  CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w"  -o  /bin/kubernetes-job *.go

FROM scratch
COPY --from=build /bin/kubernetes-job /bin/kubernetes-job

ENTRYPOINT ["/bin/kubernetes-job"]
