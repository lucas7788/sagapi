FROM ubuntu:18.04

WORKDIR /app

COPY ./sagapi /app/

CMD ["/app/sagapi", "--config", "/appconfig/config.json", "--loglevel", "1"]
