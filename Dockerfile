FROM golang:1.22.6 as builder
WORKDIR /app
COPY ./ ./
RUN CGO_ENABLED=0 go build -o sequencer ./cmd/sequencer


# Empty image to minimise footprint and maximise security
FROM scratch

COPY --from=builder /app/sequencer .

ENV VSOCK_PORT 8080

EXPOSE 8080

CMD ["./sequencer"]
