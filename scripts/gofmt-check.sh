#!/usr/bin/env sh

unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
  echo "There are unformatted files:"
  echo $unformatted
  exit 1
fi
