FROM postgres:latest
#RUN localedef -i 0 -c -f UTF-8 -A /usr/share/locale/locale.alias ru_RU.UTF-8
ENV LANG ru_RU.utf8
WORKDIR /app

COPY main /app
COPY .env /app
COPY client_secret.json /app
COPY key.pem /
ENV TG_API_KEY="1088448942:AAGbDckx7aVCoa005afOE2bVwVejgiPMS4c"

# lN4qnhFs9SL3VREzS7Ed -pg aws parole

ENV DATABASE_URL="postgres://postgres:holde_tg_bot@localhost/holdedb"

EXPOSE 80 443 22 5432
VOLUME holde_pgdata /pgdata
ENV PGDATA=pgdata

# COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# RUN chmod 644 /etc/ssl/certs/ca-certificates.crt && update-ca-certificates

ENV POSTGRES_PASSWORD=tg_holde_bot
#ENV POSTGRES_USER=tg_holde_bot
ENV POSTGRES_DB=holdedb




#CMD ["/app/main"]