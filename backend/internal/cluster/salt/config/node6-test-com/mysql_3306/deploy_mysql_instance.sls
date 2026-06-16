include:
  - init.files.mysql_dep_pkg
  - user.files.mysql
{% set mysql_version = 8025  %}
{% set toolkits = ['percona-release-latest.noarch.rpm','percona-toolkit-3.3.0-1.el7.x86_64.rpm','percona-xtrabackup-80-8.0.25-17.1.el7.x86_64.rpm'] %}
{% if  mysql_version | int > 8000 %}
{% set xtrabackup_version = mysql_version %}
{% set xtrabackup_tarball = "percona-xtrabackup-8.0.25-17-Linux-x86_64.glibc2.12.tar.gz"  %}
{% else %}
{% set xtrabackup_version = "24" %}
{% set xtrabackup_tarball = "percona-xtrabackup-2.4.22-Linux-x86_64.glibc2.12.tar.gz"  %}
{% endif %}

{% set xtrabackup_unarchive_dir = xtrabackup_tarball |replace(".tar.gz","") %}

{% if mysql_version == 5729 %}
{% set mysql_packet ="Percona-Server-5.7.29-32-Linux.x86_64.ssl101.tar.gz" %}
{% elif mysql_version == 5731 %}
{% set mysql_packet ="Percona-Server-5.7.31-34-Linux.x86_64.glibc2.12.tar.gz" %}
{% elif mysql_version == 5732 %}
{% set mysql_packet ="Percona-Server-5.7.32-35-Linux.x86_64.glibc2.12.tar.gz" %}
{% elif mysql_version == 8025 %}
{% set mysql_packet ="Percona-Server-8.0.25-15-Linux.x86_64.glibc2.12.tar.gz" %}
{% else %}
failure:
  test.fail_without_changes:
    - name: "mysql version is error."
    - failhard: True
{% endif %}

{% set mysql_packet_unarchive_dir =  mysql_packet |replace(".tar.gz","") %}


net.ipv4.ip_local_port_range:
  sysctl.present:
    - value: 10000 65000

# 检查mysql没有了是否为空
{% if not salt['file.directory_exists']('/data/mysql_3306/') %}
init_mysql_dir:
  cmd.run:
    - name: mkdir -p /data/mysql_3306/{data,logs,run,tmp,audit} && chown -R mysql.mysql /data/mysql_3306/
{% else %}
failure:
  test.fail_without_changes:
    - name: "Directory exists"
    - failhard: True
{% endif %}


# 安装percona源
sync_perocna_rpm:
  file.managed:
{% for rpm in toolkits %}
    - name: /tmp/{{rpm}}
    - source: salt://percona/files/{{rpm}}
    - unless: test -f /tmp/{{rpm}}
{% endfor %}


install_percona_rpm:
  cmd.run:
{% for rpm in  toolkits %}
    - name: yum localinstall -y /tmp/{{rpm}}
    - onlyif: test -f /tmp/{{rpm}}
    - require:
      - file: sync_perocna_rpm
{% endfor %}

xtrabackup_install:
  file.managed:
    - name: /usr/local/src/{{xtrabackup_tarball}}
    - source: salt://percona/files/{{xtrabackup_tarball}}
    - unless: test -f /usr/local/src/{{xtrabackup_tarball}}

  cmd.run:
    - name: tar -zxf /usr/local/src/{{xtrabackup_tarball}} -C /usr/local/ && ln -s /usr/local/{{ xtrabackup_unarchive_dir }} /usr/local/xtrabackup{{ xtrabackup_version }}
    - unless: test -d /usr/local/xtrabackup{{ xtrabackup_version }}
    - require:
      - file: xtrabackup_install


mysql-install:
  file.managed:
    - name: /usr/local/src/{{mysql_packet}}
    - source: salt://percona/files/{{mysql_packet}}
    - user: mysql
    - group: mysql
  cmd.run:
    - name: tar -zxf /usr/local/src/{{ mysql_packet }} -C /usr/local/ && ln -s /usr/local/{{ mysql_packet_unarchive_dir }} /usr/local/mysql{{ mysql_version }}
    - unless: test -d /usr/local/mysql{{ mysql_version }} && test -d /usr/local/{{ mysql_packet_unarchive_dir }}
    - require:
      - file: mysql-install

# 创建日志文件
mysql_log_cmd:
  file.managed:
    - name: /data/mysql_3306/logs/mysqld.log
    - user: mysql
    - group: mysql
    - mode: 0755
    - unless: test -f /data/mysql_3306/logs/mysqld.log

mysql_client:
  file.append:
    - name: /etc/profile
    - text: "export PATH=/usr/local/mysql{{ mysql_version }}/bin:$PATH"


mysql_cnf:
  file.managed:
    - name: /data/mysql_3306/my_3306.cnf
    - source: salt://node6-test-com/mysql_3306/my_3306.cnf
    - user: mysql
    - group: mysql
    - mode: 0755

# 生成/etc/my.cnf
client_cnf:
  file.managed:
    - name: /etc/my.cnf
    - source: salt://mysql/files/etc_my.cnf
    - user: mysql
    - group: mysql
    - mode: 0755

async_init_mysql:
  file.managed:
    - name: /tmp/init_mysql_3306.sh
    - source: salt://node6-test-com/mysql_3306/init_mysql_3306.sh
  cmd.run:
    - name: /bin/sh /tmp/init_mysql_3306.sh
    - onlyif: test -z "$(ls -A /data/mysql_3306/data/)"
    - require:
      - file: mysql_cnf
      - cmd: init_mysql_dir

start_mysql:
  cmd.run:
    - name: /usr/local/mysql{{ mysql_version }}/bin/mysqld_safe --defaults-file=/data/mysql_3306/my_3306.cnf &
    - user: mysql
    - bg: True
    - unless: ps -ef|grep -v 'grep' |grep 3306
    - require:
      - cmd: async_init_mysql

update_passwd:
  file.managed:
    - name: /tmp/cuser.sql
    - source: salt://mysql/files/cuser.sql
  cmd.run:
    - name: sleep 15 && /usr/local/mysql{{ mysql_version }}/bin/mysql -uroot -p'' --connect-expired-password -S /data/mysql_3306/run/mysql.sock -e "set global read_only = 0; set session sql_log_bin = 0; ALTER USER 'root'@'localhost' IDENTIFIED BY '123123'; flush privileges;" && /usr/local/mysql{{ mysql_version }}/bin/mysql -uroot -p'123123' -S /data/mysql_3306/run/mysql.sock -e "source /tmp/cuser.sql"
    - unless: test -f /data/mysql_3306/run/mysql.sock
    - require:
      - file: update_passwd

remove_tmp_file:
  cmd.run:
    - name: rm -f /tmp/init_mysql_3306.sh /tmp/cuser.sql
