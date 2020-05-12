FROM golang:alpine as builder

# Install bash for in-container debugging
RUN apk update && \
    apk add --no-cache bash make git gcc musl-dev

WORKDIR /src/github.com/aretas77/iot-controller

# Need to fetch the dependencies
COPY ./go.mod ./go.sum ./

# Remove some dependencies
RUN sed -ri '/go-python3|aretas77\/paho/d' ./go.mod
RUN go mod tidy && go mod download

# Create barebone directories and copy required files
RUN mkdir -p ./web/iotctl ./types ./utils ./cmd/web ./build
COPY ./web/iotctl ./web/iotctl/
COPY ./types ./types/
COPY ./cmd/web ./cmd/web/
COPY ./utils ./utils/

# Build the app with Makefile
RUN make --directory=./cmd/web build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /src/github.com/aretas77/iot-controller/build/web /app/
WORKDIR /app
CMD ["./app", "start"]
