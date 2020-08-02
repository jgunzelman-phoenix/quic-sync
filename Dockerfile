FROM alpine:3.12.0
#Environment Variables
ENV HTTPS_PORT=8443
ENV HTTP3_PORT=8383
ENV TLS_CERT_FILE=/opt/quic-sync/config/server-certs/server.crt
ENV TLS_KEY_FILE=/opt/quic-sync/config/server-certs/server.key
ENV KAFKA_BOOTSTRAP=kafka:9092
#Add group and user
RUN addgroup -S quic-sync
RUN adduser -S -D quic-sync quic-sync
#Add Files
WORKDIR /opt/quic-sync
COPY quic-sync-server /opt/quic-sync/
COPY configs/server-certs /opt/quic-sync/config/server-certs
#Change ownership
RUN chmod -R 755 /opt/quic-sync
#ENTRYPOINT ls -l /opt/quic-sync/config
ENTRYPOINT /opt/quic-sync/quic-sync-server --https-port=${HTTPS_PORT} --http3-port=${HTTP3_PORT} --cert-file=${TLS_CERT_FILE} --key-file=${TLS_KEY_FILE} --kafka-bootstrap=${KAFKA_BOOTSTRAP}