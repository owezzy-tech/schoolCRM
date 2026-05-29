#!/usr/bin/env sh

set -eu

key="${GOOGLE_MAPS_API_KEY:-}"

if [ -f /srv/index.html ]; then
  escaped_key=$(printf '%s' "$key" | sed 's/[\/&]/\\&/g')
  sed -i "s/__GOOGLE_MAPS_API_KEY__/${escaped_key}/g" /srv/index.html
fi

exec "$@"
