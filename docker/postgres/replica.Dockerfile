FROM postgres:18.1-alpine3.23

COPY replica.entrypoint.sh /usr/local/bin/replica.entrypoint.sh
RUN chmod +x /usr/local/bin/replica.entrypoint.sh
