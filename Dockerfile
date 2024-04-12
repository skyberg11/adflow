FROM golang:1.19 AS build
RUN apt-get update && apt-get install gcc g++ make git

WORKDIR /go/src/backend

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV ARCH=amd64

COPY . .

RUN go mod download

RUN go build -o ./bin/backend ./cmd/main/main.go

FROM ubuntu AS runtime 
WORKDIR /

COPY --from=build /go/src/backend/bin/backend /go/bin/

EXPOSE 8080
ENTRYPOINT /go/bin/backend