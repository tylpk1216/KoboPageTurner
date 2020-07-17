#!/bin/bash

KS_HOME=$(dirname $0)

if uname -a | grep -q x86
then
    #echo "PC detected"
    . $KS_HOME/config_pc.sh
else
    . $KS_HOME/config_kobo.sh
fi
