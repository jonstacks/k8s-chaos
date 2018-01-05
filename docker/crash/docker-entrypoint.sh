#!/bin/sh

# Sleep for CRASH_DELAY seconds. Default to 10.
sleep ${CRASH_DELAY:=10}
exit 2
