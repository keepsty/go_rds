#!/bin/sh

/usr/local/mysql{{ MysqlVersion }}/bin/mysql -uadmin -p[[.AdminPassword]] -h 127.0.0.1 -P6032 <<EOF
load mysql users to runtime;
save mysql users from runtime;
save mysql users to disk;
exit
EOF
