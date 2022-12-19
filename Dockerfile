# syntax=docker/dockerfile:1

FROM golang:1.17
WORKDIR /Proftaak_slagboom
COPY go.mod .
COPY go.sum .
COPY . .
RUN go mod download
COPY *.go ./
RUN go build -o main main.go
EXPOSE 3000
CMD ["./main"]
