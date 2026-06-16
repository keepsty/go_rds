include:
  - init.files.mysql_dep_pkg
  - user.files.mysql

{% set mysql_version=8025 %}
{% set basedir="/usr/local/mysql8025" %}
{% set percona_toolkit="percona-toolkit-3.3.0-1.el7.x86_64.rpm" %}
{% set xtrabackup_24_rpm="percona-xtrabackup-24-2.4.15-1.el7.x86_64.rpm" %}
{% set xtrabackup_24_tarball="percona-xtrabackup-2.4.22-Linux-x86_64.glibc2.12.tar.gz" %}
{% set xtrabackup_24_tarball_dir="/usr/local/xtrabackup24" %}
{% set xtrabackup_80_rpm="percona-xtrabackup-80-8.0.25-17.1.el7.x86_64.rpm" %}
{% set xtrabackup_80_tarball="percona-xtrabackup-8.0.25-17-Linux-x86_64.glibc2.12.tar.gz" %}
{% set xtrabackup_80_tarball_dir="/usr/local/xtrabackup80" %}

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
    - name: {{ mysql_version }}" is error."
    - failhard: True
{% endif %}

{% set mysql_packet_unarchive_dir =  mysql_packet |replace(".tar.gz","") %}

sshkeys:
  ssh_auth.present:
    - user: root
    - enc: ssh-rsa
    - names:
      - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAi9Kp2AD7JzgKUslRwxtwez/YBYzOVAF6ZY5lZLXPohEhPiWWPHyD1wIYrNjsiR0dDZBNosmaLu+Bc4hevo/lSBWprldaiXkczLqLybv9NNNK2REdEzwRc3/uFgLzufzmOdjdlVvIzn7r35KkMpP7V6qIjB6k63BHvAB7FgfOq2WBddRuh04Me9jkEq01MZmrVvFBb1SUwkXYFc/HB0mhq8XKUzOAmmJiwM2QaR2HPfwuSZ5iDYXSrG0NGzaCybPWx7y1GuHjsU0wf6vSLGc8/gEH0+Q+h8e0tfQhGkrg+6X+Mzra5Pw4+XyCb33KhjFLju23uCq6aPT7GJsqdGs/ root@node3

install_repl:
  cmd.run:
    - name: yum -y install epel-release



install_percona_dependent:
  file.managed:
    - name: /tmp/percona-release-latest.noarch.rpm
    {% if mysql_version > 8000 %}
    - source: salt://percona/files/percona-release-latest.noarch.80.rpm
    {% else  %}
    - source: salt://percona/files/percona-release-latest.noarch.57.rpm
    {% endif %}
  cmd.run:
    - name: yum install -y /tmp/percona-release-latest.noarch.rpm
    - required:
      - file: install_percona_dependent


install_percona_toolkit:
  file.managed:
    - name: /tmp/{{ percona_toolkit }}
    - source: salt://percona/files/{{ percona_toolkit }}
  cmd.run:
    - name: yum install -y /tmp/{{ percona_toolkit }}
    - required:
      - file: install_percona_toolkit

{% if (salt['cmd.run']('rpm -qa|grep percona-xtrabackup-80|wc -l') | int == 0) and (mysql_version < 8000) %}
install_xtrabackup_24_rpm:
  file.managed:
    - name: /tmp/{{ xtrabackup_24_rpm }}
    - source: salt://percona/files/{{ xtrabackup_24_rpm }}
  cmd.run:
    - name: yum install -y /tmp/{{ xtrabackup_24_rpm }}
    - required:
      - file: /tmp/{{ xtrabackup_24_rpm }}
{% endif %}


{% if (salt['cmd.run']('rpm -qa|grep percona-xtrabackup-80|wc -l') | int != 0) and (mysql_version < 8000) %}
install_xtrabackup_24_tarball:
  file.managed:
    - name: /tmp/{{ xtrabackup_24_tarball }}
    - source: salt://percona/files/{{ xtrabackup_24_tarball }}
  cmd.run:
    - name: tar -zxf /tmp/{{ xtrabackup_24_tarball }} -C /usr/local/src/
    - required:
      - file: install_xtrabackup_24_tarball

/usr/local/src/{{ xtrabackup_24_tarball.replace('.tar.gz','') }}:
  file.symlink:
    - target: /usr/local/{{ xtrabackup_24_tarball_dirname }}
{% endif %}


{% if (salt['cmd.run']('rpm -qa|grep percona-xtrabackup-24|wc -l') | int == 0) and (mysql_version > 8000) %}
install_xtrabackup_80_rpm:
  file.managed:
    - name: /tmp/{{ xtrabackup_80_rpm }}
    - source: salt://percona/files/{{ xtrabackup_80_rpm }}
  cmd.run:
    - name: yum install -y /tmp/{{ xtrabackup_80_rpm }}
    - required:
      - file: /tmp/{{ xtrabackup_80_rpm }}
{% endif %}


{% if (salt['cmd.run']('rpm -qa|grep percona-xtrabackup-24|wc -l') | int != 0) and (mysql_version > 8000) %}
install_xtrabackup_80_tarball:
  file.managed:
    - name: /tmp/{{ xtrabackup_80_tarball }}
    - source: salt://percona/files/{{ xtrabackup_80_tarball }}
  cmd.run:
    - name: tar -zxf /tmp/{{ xtrabackup_80_tarball }} -C /usr/local/src/
    - required:
      - file: install_xtrabackup_24_tarball

/usr/local/src/{{ xtrabackup_80_tarball.replace('.tar.gz','') }}:
  file.symlink:
    - target: /usr/local/{{ xtrabackup_80_tarball_dirname }}
{% endif %}

unarchive_mysql_base:
  file.managed:
    - name: /usr/local/src/{{ mysql_packet }}
    - source: salt://percona/files/{{ mysql_packet }}
  cmd.run:
    - name: tar -zxf /usr/local/src/{{ mysql_packet }} -C /usr/local/ && ln -s /usr/local/{{ mysql_packet_unarchive_dir }} /usr/local/mysql{{ mysql_version }}
    - unless: test -d /usr/local/{{ mysql_packet_unarchive_dir }}
    - required:
      - file: unarchive_mysql_base

mysql_client:
  file.append:
    - name: /etc/profile
    - text: "export PATH=/usr/local/mysql{{ mysql_version | replace('.','') }}/bin:$PATH"

flush_os_cache:                          #脚本脚识
   cron.present:               #模板:cron 计划任务     功能：present
       - name: "sync && echo 3 >/proc/sys/vm/drop_caches && sleep 2 && echo 0>/proc/sys/vm/drop_caches"
       - user: root                  #添加到root的计划列表
       - minute: "0"
       - hour: "12"

