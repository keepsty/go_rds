{% set min_open_fds=65535 %}
{% set min_process=65535 %}
{% set disk_scheduler = 'noop' %}

{% if (salt['cmd.run']('cat /etc/selinux/config') ) == "SELINUX=disabled" %}
selinux_check_s:
  cmd.run:
    - name: echo 'selinux disabled'
{% else %}
failure:
  test.fail_without_changes:
    - name: "selinux is not disabled"
    - failhard: True
{% endif %}

pkg-install:
  pkg.installed:
    - pkgs:
      - numactl

{% if (salt['cmd.run']('numactl --hardware | awk "NR==1{print $2}"') )|int >1 %}
failure:
  test.fail_without_changes:
    - name: "selinux is not disabled"
    - failhard: True
{% else %}
selinux_check_s:
  cmd.run:
    - name: echo 'numa has one'
{% endif %}

{% for disk in grains['disks'] %}

{% if (disk != 'sr0') and salt['cmd.run']("cat /sys/block/{{ disk }}/queue/scheduler | egrep -q '\[{{ disk_scheduler }}\]|none' |wc -l ")|int <1 %}
check_disk_scheduler_{{ disk }}:
  cmd.run:
    - name: echo "{{ disk_scheduler }}" > /sys/block/{{ disk }}/queue/scheduler
    - failhard: True

{% endif %}
{% endfor %}


{% set data_dir = '/data' %}
{% if (salt['cmd.run'](" df -h '{{data_dir}}' | tail -n1 | awk '{print $(NF -1)}' ")|replace('%','')|int > 80) or not (salt['file.directory_exists'](data_dir)) %}
failure_dir_free:
  test.fail_without_changes:
    - name: "disk {{data_dir}} has used > 80% or dir {{data_dir}} does not exists"
    - failhard: True
{% endif %}




ft_type:
  cmd.run:
    - name: echo "{{data_dir}} file type, `df -T '/data' | awk 'NR==2'|awk '{print $(2)}'`"

{% if (salt['cmd.run']("cat /proc/mounts|grep `df -h '{{data_dir}}' | tail -n1 | awk '{print $1}'`|grep 'barrier=1' |wc -l ")|int >0) %}
failure:
  test.fail_without_changes:
    - name: "linux barrier has used, check failed"
    - failhard: True
{%else%}
ft_type:
  cmd.run:
    - name: echo "{{data_dir}} linux barrier has not used, check passed"
{% endif %}





