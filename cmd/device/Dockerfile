FROM golang:latest as builder

# Install bash for in-container debugging
RUN apt-get update && apt-get update && apt-get install -y \
    make git gcc musl-dev python3-pip python3.7-dev python3.7
RUN pip3 install tensorflow

WORKDIR /src/github.com/aretas77/iot-controller

# Need to fetch the dependencies
COPY ./go.mod ./go.sum ./
RUN go mod tidy && go mod download

# Create barebone directories and copy required files
RUN mkdir -p ./device ./types ./utils ./build ./scripts ./clients ./cmd/device \
    ./configs ./cmd/data
COPY ./device ./device/
COPY ./configs/device.yaml ./configs/device.yaml
COPY ./types ./types/
COPY ./utils ./utils/
COPY ./scripts ./scripts/
COPY ./clients ./clients/
COPY ./Makefile ./Makefile
COPY ./cmd/device ./cmd/device/
COPY ./cmd/data ./cmd/data/

ENV PYTHONPATH "/bin/python3"

# Build the app with Makefile
RUN make --directory=./cmd/device build
#CMD ["/build/device", "start"]
