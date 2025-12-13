FROM golang:1.25.5-alpine3.23 AS builder

RUN apk add make git \
    && go install github.com/google/wire/cmd/wire@v0.7.0
    #&& go install github.com/pressly/goose/v3/cmd/goose@v3.24.2

# Build and cache the dependencies
WORKDIR /srv
COPY go.mod go.sum ./

RUN go mod download

# Copy the actual code files and build the application
COPY . .

RUN make build app=http
    #&& make build app=scheduler

# Build the final image
FROM alpine:3.23

ARG USER_GID=1000
ARG USER_UID=1000
RUN apk add make \
  && addgroup -g $USER_GID app \
  && adduser -u $USER_UID -G app -D app

WORKDIR /srv
USER app

#COPY --from=builder /go/bin/goose /usr/local/bin/goose
#COPY migrations migrations
COPY --chown=app:app go.mod Makefile ./
COPY --chown=app:app .env* ./
#COPY web web
COPY --chown=app:app --from=builder /srv/build ./build

#EXPOSE 8080
CMD ["./build/http"]
