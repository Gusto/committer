#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]:-$0}"; )" &> /dev/null && pwd 2> /dev/null; )";
ROOT_DIR="$SCRIPT_DIR/.."
[ -f $ROOT_DIR/.bin/crom ] && exit 0

CROM_VERSION="0.4.2"
if [ "$(uname)" == "Darwin" ]; then
  if [ "$(arch)" == "arm64" ]; then
    export OS_ARCH="darwin-aarch64"
  else
    export OS_ARCH="darwin-x86_64"
  fi
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    export OS_ARCH="linux-gnu-x86_64"
fi

[ ! -d $ROOT_DIR/.bin ] && mkdir $ROOT_DIR/.bin
curl --location https://github.com/ethankhall/crom/releases/download/v$CROM_VERSION/crom-v$CROM_VERSION-$OS_ARCH.tgz | tar -xvz  -C $ROOT_DIR/.bin
chmod +x $ROOT_DIR/.bin/crom
