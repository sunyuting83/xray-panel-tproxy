#!/bin/sh

#set env

basepath=$(cd `dirname $0`; pwd)
pidpath="$basepath/Core/app_pid.pid"
configpath="$basepath/Core/config.json"
XrayCore="$basepath/Core/xray"

case $1 in
    start)
        $XrayCore run -config $configpath 2>&1 &
        echo $! > $pidpath
        echo "The process $! is running..."
        ;;
    stop)
        if [ -f $pidpath ];then
            pid=`cat $pidpath`
            # echo $pid
            kill -9 $pid
            echo "The process $pid is stop..."
            rm -rf $pidpath
        else
            echo "The process not running..."
        fi
        ;;
    reload)
        if [ -f $pidpath ];then
            pid=`cat $pidpath`
            # echo $pid
            kill -9 $pid
            echo "The process $pid was reloaded..."
            $XrayCore run -config $configpath 2>&1 &
            echo $! > $pidpath
            echo "The process $! is running..."
        else
            echo "The process not running..."
        fi
        ;;
    status)
        if [ -f $pidpath ];then
            pid=`cat $pidpath`
            echo "The process $pid is running"
        else
            echo "The process not running..."
        fi
        ;;
    *)
        echo "Usage:{start|stop|reload|status}"
        ;;

esac