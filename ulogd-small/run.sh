#!/bin/sh

#ESADDR=`ip route | grep default | awk '{for(i=1;i<=NF;i++){if($i=="via"){print $(i+1)}}}'`
#if [ "$ESADDR" = "" ] ; then
#	echo "Could not find elastic search address" >&2
#	exit 1
#fi

# Temprorary
ESADDR=52.42.55.249

ulogd -d
while [ ! -f /var/log/ulogd_flow_events.log ] ; do sleep 1 ; done
tail -f /var/log/ulogd_flow_events.log | /usr/bin/ulogd2json | nc -v $ESADDR 55555
