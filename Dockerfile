#Build stage
FROM golang:alpine AS build-env
ENV APP hello-world-webbapp
ENV GOOS linux
ENV GOARCH amd64
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -a -o server

# Final
FROM ubuntu
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /src/server .
CMD ["./server"]
