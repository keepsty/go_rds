{% set proxysql_rpm='[[.ProxysqlRpm ]]' %}
{% set proxysql_port=[[.AdminPort ]] %}
{% set proxysql_datadir='[[.ProxysqlDir ]]' %}
{% set proxysql_conf_dir = '[[.ProxysqlConfDir ]]' %}
{% set netstat_admin_port_check = 'netstat -tulpn|grep [[.AdminPort ]] |wc -l ' %}
{% set netstat_buss_port_check = 'netstat -tulpn|grep 1[[.AdminPort ]] |wc -l ' %}
{% set lsof_admin_port_check = 'lsof -i:[[.AdminPort ]] |wc -l ' %}


proxysql.packages:
  pkg.installed:
    - pkgs:
      - perl-DBD-MySQL
      - gnutls
      - python-devel
      - MySQL-python

sysctl_net.ipv4.ip_local_reserved_ports:
  sysctl.present:
    - name: net.ipv4.ip_local_reserved_ports
    - value: 6032-6060,16032-16060


proxy_rpm_manage:
  file.managed:
    - name: /tmp/{{proxysql_rpm}}
    - source: salt://proxysql/files/{{proxysql_rpm}}
    - unless: test -f /tmp/{{proxysql_rpm}}
  cmd.run:
    - name: yum install -y /tmp/{{proxysql_rpm}}
    - required:
      - file: proxy_rpm_manage

{% for f in ['/etc/proxysql.cnf','/etc/systemd/system/proxysql.service','/etc/systemd/system/proxysql.service.bak'] %}
mv_proxysql_service_{{ f }}:
  cmd.run:
    - name: mv {{ f }} {{ f }}_`date "+%Y%m%d"`
    - onlyif: test -f {{ f }}
{% endfor %}


{% if salt['file.directory_exists']( "{{ proxysql_datadir }}" ) %}
failure:
  test.fail_without_changes:
    - name: {{ proxysql_datadir }}"directory exists"
    - failhard: True
{% else %}
init_proxysql_dir:
  cmd.run:
    - name: mkdir -p {{ proxysql_datadir }}
{% endif %}



{% if (salt['cmd.run'](netstat_admin_port_check) | int == 0)  and (salt['cmd.run'](netstat_buss_port_check) | int == 0) and (salt['cmd.run'](lsof_admin_port_check)) |int ==0 %}
install_proxysql:
  file.managed:
    - name: {{ proxysql_datadir }}/proxysql.cnf
    - source: salt://{{ proxysql_conf_dir }}/proxysql.cnf
    - mode: 0644
    - required:
      - cmd: init_proxysql_dir
  cmd.run:
    - name: proxysql --idle-threads --no-version-check -c {{ proxysql_datadir }}/proxysql.cnf;sleep 30;
    - failhard: True
    - required:
      - file: install_proxysql

init_proxysql_user:
  file.managed:
    - name: /tmp/init_proxysql_user.sh
    - source: salt://{{ proxysql_conf_dir }}/init_proxysql_user.sh
    - template: jinja    # 增加这行表示开启模板
    - defaults:          # 下面设定变量的值
      MysqlVersion: [[.MysqlVersion ]]
  cmd.run:
    - name: sh /tmp/init_proxysql_user.sh
    - failhard: True
    - required:
      - file: init_proxysql_user

{% else %}
failure:
  test.fail_without_changes:
    - name: "port {{proxysql_port}} has been used"
    - failhard: True
{% endif %}

file_remove:
  cmd.run:
    - name: rm -rf /tmp/init_proxysql_user.sh
