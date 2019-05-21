#!/bin/bash
#set -x
# Setup:
# 1 - vi $HOME/rubetekVideoCleaner/rvc_service
# 2 - chmod +x $HOME/rubetekVideoCleaner/rvc_service
# 3 - crontab -e
# 4 - */5 * * * *  /home/pi/rubetekVideoCleaner/rvc_service monitor

#
# Var definitions
#
RVC="/home/pi/rubetekVideoCleaner/rubetekVideoCleaner -cons -limit 50 -count 120 &"

#
#Functions
#
START(){
  eval $RVC
}

STOP(){
  ps -ef | grep rubetekVideoCleaner | awk '{print $2}' | xargs kill -9
}


PROCMON(){
  STATUS=`ps -ef | grep rubetekVideoCleaner | grep -v grep | wc -l`
  # STATUS = 0 means that has no processes running:
  if [ "$STATUS" == 0 ]; then
          START
  fi
}


#
# Script options
#
case $1 in
  start)
    START
    ;;

  stop)
    STOP
    ;;

  restart)
    STOP
    START
    ;;

  monitor)
	  PROCMON
	  ;;

  *)
    echo -e "Try $0 {start|stop|monitor|restart}"
    ;;
esac
