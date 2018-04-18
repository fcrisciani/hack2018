FROM alpine:latest

RUN apk add --update \
    ulogd \
    netcat-openbsd \
    curl \
    bash

COPY ulogd.conf /etc
COPY ulogd2json /usr/bin
COPY run.sh /usr/bin

ENTRYPOINT ["run.sh"]
