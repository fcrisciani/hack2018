#!/bin/sh
/usr/bin/funnel &
/usr/share/logstash/bin/logstash -f /etc/logstash.conf
