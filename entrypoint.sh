#!/bin/sh

set -e

check_and_delete_rule() {
    # Run ip rule show and check if fwmark 0x1 is in the output
    if ip rule show | grep -q "fwmark 0x1"; then
        # If fwmark 0x1 is present, delete the rule and route
        echo "Deleting rule and route..."
        ip rule delete fwmark 1 table 100
        ip route delete local default dev lo table 100
    else
        echo "Rule not found. Nothing to delete."
    fi
}

# 调用函数


reset_iptables(){
    check_and_delete_rule
    ip rule add fwmark 1 table 100
    ip route add local default dev lo table 100
    iptables -P INPUT ACCEPT
    iptables -P FORWARD ACCEPT
    iptables -P OUTPUT ACCEPT
    iptables -t nat -F
    iptables -t mangle -F
    iptables -t mangle -X
    iptables -t mangle -Z
    iptables -F
    iptables -X
}

set_xray_iptables(){
    iptables -t mangle -N XRAY
    iptables -t mangle -A XRAY -d 10.0.0.0/8 -j RETURN
    iptables -t mangle -A XRAY -d 100.64.0.0/10 -j RETURN
    iptables -t mangle -A XRAY -d 127.0.0.0/8 -j RETURN
    iptables -t mangle -A XRAY -d 169.254.0.0/16 -j RETURN
    iptables -t mangle -A XRAY -d 172.16.0.0/12 -j RETURN
    iptables -t mangle -A XRAY -d 192.0.0.0/24 -j RETURN
    iptables -t mangle -A XRAY -d 224.0.0.0/4 -j RETURN
    iptables -t mangle -A XRAY -d 240.0.0.0/4 -j RETURN
    iptables -t mangle -A XRAY -d 255.255.255.255/32 -j RETURN
    iptables -t mangle -A XRAY -s 192.168.1.45 -j RETURN -m mark --mark 1
    iptables -t mangle -A XRAY -d 192.168.0.0/16 -p tcp ! --dport 53 -j RETURN
    iptables -t mangle -A XRAY -d 192.168.0.0/16 -p udp ! --dport 53 -j RETURN
    iptables -t mangle -A XRAY -j RETURN -m mark --mark 0xff
    iptables -t mangle -A XRAY -p tcp -j TPROXY --on-port 7892 --tproxy-mark 1
    iptables -t mangle -A XRAY -p udp -j TPROXY --on-port 7892 --tproxy-mark 1
    iptables -t mangle -A PREROUTING -j XRAY

    iptables -t mangle -N DIVERT
    iptables -t mangle -A DIVERT -j MARK --set-mark 1
    iptables -t mangle -A DIVERT -j ACCEPT
    iptables -t mangle -I PREROUTING -p tcp -m socket -j DIVERT
}

reset_iptables
set_xray_iptables

/xpanel/server
exec "$@"
