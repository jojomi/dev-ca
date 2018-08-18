#!/bin/sh

set -ex

find dist -type f -name 'dev-ca' -exec upx -9 {} \;
