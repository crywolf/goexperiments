FROM golang:1.9.3
WORKDIR /go/src/github.com/crywolf/goexperiments/search/
COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build .

FROM scratch
COPY --from=0 /go/src/github.com/crywolf/goexperiments/search/search .
ENTRYPOINT ["/search"]
