#!/bin/sh

# udev kills slow scripts
if [ "$SETSID" != "1" ]
then
    SETSID=1 setsid "$0" "$@" &
    exit
fi

# load config
. $(dirname $0)/config.sh

# create work dirs
[ ! -e "$Logs" ] && mkdir -p "$Logs" >/dev/null 2>&1
[ ! -e "$UserConfig" ] && echo #UNINSTALL > $UserConfig

# copy files to $Logs
[ ! -e "$Logs/$Server" ] && cp $KS_HOME/$Server $Logs
[ ! -e "$Logs/$TestScript" ] && cp $KS_HOME/$TestScript $Logs
[ ! -e "$Logs/run.sh" ] && cp $KS_HOME/run.sh $Logs
[ ! -e "$Logs/config.sh" ] && cp $KS_HOME/config* $Logs

# output to log
$Logs/run.sh >> $Logs/koboserver.log 2>&1 &
