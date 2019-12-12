#!/bin/sh

BUILD_DIR=/usr/share/nginx/html
OUT_FILE=$BUILD_DIR/config.js

if [ -n "${API_ENDPOINT+set}" ]; then
  echo "window.API_ENDPOINT = '${API_ENDPOINT}';" >> $OUT_FILE
else
  echo "console.log('API_ENDPOINT environment variable not set and required to use this docker image');" >> $OUT_FILE
fi

exec "$@"