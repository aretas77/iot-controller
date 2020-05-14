FROM golang:alpine as builder

# Install bash for in-container debugging
RUN apk update && apk add --no-cache make git gcc musl-dev

WORKDIR /src/github.com/aretas77/iot-controller

# Need to fetch the dependencies
COPY ./go.mod ./go.sum ./

# Remove some dependencies
RUN sed -ri '/go-python3|aretas77\/paho/d' ./go.mod
RUN go mod tidy && go mod download

# Create barebone directories and copy required files
RUN mkdir -p ./web/iotctl ./types ./utils ./cmd/web ./build ./scripts
COPY ./web/iotctl ./web/iotctl/
COPY ./types ./types/
COPY ./cmd/web ./cmd/web/
COPY ./utils ./utils/
COPY ./scripts ./scripts/

# Build the app with Makefile
RUN make --directory=./cmd/web build

FROM alpine
RUN apk add --no-cache bash tzdata
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /src/github.com/aretas77/iot-controller/build/web /app/
COPY --from=builder /src/github.com/aretas77/iot-controller/scripts/wait-for-it.sh /app/
WORKDIR /app
#CMD ["./web", "start"]
