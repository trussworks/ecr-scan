FROM golang:1.15.6 AS base

WORKDIR /app

FROM base AS modules

COPY go.mod go.sum ./
RUN go mod download
