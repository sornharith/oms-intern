FROM golang:alpine as builder

WORKDIR /go/src/app

ENV GO111MODULE=on

RUN go install github.com/cespare/reflex@v0.3.1

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./run .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

#Copy executable from builder
COPY --from=builder /go/src/app/run .

EXPOSE 8080
CMD ["./run"]