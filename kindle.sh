#!/bin/sh


configFile=/root/config.yml

if[ ! -f $configFile ]; then
  echo "请配置小说文件"
fi

if [ $(($randomNum % 2)) -eq 0 ]; then
   /root -d /root -f $configFile
fi