FROM postgres:14

COPY data/brindle.sql /docker-entrypoint-initdb.d/brindle.sql
COPY data/pg_hba.conf /var/lib/postgresql/data/pg_hba.conf

ENV PGDATA=/var/lib/postgresql/data/send
ENV POSTGRES_DB=send