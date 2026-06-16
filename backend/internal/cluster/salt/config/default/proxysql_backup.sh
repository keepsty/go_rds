#!/bin/bash
# __author__ = 'Fan()'
# Date: 2020-03-16

export PATH=/usr/local/mysql/bin:/usr/local/java/bin:/sbin:/bin:/usr/sbin:/usr/bin:$PATH

save_days=60
work_dir=/backup/scripts
list_file=${work_dir}/proxysql.db.list
log_file=/backup/proxysql_backup.log.`date +\\%Y\\%m\\%d`

find /data -maxdepth 2 -type f -name "proxysql.db"|xargs -i ls -l {} > ${list_file}

if [ ! -f "${list_file}" ]; then
    echo "[$(date "+%F %H:%M:%S")] [OK] No proxysql.db, exit" >>${log_file}
    exit 0
fi

proxysql_db_array=($(awk '{print $NF}' ${list_file}))
for ((i=0;i<${#proxysql_db_array[@]};i++))
do
    if [ -f "${proxysql_db_array[i]}" ]
    then
       proxysql_admin_port=$(echo ${proxysql_db_array[i]}|awk -F'proxysql_' '{print $2}'|awk -F'/' '{print $1}')
       /bin/cp ${proxysql_db_array[i]} /backup/proxysql_${proxysql_admin_port}.db.`date +\%Y\%m\%d` && /usr/bin/sqlite3 ${proxysql_db_array[i]} .dump > /backup/proxysql_${proxysql_admin_port}.sql.`date +\%Y\%m\%d`
       if [ $? -eq 0 ]
       then
            echo "[$(date "+%F %H:%M:%S")] [OK] Backup success:${proxysql_db_array[i]}" >>${log_file}
        else
            echo "[$(date "+%F %H:%M:%S")] [ERROR] Backup failed:${proxysql_db_array[i]}" >> ${log_file}
        fi
    else
        echo "[$(date "+%F %H:%M:%S")] [ERROR] ${proxysql_db_array[i]}:No such file." >> ${log_file}
    fi
done

if [ -f "/var/lib/proxysql/proxysql.db" ]; then
    proxysql_admin_port=$(cat /etc/proxysql.cnf | grep "mysql_ifaces"| grep -v '^#'| grep -v '^$'|awk -F':' '{print $2}'|awk -F'\"' '{print $1}')
    /bin/cp /var/lib/proxysql/proxysql.db /backup/proxysql_${proxysql_admin_port}.db.`date +\%Y\%m\%d` && /usr/bin/sqlite3 /var/lib/proxysql/proxysql.db .dump > /backup/proxysql_${proxysql_admin_port}.sql.`date +\%Y\%m\%d`
    if [ $? -eq 0 ]
    then
        echo "[$(date "+%F %H:%M:%S")] [OK] Backup success:/var/lib/proxysql/proxysql.db" >>${log_file}
    else
        echo "[$(date "+%F %H:%M:%S")] [ERROR] Backup failed:/var/lib/proxysql/proxysql.db" >> ${log_file}
    fi
fi
find /backup/* -maxdepth 0 -name 'proxysql*' -type f -mtime +${save_days} -exec rm -f {} \;
