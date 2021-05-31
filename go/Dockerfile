FROM golang:1.16
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build .
ENTRYPOINT ["/app/opa_scorecard_exporter","--incluster=true"]
