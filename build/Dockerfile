FROM golang:latest
#RUN localedef -i 0 -c -f UTF-8 -A /usr/share/locale/locale.alias ru_RU.UTF-8
WORKDIR /app

COPY . /app

RUN chmod +x wait-for-it.sh


ENV TG_API_KEY="1088448942:AAGbDckx7aVCoa005afOE2bVwVejgiPMS4c"

# lN4qnhFs9SL3VREzS7Ed -pg aws parole

ENV DATABASE_URL="postgres://postgres:tg_holde_bot@db/holdedb"


EXPOSE 80 443 22 5432

COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
RUN chmod 644 /etc/ssl/certs/ca-certificates.crt && update-ca-certificates


