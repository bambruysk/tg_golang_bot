FROM postgres:latest
#ENV LANG ru_RU.utf8

# lN4qnhFs9SL3VREzS7Ed -pg aws parole
EXPOSE 22 5432
VOLUME holde_pgdata_1 /pgdata_1
ENV PGDATA=pgdata_1

# COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# RUN chmod 644 /etc/ssl/certs/ca-certificates.crt && update-ca-certificates

ENV POSTGRES_PASSWORD=tg_holde_bot
#ENV POSTGRES_USER=tg_holde_bot
ENV POSTGRES_DB=holdedb

