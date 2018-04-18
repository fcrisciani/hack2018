FROM ubuntu:latest

RUN apt-get update && \
    apt-get -y install git \
    autoconf \
    libtool \
    libnfnetlink-dev \
    libnfnetlink0 \
    libnetfilter-log-dev \
    libnetfilter-conntrack-dev \
    libmnl-dev \
    libnetfilter-acct-dev \
    libjansson-dev \
    curl

WORKDIR /src
RUN git clone git://git.netfilter.org/ulogd2
RUN cd ulogd2 && ./autogen.sh
RUN cd ulogd2 && ./configure
RUN cd ulogd2 && make
RUN cd ulogd2 && make install
COPY ./ulogd.conf /usr/local/etc/ulogd.conf
COPY ./run.sh /usr/local/bin/run.sh

RUN curl https://artifacts.elastic.co/GPG-KEY-elasticsearch | apt-key add - && \
    apt-get install apt-transport-https && \
    echo "deb https://artifacts.elastic.co/packages/6.x/apt stable main" | tee -a /etc/apt/sources.list.d/elastic-6.x.list && \
    apt-get update && apt-get -y install default-jre && apt-get -y install logstash && \
    apt-get clean

COPY ./logstash.conf /etc/logstash.conf

ENTRYPOINT ["run.sh"]
