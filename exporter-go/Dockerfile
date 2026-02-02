FROM golang:1.25 as build
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build .
# Now copy it into our base image.
FROM gcr.io/distroless/static
COPY --from=build /app/opa_scorecard_exporter /app/opa_scorecard_exporter
CMD ["/app/opa_scorecard_exporter","--incluster=true"]
