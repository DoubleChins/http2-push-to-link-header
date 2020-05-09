FROM golang:alpine AS builder
WORKDIR /go/
COPY . .
RUN go build -o /go/bin/http2-linker

FROM scratch
COPY --from=builder /go/bin/http2-linker /http2-linker
ENTRYPOINT ["/http2-linker"]