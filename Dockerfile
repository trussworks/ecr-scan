FROM golang:1.20.1 AS base

WORKDIR /app

FROM base AS modules

COPY go.mod go.sum ./
RUN go mod download
