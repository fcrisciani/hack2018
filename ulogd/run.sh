#!/bin/sh

sed -ie "s/MYHOSTNAME/"`hostname`"/" /usr/local/etc/ulogd.conf
ulogd -d
/usr/share/logstash/bin/logstash -f /etc/logstash.conf
