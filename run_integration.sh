#!/bin/sh

set -e

sleep 5

go test ./... -p 1 -tags=integration -count=1
