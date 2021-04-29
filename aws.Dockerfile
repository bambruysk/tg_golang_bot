

FROM golang:latest
COPY main .
#COPY .env .
COPY client_secret.json .
# COPY key.pem /
ENV TG_API_KEY="1088448942:AAGbDckx7aVCoa005afOE2bVwVejgiPMS4c"

ENV DATABASE_URL="postgres://postgres:s2tj33zm_08@holde-aurora-pg-db.cluster-ctbh4jkhwbkw.eu-central-1.rds.amazonaws.com:5432/holdedb"

EXPOSE 80 443 22 5432

COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
RUN chmod 644 /etc/ssl/certs/ca-certificates.crt && update-ca-certificates

CMD ["/go/main"]