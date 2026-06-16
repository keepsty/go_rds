#!/bin/sh
chown -R mysql.mysql /data/mysql_3306
/usr/local/mysql8025/bin/mysqld --defaults-file=/data/mysql_3306/my_3306.cnf --datadir=/data/mysql_3306/data/ --initialize-insecure  --user=mysql --basedir=/usr/local/mysql8025
chown -R mysql.mysql /data/mysql_3306
