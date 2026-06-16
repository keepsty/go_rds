/etc/sysctl.conf:
  file.managed:
    - source: salt://linux_init/config/sysctl.conf
    - user: root
    - group: root
    - mode: 644
