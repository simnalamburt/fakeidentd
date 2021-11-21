FROM golang:alpine AS builder
WORKDIR /go/src/fakeidentd
COPY . .
RUN CGO_ENABLED=0 go install

FROM alpine
COPY --from=builder /go/bin/fakeidentd /usr/local/bin
CMD ["fakeidentd"]
