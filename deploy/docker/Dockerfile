FROM debian:bullseye

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app
COPY /deploy/bin/pigeomail /app
COPY /deploy/migrations /app/migrations

EXPOSE 21201-21025

ENTRYPOINT ["./pigeomail"]
