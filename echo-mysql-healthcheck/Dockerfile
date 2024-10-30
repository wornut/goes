FROM golang:1.23-alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian12
COPY --from=builder /go/bin/app / 
EXPOSE 8080
CMD ["/app"]
