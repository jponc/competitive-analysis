#!/bin/sh

set -e

sleep 2
ts-node --transpile-only -e 'require("./UpMigration.ts").handler()'
