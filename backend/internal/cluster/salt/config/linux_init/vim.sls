/etc/vimrc:
  file.managed:
    - source: salt://init/config/vimrc
    - user: root
    - group: root
    - mode: 644
    - backup: '*'
