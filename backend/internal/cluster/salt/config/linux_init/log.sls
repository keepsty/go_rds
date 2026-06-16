## 使用logger将输入的命令写入到memssages的一个简单功能，使用SaltStack的File模块的Append方法。建议将memssages日志文件进行统一收集管理，建议使用ELK Stack(Elasticsearch、LogStach、Kibana)。

append_log:
  file.append:
    - name: /etc/bashrc
    - text:
      - export PROMPT_COMMAND='{ msg=$(history 1 | { read x y; echo $y; });logger "[euid=$(whoami)]":$(who am i):[`pwd`]"$msg"; }'
  cmd.run:
    - name: source /etc/bashrc
