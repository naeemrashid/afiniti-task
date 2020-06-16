FROM golang:1.14-alpine AS build

WORKDIR /go/src/palindrome-calculator/
COPY src/ ./src/
RUN apk add --no-cache git && cd src/ &&  CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w"  -o  /bin/palindrome-calculator *.go

FROM scratch
COPY --from=build /bin/palindrome-calculator /bin/palindrome-calculator

ENTRYPOINT ["/bin/palindrome-calculator"]
