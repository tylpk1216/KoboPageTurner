#!/bin/sh

#load config
. $(dirname $0)/config.sh

#check if KoboServer contains the line "UNINSTALL"
if grep -q '^UNINSTALL$' $UserConfig; then
    echo "Uninstalling KoboServer!"
    $KS_HOME/uninstall.sh
    exit 0
fi

$Logs/$Server &

echo "`$Dt` done"
