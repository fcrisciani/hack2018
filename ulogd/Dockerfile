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
    libjansson-dev && \
    apt-get clean

WORKDIR /src
RUN git clone git://git.netfilter.org/ulogd2
RUN cd ulogd2 && ./autogen.sh
RUN cd ulogd2 && ./configure
RUN cd ulogd2 && make
RUN cd ulogd2 && make install
COPY ./ulogd.conf /usr/local/etc/ulogd.conf
COPY ./run.sh /usr/local/bin/run.sh
ENTRYPOINT ["run.sh"]
