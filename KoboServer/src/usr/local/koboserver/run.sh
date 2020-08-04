#!/bin/sh

# load config
. $(dirname $0)/config.sh

# The $KS_HOME is $Logs at this point because we call run.sh in $Logs.

# check if KoboServer contains the line "uninstall=true"
if grep -q '^uninstall=true\s' $UserConfig; then
    echo "Uninstalling KoboServer!"
    $KS_HOME/uninstall.sh
    exit 0
fi

# check instance exist
running=$KS_HOME/PID
if [ -e "$running" ]; then
    echo "Server is running."
    exit 0
fi

echo "Run" > $running
$KS_HOME/$TestScript
$KS_HOME/$Server &
#echo "`$Dt` done"
