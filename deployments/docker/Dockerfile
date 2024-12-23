ARG GOLANG_VERSION=1.22
ARG ALPINE_VERSION=3.20

FROM golang:${GOLANG_VERSION}-alpine AS builder
ENV APP=/app
WORKDIR $APP

# Install app dependencies
COPY go.mod go.sum ./
RUN go mod download

# Make binary files
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN mkdir -p $APP/bin && go build -o $APP/bin/app ./cmd/app


FROM alpine:${ALPINE_VERSION}
ENV APP=/app
WORKDIR $APP

# Install additional tools
RUN apk add curl

# Copy project binaries from builder
ENV PATH="$APP:$PATH"
COPY --from=builder $APP/bin $APP

# Create app user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && chown appuser:appgroup $APP
USER appuser
