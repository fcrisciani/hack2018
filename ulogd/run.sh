#!/bin/sh

sed -ie "s/MYHOSTNAME/"`hostname`"/" /usr/local/etc/ulogd.conf
exec ulogd -v 
