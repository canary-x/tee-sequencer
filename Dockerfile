# How to build and run from this projects' root dir
# docker build -t com.github.canary-x.tee-sequencer .
# docker run -it --rm com.github.canary-x.tee-sequencer
FROM golang:1.22.6 as builder
WORKDIR /app
COPY ./ ./
RUN CGO_ENABLED=0 go build -o sequencer ./cmd/sequencer


# Empty image to minimise footprint and maximise security
FROM scratch

COPY --from=builder /app/sequencer .

ENV VSOCK_PORT 8080

CMD ["./sequencer"]
