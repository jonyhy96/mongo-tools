FROM golang:1.13 as builder

ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on

WORKDIR /app/mongo-tools

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make build


FROM alpine
COPY --from=builder /app/mongo-tools/bin /usr/local/bin
