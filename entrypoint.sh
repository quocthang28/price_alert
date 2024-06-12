#!/bin/sh

if [ ! -d /app/config ]; then
   mkdir -p /app/config
fi

if [ ! -f ./config/config.json ]; then
   mv config.json /app/config
fi

./main