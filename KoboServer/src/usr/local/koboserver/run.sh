#!/bin/sh

# load config
. $(dirname $0)/config.sh

# check if KoboServer contains the line "UNINSTALL"
if grep -q '^uninstall=true\s' $UserConfig; then
    echo "Uninstalling KoboServer!"
    $KS_HOME/uninstall.sh
    exit 0
fi

# check instance exist
running=$Logs/PID
if [ -e "$running" ]; then
    echo "Server is running."
    exit 0
fi

echo "Run" > $running
$Logs/$TestScript
$Logs/$Server &
#echo "`$Dt` done"
