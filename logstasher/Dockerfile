FROM golang:1.9 as build

RUN mkdir -p /build
WORKDIR /build
COPY funnel.go funnel.go
RUN go build -ldflags "-linkmode external -extldflags -static" -a funnel.go

FROM ubuntu:latest
RUN apt-get update && \
    apt-get -y install curl apt-transport-https && \
    curl https://artifacts.elastic.co/GPG-KEY-elasticsearch | apt-key add - && \
    echo "deb https://artifacts.elastic.co/packages/6.x/apt stable main" | tee -a /etc/apt/sources.list.d/elastic-6.x.list && \
    apt-get update && \
    apt-get -y install default-jre && \
    apt-get -y install logstash && \
    apt-get clean

COPY --from=build /build/funnel /usr/bin/funnel
COPY ./logstash.conf /etc/logstash.conf
COPY run.sh /usr/bin/run.sh

ENTRYPOINT ["run.sh"]
