yum_repo_release:
  pkg.installed:
    - sources:
      - epel-release: https://mirrors.aliyun.com/centos/7/extras/x86_64/Packages/epel-release-7-11.noarch.rpm
    - unless: rpm -qa | grep  epel-release-7
