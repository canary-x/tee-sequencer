# How to build and run from this projects' root dir
# docker build -t com.github.canary-x.tee-sequencer .
# docker run -it --rm com.github.canary-x.tee-sequencer
FROM golang:1.22.6 as builder
WORKDIR /app
COPY ./ ./
RUN CGO_ENABLED=0 go build -o sequencer ./cmd/sequencer


FROM gcr.io/distroless/base

# required for Go to operate, as the the home env var is not set in the nitro enclave
ENV HOME /root

COPY --from=builder /app/sequencer .

CMD ["./sequencer"]
