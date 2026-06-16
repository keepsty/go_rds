# 建议在生产服务器对ssh配件文件进行统一管理，修改默认的连接端口
sync-ssh:
  file.managed:
    - name: /etc/ssh/sshd_config
    - source: salt://linux_init/config/sshd_config
    - user: root
    - group: root
    - mode: 644
  cmd.run:
    - name: /etc/init.d/sshd restart
    - require:
      - file: sync-ssh
  service.running:
    - name: sshd
    - enable: True
    - reload: True
    - require:
      - file: sync-ssh
