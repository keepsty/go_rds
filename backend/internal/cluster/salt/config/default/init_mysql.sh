#!/bin/sh
chown -R mysql.mysql /data/mysql_[[.Port]]
/usr/local/mysql[[.Version]]/bin/mysqld --defaults-file=/data/mysql_[[.Port]]/my_[[.Port]].cnf --datadir=/data/mysql_[[.Port]]/data/ --initialize-insecure  --user=mysql --basedir=/usr/local/mysql[[.Version]]
chown -R mysql.mysql /data/mysql_[[.Port]]
