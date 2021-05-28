FROM golang:1.16
WORKDIR /app
COPY . .
RUN go build .
ENTRYPOINT ["/app/opa_scorecard_exporter","--incluster=true"]
