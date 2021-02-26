FROM golang:1.15 as builder

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build -o /go/bin/apiserver -ldflags '-s -w' cmd/apiserver/main.go
RUN go build -o /go/bin/agent -ldflags '-s -w' cmd/agent/main.go

FROM scratch

COPY --from=builder /go/bin/apiserver /usr/local/bin/apiserver
COPY --from=builder /go/bin/agent /usr/local/bin/agent
