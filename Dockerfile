# ---- Dependencies ----
FROM golang:1.14-alpine AS build
WORKDIR /app
RUN apk add git gcc musl-dev
COPY ./src/go.mod .
COPY ./src/go.sum .
RUN go mod download
COPY ./src .
RUN go build -o main mean-reversion.go

# ---- Release ----
FROM golang:1.14-alpine
WORKDIR /app
COPY --from=build /app/main /app/main

CMD "./main"
