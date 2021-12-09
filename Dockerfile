FROM golang:1.17.2-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

ENV CGO_ENABLED=0
RUN go build -o /tfc-runtasks-server

FROM alpine:3.15
WORKDIR /

COPY --from=build /tfc-runtasks-server /tfc-runtasks-server

EXPOSE 10000

CMD ["/tfc-runtasks-server"]
