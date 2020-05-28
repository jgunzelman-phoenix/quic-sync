#QUIC SYNC

Quic Sync is a microservice that links kafka topics over quic http3 connection.

**Why use this over mirrormaker?**
   -  QUIC!!! Quic is a new application level network protocol that uses UDP as it's transport layer and supports most of TCPs features.  Best of all http3 runs on top of quic so this microservice uses REST to exchange data.
   - For disadvantaged and poorly connected clusters this can be used as an alternative to mirrormaker when connectivity errors pervent mirrormaker from working.
   - Topic Filtering,  While currently not a feature we will be adding message filtering as an option to topic connections.

# Build Instructions

The service coded in go and will need to be cross compiled to the target OS and Arch to execute.  We have provided an example to build the service as a alpine linux service container.

**Prerequisites**
   * Mac OS
      * Install go with cross compiling enabled
      ```bash
      brew install go --with-cc-common
      ```
   * Windows
      * Install go with cross compiling enabled

## Build executable

Build the executable for alpine linux 

GOOS=OS to build 

1. Build the service for linux and amd64
   ```bash
   GOOS=linux GOARCH=amd64 go build -o quic-sync-server
   ```