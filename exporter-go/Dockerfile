FROM golang:1.17.10 as build
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build .
# Now copy it into our base image.
FROM gcr.io/distroless/static
COPY --from=build /app/exporter-go /app/exporter-go
CMD ["/app/exporter-go","--incluster=true"]
