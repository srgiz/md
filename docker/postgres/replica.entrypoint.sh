#!/bin/bash
set -e

PGDATA=/var/lib/postgresql/data
#rm -r "$PGDATA"

# Helper: wait for database-master to be ready
wait_for_primary() {
  echo "Waiting for ${PGHOST} to be ready..."
  until pg_isready -h "${PGHOST}" -p 5432 -U postgres >/dev/null 2>&1; do
    sleep 1
  done
  echo "Primary is ready."
}

# If PGDATA is empty, initialize by taking base backup from primary
if [ ! -s "${PGDATA}/PG_VERSION" ]; then
  rm -r "$PGDATA"

  wait_for_primary

  echo "Running pg_basebackup to pull base backup from primary..."
  #export PGPASSWORD="root"
  pg_basebackup -h "${PGHOST}" -D "${PGDATA}" -U "${PGUSER}" -v -P --wal-method=stream

  # Configure the standby to connect to primary
  #echo "primary_conninfo = 'host=${PGHOST} port=5432 user=${PGUSER} password=${PGPASSWORD} application_name=pg-replica'" >> "${PGDATA}/postgresql.conf"
  echo "primary_conninfo = 'host=${PGHOST} port=5432 user=${PGUSER} password=${PGPASSWORD}'" >> "${PGDATA}/postgresql.conf"
  # create standby.signal to tell Postgres to start in standby mode (PG12+)
  touch "${PGDATA}/standby.signal"

  # (Optional) set hot_standby on standby (already default in many images)
  echo "hot_standby = on" >> "${PGDATA}/postgresql.conf"

  echo "Replica base backup and standby config completed."
fi

# Finally exec the original entrypoint to start postgres normally
exec docker-entrypoint.sh postgres
