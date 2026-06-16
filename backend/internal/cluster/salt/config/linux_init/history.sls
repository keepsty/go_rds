/etc/bashrc:
  file.append:
    - text:
      - HISTTIMEFORMAT="%F %T `whoami` "
  cmd.run:
    - name: source /etc/bashrc



