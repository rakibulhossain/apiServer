FROM golang:1.17.0-alpine3.14 AS builder

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .
#https://medium.com/@diogok/on-golang-static-binaries-cross-compiling-and-plugins-1aed33499671
#https://tutorialedge.net/golang/go-multi-stage-docker-tutorial/
RUN CGO_ENABLED=0 GOOS=linux go build -o main .


FROM alpine:latest AS production

COPY --from=builder /app .

EXPOSE 8081

CMD ["./main"]
 


