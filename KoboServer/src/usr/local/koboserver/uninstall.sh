#!/bin/sh

# Uninstall koboserver

rm -f /etc/udev/rules.d/98-koboserver.rules
rm -rf /usr/local/koboserver/
rm -rf /mnt/onboard/.koboserver
