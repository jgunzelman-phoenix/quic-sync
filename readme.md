# QUIC SYNC

Quic Sync is a microservice that links kafka topics over quic http3 connection.

## System Design
[](./docs/diagrams/standalone.svg)

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
      2. Install Docker Desktop
   * Windows
      1. Install go with cross compiling enabled
      2. Install dep package manager for go
      3. Install Docker Desktop

## Running Project

Build the executable for alpine linux 
1. Pull dependecies with dep
   ```bash
   dep ensure
   ```
### Program arguments

 | Argument | Default Value | Description | Required |
| ----------- | ----------- | -------- | ---- |
| --https-port      | 8443                      | HTTPS port to bind too. | No
| --cert-file       | ./config/server-certs/server.crt | Certificate file for TLS | Yes
| --key-file        | ./config/server-certs/server.key | Key file for TLS | Yes
| --trusted-dir     | ./config/trusted  | Folder for trusted certs to be added
| --kafka-bootstrap | localhost:9092      | Bootstrap servers for Kafka brokers

# Docker Container Info

Here is all you need to know about the quic sync container

## Build Container

Build container for Quic Sync

1. Build the service for linux and amd64 (or what ever arch you are running)
   ```bash
   GOOS=linux GOARCH=amd64 go build -o quic-sync-server
   ```
1. Build Container
   ```bash
   docker build . -t quic-sync:latest
   ```

## Environment Variables

| Variable Name       | Default Value             | Description |
| -----------         | -----------               | -------- |
| HTTPS_PORT          | 8443                      | HTTPS port to bind too.
| HTTP3_PORT          | 8383                      | QUIC HTTP3 port
| TLS_CERT_FILE       | /opt/quic-sync/server.crt | Certificate file for TLS
| TLS_KEY_FILE        | /opt/quic-sync/server.key | Key file for TLS
| KAFKA_BOOTSTRAP     | kafka:9092                | Bootstrap servers for Kafka brokers

## Running

There is a docker swarm stack yaml included that has all the dependecies for quic sync. It also provides a couple of tools to help with testing integration. 

**Running Swarm Stack**
  * Start
     ```bash
     docker stack deploy -c teststack/quic-sync-test-stack.yaml quic-sync
     ```
  * Stop
    ```bash
    docker stack remove quic-sync
    ```
**Service Endpoints**
| Service             | Endpoint                    | 
| -----------         | -----------                 | 
| Kafka UI            | https://localhost/kafka     | 
| Swagger UI          | https://localhost/swagger   | 
| QUIC Sync REST      | https://localhost/quic-sync | 
