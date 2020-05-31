# QUIC SYNC

Quic Sync is a microservice that links kafka topics over quic http3 connection.

**Why use this over mirrormaker?**
   -  QUIC!!! Quic is a new application level network protocol that uses UDP as it's transport layer and supports most of TCPs features.  Best of all http3 runs on top of quic so this microservice uses REST to exchange data.
      - https://www.chromium.org/quic
   - For disadvantaged and poorly connected clusters this can be used as an alternative to mirrormaker when connectivity errors pervent mirrormaker from working.
   - Topic Filtering,  While currently not a feature we will be adding message filtering as an option to topic connections.

# Build Instructions

The service coded in go and will need to be cross compiled to the target OS and Arch to execute.  We have provided an example to build the service as a alpine linux service container.

**Prerequisites**
   * Mac OS
      1. Install go with cross compiling enabled and dep the go package manager
         ```bash
         brew install go --with-cc-common
         brew install dep
         ```
   * Windows
      1. Install go with cross compiling enabled
      2. Install dep package manager for go

## Build Executable

Build the executable for alpine linux 
1. Pull dependecies with dep
   ```bash
   dep ensure
   ```
1. Build the service for linux and amd64
   ```bash
   GOOS=linux GOARCH=amd64 go build -o quic-sync-server
   ```
## Build Container

Build container for Quic Sync

1. Build Container
   ```bash
   docker build . -t quic-sync:latest
   ```

# Docker Container Info

Below are the environment variable and volumes you need to mount to run this docker container

## Environment Variables

| Variable Name | Default Value | Description |
| ----------- | ----------- | -------- |
| WEB_PORT        | 8443                      | HTTPS port to bind too.
| TLS_CERT_FILE       | /opt/quic_sync/server.crt | Certificate file for TLS
| TLS_KEY_FILE        | /opt/quic_sync/server.key | Key file for TLS
| KAFKA_BOOTSTRAP | kafka:9092            | Bootstrap servers for Kafka brokers

## Running

You have two options for running the container either by itself or via docker swarm with kafka

* Just the container
   ```bash
   docker run -p 8443:8443 quic-sync:latest
   ```
* Swarm stack
  * Start
     ```bash
     docker stack deploy -c test-stack/quic-sync-test-stack.yaml quic-sync
     ```
  * Stop
    ```bash
    docker stack remove quic-sync
    ```