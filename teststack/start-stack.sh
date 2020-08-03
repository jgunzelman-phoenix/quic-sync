docker build -t quic-sync-proxy teststack/nginx
docker stack deploy -c teststack/quic-sync-test-stack.yaml quic-sync