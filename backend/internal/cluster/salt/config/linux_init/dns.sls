/etc/resolv.conf:
  file.managed:
    - source: salt://linux_init/config/resolv.conf
    - user: root
    - group: root
    - mode: 644
    - backup: '*'
