FROM postgres:18.1-alpine3.23

COPY master.initdb.sh /docker-entrypoint-initdb.d/master.initdb.sh
RUN chmod +x /docker-entrypoint-initdb.d/master.initdb.sh
