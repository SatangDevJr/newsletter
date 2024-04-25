FROM golang:1.18-buster
RUN export GO111MODULE=on
WORKDIR /go/src/app
RUN apt-get update \
    && apt-get install -y --force-yes --no-install-recommends \
    apt-transport-https \
    curl \
    ca-certificates 
COPY go.mod go.sum air.conf ./
COPY apidocs ./apidocs
COPY src ./src
COPY scripts/air/air .
RUN chmod +x air
RUN go mod download
ENV PORT=8000 \
    APP_VERSION=dev \
    DB_NAME= \
    DB_USERNAME= \
    DB_PASSWORD= \
    DB_HOST= \
    DB_PORT= \
    STAGE=
CMD ./air -c air.conf