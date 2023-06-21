FROM golang:1.20.5-alpine as dev-env

COPY . /go/src/redanthrax/gorlami-server
WORKDIR /go/src/redanthrax/gorlami-server

COPY go.mod ./
#COPY go.sum ./
RUN go mod download

COPY *.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /server

FROM alpine:latest
RUN mkdir /app

COPY --from=dev-env /server ./
CMD ["./server"]