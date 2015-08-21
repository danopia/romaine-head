#!/bin/sh
set -e # stop on error

# Copyright Daniel Lamando 2015
# Builds and packages a go program for dev-mode chromebooks running Crouton
# Outputted file can be easily downloaded (to ~/Downloads by Chrome)
# The program can be launched from the chromeos shell as well as from a chroot:
# $ sh ~/Downloads/romaine-head.run -- [options]
# When launched on the chromeos side, uses root to install to $PATH.

# config for the entire script
APP_NAME=romaine-head
GO=/usr/local/go/bin/go
GO_PACKAGE=github.com/danopia/$APP_NAME
TEMP=/tmp/$APP_NAME
OUTPUT_PATH=~/Downloads/$APP_NAME.run
INSTALL_PATH=/usr/local/bin/$APP_NAME

# build binary
mkdir -p $TEMP
$GO build -o $TEMP/$APP_NAME $GO_PACKAGE

# add an install/launch script
cat << EOF > $TEMP/run
  if [ \$USER = root ]; then
    cp $APP_NAME $INSTALL_PATH
    $INSTALL_PATH \$*
  elif [ \$USER = chronos ]; then
    sudo cp $APP_NAME $INSTALL_PATH
    sudo $INSTALL_PATH \$*
  else
    # assuming we're in a chroot, just launch
    ./$APP_NAME \$*
  fi
EOF

# package self-running blob
rm $OUTPUT_PATH
makeself --nocomp $TEMP $OUTPUT_PATH $APP_NAME "sh ./run" 2>&1 >/dev/null \
  | grep -vE "(^./|test: x|^Header is \d+ lines)" # ignore common output
