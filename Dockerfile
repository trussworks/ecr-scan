FROM lambci/lambda:build-go1.x AS base

WORKDIR /app

FROM base AS modules

COPY go.mod go.sum ./
RUN go mod download
