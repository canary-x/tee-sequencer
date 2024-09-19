FROM golang:1.22.6 as builder
WORKDIR /app
COPY ./ ./
RUN CGO_ENABLED=0 go build -o sequencer ./cmd/sequencer


# Empty image to minimise footprint and maximise security
FROM scratch

COPY --from=builder /app/sequencer .

ENV VSOCK_PORT 8080
ENV LOG_VSOCK_PORT 8090
# Usually the host is CID 2, but on EC2 running Nitro it's 3
ENV LOG_VSOCK_CID 3

# Do not use port 9000 for anything, as Nitro needs that free to connect to the instance

EXPOSE 8080

CMD ["./sequencer"]
