FROM alpine:3.12.0
#Environment Variables
ENV WEB_PORT=8443
ENV TLS_CERT_FILE=/opt/quic_sync/server.crt
ENV TLS_KEY_FILE=/opt/quic_sync/server.key 
ENV KAFKA_BOOTSTRAP=kafka:9092
#Add group and user
RUN addgroup -S quic-sync
RUN adduser -S -D quic-sync quic-sync
#Add Files
WORKDIR /opt/quic_sync
ADD [--chown=quic-sync:quic-sync] quic-sync-server /opt/quic_sync/
ADD [--chown=quic-sync:quic-sync] default-certs/* /opt/quic_sync/
#Change ownership
RUN chmod -R 755 /opt/quic_sync

ENTRYPOINT [ "/opt/quic_sync/quic-sync-server --web-port=${WEB_PORT} --cert-file=${TLS_CERT_FILE} --key-file=${TLS_KEY_FILE} --kafka-bootstrap=${KAFKA_BOOTSTRAP}" ]