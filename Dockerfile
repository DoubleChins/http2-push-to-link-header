FROM golang:alpine AS builder
WORKDIR /app/
COPY . .
RUN go build -o /app/bin/http2-linker

FROM scratch
COPY --from=builder /app/bin/http2-linker /http2-linker
ENTRYPOINT ["/http2-linker"]